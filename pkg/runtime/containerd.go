package runtime

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/containerd/containerd"
	tasks "github.com/containerd/containerd/api/services/tasks/v1"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	namespaceutils "github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/containerd/plugin"
	"github.com/containerd/containerd/remotes"
	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/ernoaapa/eliot/pkg/progress"
	opts "github.com/ernoaapa/eliot/pkg/runtime/containerd"
	"github.com/ernoaapa/eliot/pkg/runtime/containerd/extensions"
	"github.com/ernoaapa/eliot/pkg/runtime/containerd/mapping"
	imagespecs "github.com/opencontainers/image-spec/specs-go/v1"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
)

// ContainerdClient is containerd client wrapper
type ContainerdClient struct {
	context     context.Context
	timeout     time.Duration
	snapshotter string
	address     string
	hostname    string
}

// NewContainerdClient creates new containerd client with given timeout
func NewContainerdClient(context context.Context, timeout time.Duration, snapshotter, address, hostname string) *ContainerdClient {
	return &ContainerdClient{
		context:     context,
		timeout:     timeout,
		address:     address,
		snapshotter: snapshotter,
		hostname:    hostname,
	}
}

func (c *ContainerdClient) getContext() (context.Context, context.CancelFunc) {
	var (
		ctx    = c.context
		cancel context.CancelFunc
	)

	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	return ctx, cancel
}

func (c *ContainerdClient) getConnection(namespace string) (*containerd.Client, error) {
	client, err := containerd.New(c.address, containerd.WithDefaultNamespace(namespace))
	if err != nil {
		return client, errors.Wrapf(err, "Unable to create connection to containerd")
	}
	return client, nil
}

// GetPods return all containers active in containerd grouped by pods
func (c *ContainerdClient) GetPods(namespace string) ([]model.Pod, error) {
	pods := map[string]*model.Pod{}
	ctx, cancel := c.getContext()
	defer cancel()

	client, err := c.getConnection(namespace)
	if err != nil {
		return nil, err
	}

	containers, err := client.Containers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting list of containers")
	}

	for _, container := range containers {
		info, err := container.Info(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "Error while fetching container info")
		}
		podName := mapping.GetPodName(info)
		if _, ok := pods[podName]; !ok {
			pod := mapping.InitialisePodModel(info, namespace, podName, c.hostname)
			pods[pod.Metadata.Name] = &pod
		}

		pods[podName].AppendContainer(
			mapping.MapContainerToInternalModel(info),
			mapping.MapContainerStatusToInternalModel(info, resolveContainerStatus(ctx, container)),
		)
	}

	return getValues(pods), nil
}

func resolveContainerStatus(ctx context.Context, container containerd.Container) containerd.Status {
	status := containerd.Status{}
	task, err := container.Task(ctx, nil)
	if err != nil {
		log.Warnf("Cannot resolve container status, failed to fetch task, will mark as unknown. Error: %s", err)
	} else {
		status, err = task.Status(ctx)
		if err != nil {
			log.Warnf("Cannot resolve container status, failed to fetch status, will mark as unknown. Error: %s", err)
		}
	}

	return status
}

// GetPod return pod by name
func (c *ContainerdClient) GetPod(namespace, podName string) (model.Pod, error) {
	pods, err := c.GetPods(namespace)
	if err != nil {
		return model.Pod{}, err
	}
	for _, pod := range pods {
		if pod.Metadata.Name == podName {
			return pod, nil
		}
	}
	return model.Pod{}, ErrWithMessagef(ErrNotFound, "Pod in namespace [%s] with name [%s] not found", namespace, podName)
}

// CreateContainer creates given container
func (c *ContainerdClient) CreateContainer(pod model.Pod, container model.Container) (status model.ContainerStatus, err error) {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connectionErr := c.getConnection(pod.Metadata.Namespace)
	if connectionErr != nil {
		return status, connectionErr
	}

	image, imageErr := client.GetImage(ctx, container.Image)
	if imageErr != nil {
		return status, imageErr
	}

	specOpts := []oci.SpecOpts{
		oci.WithImageConfig(image),
	}

	if len(container.Args) > 0 {
		specOpts = append(specOpts, oci.WithProcessArgs(container.Args...))
	}

	if container.Tty {
		specOpts = append(specOpts, oci.WithTTY)
	}

	if container.WorkingDir != "" {
		specOpts = append(specOpts, opts.WithCwd(container.WorkingDir))
	}

	if len(container.Env) > 0 {
		log.Debugf("Adding %d environment variables", len(container.Env))
		specOpts = append(specOpts, opts.WithEnv(container.Env))
	}

	if len(container.Mounts) > 0 {
		err := ensureMountSourceDirExists(container.Mounts)
		if err != nil {
			return status, errors.Wrapf(err, "Error while ensuring mount source directories exist")
		}
		log.Debugf("Adding %d mounts to container", len(container.Mounts))
		specOpts = append(specOpts, opts.WithMounts(container.Mounts))
	}

	log.Debugf("I am before devices")
	if (len(container.Devices) > 0) {
		for _, dev := range container.Devices {
			log.Debugf("Adding device to container...")
			specOpts = append(specOpts, opts.WithDevice(dev.DeviceType, int64(dev.MajorId), int64(dev.MinorId)))
		}
	}
	log.Debugf("I am behind devices")

	if pod.Spec.HostNetwork {
		specOpts = append(specOpts,
			oci.WithHostNamespace(specs.NetworkNamespace),
			oci.WithHostHostsFile,
			oci.WithHostResolvconf,
		)
	}

	if pod.Spec.HostPID {
		specOpts = append(specOpts, oci.WithHostNamespace(specs.PIDNamespace))
	}

	id := xid.New()
	containerOpts := []containerd.NewContainerOpts{
		containerd.WithContainerLabels(mapping.NewLabels(pod, container)),
		containerd.WithNewSpec(specOpts...),
		containerd.WithSnapshotter(c.snapshotter),
		containerd.WithNewSnapshot(id.String(), image),
		containerd.WithRuntime(fmt.Sprintf("%s.%s", plugin.RuntimePlugin, "linux"), nil),
		extensions.WithLifecycleExtension,
	}

	if container.Pipe != nil {
		containerOpts = append(containerOpts, extensions.WithPipeExtension(
			mapping.MapPipeToContainerdModel(*container.Pipe),
		))
	}

	log.Debugf("Create new container from image %s...", image.Name())
	created, err := client.NewContainer(
		namespaceutils.WithNamespace(ctx, pod.Metadata.Namespace),
		id.String(),
		containerOpts...,
	)
	if err != nil {
		return status, errors.Wrapf(err, "Failed to create new container from image %s", image.Name())
	}

	info, err := created.Info(ctx)
	if err != nil {
		return status, errors.Wrap(err, "Error while fetching container info")
	}

	return mapping.MapContainerStatusToInternalModel(info, resolveContainerStatus(ctx, created)), nil
}

// StartContainer starts the pre-created container
func (c *ContainerdClient) StartContainer(namespace, id string, ioSet IOSet) (result model.ContainerStatus, err error) {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connectionErr := c.getConnection(namespace)
	if connectionErr != nil {
		return result, connectionErr
	}

	container, err := client.LoadContainer(ctx, id)
	if err != nil {
		return result, errors.Wrapf(err, "Failed to load container [%s], cannot start it", id)
	}

	info, err := container.Info(ctx)
	if err != nil {
		return result, errors.Wrap(err, "Error while fetching container info")
	}

	log.Debugf("Create task in container: %s", container.ID())
	io, err := opts.NewDirectIO(ctx, ioSet.Stdin, ioSet.Stdout, ioSet.Stderr, mapping.RequireTty(info))
	if err != nil {
		return result, errors.Wrapf(err, "Error while creating container task IO")
	}

	if task, err := container.Task(ctx, nil); err != nil {
		if !errdefs.IsNotFound(err) {
			return result, errors.Wrapf(err, "Error while resolving container task status")
		}
	} else {
		if err := ensureTaskStopped(ctx, task); err != nil {
			return result, errors.Wrapf(err, "Failed to ensure task is stopped")
		}
		if _, err := task.Delete(ctx); err != nil {
			return result, errors.Wrapf(err, "Error while cleaning up old container task")
		}
	}

	task, err := container.NewTask(ctx, io.IOCreate)
	if err != nil {
		return result, errors.Wrapf(err, "Error while creating task for container [%s]", container.ID())
	}

	log.Debugln("Starting task...")
	err = task.Start(ctx)
	if err != nil {
		return result, errors.Wrapf(err, "Failed to start task in container [%s]", container.ID())
	}
	log.Debugf("Task started (pid %d)", task.Pid())

	if err := container.Update(ctx, extensions.IncrementRestart); err != nil {
		return result, errors.Wrapf(err, "Failed to increment container [%s] start counter", container.ID())
	}

	return mapping.MapContainerStatusToInternalModel(info, resolveContainerStatus(ctx, container)), nil
}

func ensureTaskStopped(ctx context.Context, task containerd.Task) error {
	status, err := task.Status(ctx)
	if err != nil {
		return errors.Wrapf(err, "Failed to resolve task status")
	}
	switch status.Status {
	case containerd.Running, containerd.Paused, containerd.Pausing:
		return task.Kill(ctx, syscall.SIGTERM)
	}
	return nil
}

// StopContainer stops given container
func (c *ContainerdClient) StopContainer(namespace, name string) (result model.ContainerStatus, err error) {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connectionErr := c.getConnection(namespace)
	if connectionErr != nil {
		return result, connectionErr
	}

	container, err := client.LoadContainer(ctx, name)
	if err != nil {
		return result, errors.Wrapf(err, "Failed to load container [%s], cannot stop it", name)
	}

	info, err := container.Info(ctx)
	if err != nil {
		return result, errors.Wrap(err, "Error while fetching container info")
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return result, errors.Wrap(err, "Fetching container task returned unexpected error")
		}
	}

	if task != nil {
		if err := ensureTaskStopped(ctx, task); err != nil {
			log.Warnf("Failed to kill task with SIGTERM, will next force kill. Error: %s", err)
		}

		_, taskDeleteErr := task.Delete(ctx, containerd.WithProcessKill)
		if err != nil {
			return result, errors.Wrapf(taskDeleteErr, "Container task deletion returned error")
		}
	}

	if err := container.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
		// Someone might already deleted it...
		if !errdefs.IsNotFound(err) {
			return result, errors.Wrapf(err, "Failed to delete container [%s]", container.ID())
		}
	}

	return model.ContainerStatus{
		ContainerID: info.ID,
		Image:       info.Image,
		State:       "stopped",
	}, nil
}

// Signal will send a syscall.Signal to the container task process
func (c *ContainerdClient) Signal(namespace, name string, signal syscall.Signal) error {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connectionErr := c.getConnection(namespace)
	if connectionErr != nil {
		return connectionErr
	}

	container, err := client.LoadContainer(ctx, name)
	if err != nil {
		return errors.Wrapf(err, "Failed to load container [%s], cannot send signal", name)
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "Unable to get task in container [%s], cannot send signal", name)
	}

	log.Debugf("Send signal [%s] to container all tasks", signal)
	return task.Kill(ctx, signal, containerd.WithKillAll)
}

// PullImage ensures that given container image is pulled to the namespace
func (c *ContainerdClient) PullImage(namespace, ref string, progress *progress.ImageFetch) error {
	ctx, cancel := c.getContext()
	defer cancel()

	client, err := c.getConnection(namespace)
	if err != nil {
		return err
	}

	done := make(chan struct{})
	defer close(done)
	go opts.UpdateFetchProgress(done, client, progress)

	handler := func(ctx context.Context, desc imagespecs.Descriptor) ([]imagespecs.Descriptor, error) {
		if desc.MediaType != images.MediaTypeDockerSchema1Manifest {
			progress.Add(remotes.MakeRefKey(ctx, desc), desc.Digest.String())
		}
		return nil, nil
	}

	img, err := client.Pull(
		ctx,
		ref,
		containerd.WithSchema1Conversion,
		containerd.WithImageHandler(images.HandlerFunc(handler)),
	)
	if err != nil {
		return errors.Wrapf(err, "Error while pulling image [%s] to namespace [%s]", ref, namespace)
	}

	supported, err := images.Platforms(ctx, img.ContentStore(), img.Target())
	if err != nil {
		return errors.Wrapf(err, "Error while resolving image [%s] supported platforms", ref)
	}

	if !platformExist(platforms.DefaultSpec(), supported) {
		platformNames := []string{}
		for _, platform := range supported {
			platformNames = append(platformNames, platforms.Format(platform))
		}
		return ErrWithMessagef(ErrNotSupported, "Image [%s] does not support [%s]. Supported platforms: %s", ref, platforms.Default(), strings.Join(platformNames, ","))
	}

	available, _, _, _, err := images.Check(ctx, img.ContentStore(), img.Target(), platforms.Default())
	if err != nil {
		return errors.Wrapf(err, "Error while checking image [%s] availability", ref)
	}

	if !available {
		return ErrWithMessagef(ErrNotSupported, "Image [%s] does not available for [%s/%s]", ref, runtime.GOOS, runtime.GOARCH)
	}

	if err := img.Unpack(ctx, c.snapshotter); err != nil {
		return errors.Wrapf(err, "Error while unpacking image [%s] to namespace [%s]", ref, namespace)
	}

	progress.AllDone()

	return nil
}

func platformExist(platform imagespecs.Platform, supported []imagespecs.Platform) bool {
	matcher := platforms.NewMatcher(platform)
	for _, platform := range supported {
		if matcher.Match(platform) {
			return true
		}
	}
	return false
}

// GetNamespaces return all namespaces
func (c *ContainerdClient) GetNamespaces() ([]string, error) {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connErr := c.getConnection(model.DefaultNamespace)
	if connErr != nil {
		return nil, connErr
	}

	resp, err := client.NamespaceService().List(ctx)
	if err != nil {
		return nil, err
	}

	return getNamespaces(resp), nil
}

func getNamespaces(namespaces []string) (result []string) {
	for _, namespace := range namespaces {
		if namespace != "default" {
			result = append(result, namespace)
		}
	}
	return result
}

// IsContainerRunning returns true if container running. If cannot resolve, return false with error
func (c *ContainerdClient) IsContainerRunning(namespace, name string) (bool, error) {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connErr := c.getConnection(namespace)
	if connErr != nil {
		return false, connErr
	}

	container, loadErr := client.LoadContainer(ctx, name)
	if loadErr != nil {
		return false, errors.Wrapf(loadErr, "Failed to load container [%s], cannot resolve running state", name)
	}

	_, err := container.Task(ctx, nil)
	if err != nil {
		if errdefs.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetContainerTaskStatus resolves container status or return UNKNOWN
func (c *ContainerdClient) GetContainerTaskStatus(namespace, name string) string {
	ctx, cancel := c.getContext()
	defer cancel()

	client, err := c.getConnection(namespace)
	if err != nil {
		log.Warnf("Unable to get connection for resolving task status for container %s", name)
		return "UNKNOWN"
	}

	resp, err := client.TaskService().Get(ctx, &tasks.GetRequest{
		ContainerID: name,
	})
	if err != nil {
		log.Warnf("Unable to resolve Container task status: %s", err)
		return "UNKNOWN"
	}

	return resp.Process.Status.String()
}

// Exec run command in container and hook IO to the new process
func (c *ContainerdClient) Exec(namespace, name, id string, args []string, tty bool, io AttachIO) error {
	ctx, cancel := c.getContext()
	defer cancel()
	ctx = namespaces.WithNamespace(ctx, namespace)

	client, err := c.getConnection(namespace)
	if err != nil {
		return errors.Wrapf(err, "Unable to get connection to execute command")
	}

	container, err := client.LoadContainer(ctx, name)
	if err != nil {
		return errors.Wrapf(err, "Cannot execute command in container [%s] in namespace [%s]", name, namespace)
	}

	spec, err := container.Spec(ctx)
	if err != nil {
		return err
	}

	task, taskErr := container.Task(ctx, nil)
	if taskErr != nil {
		return taskErr
	}

	pspec := spec.Process
	pspec.Terminal = tty
	pspec.Args = args

	process, err := task.Exec(ctx, id, pspec, cio.NewCreator(
		cio.WithStreams(io.Stdin, io.Stdout, io.Stderr),
		cio.WithTerminal,
	))
	if err != nil {
		return err
	}
	defer process.Delete(ctx)

	status, err := process.Wait(ctx)
	if err != nil {
		return err
	}

	if err := process.Start(ctx); err != nil {
		return err
	}

	exitStatus := <-status
	return exitStatus.Error()
}

// Attach hook IO to container main process
func (c *ContainerdClient) Attach(namespace, name string, io AttachIO) error {
	ctx, cancel := c.getContext()
	defer cancel()

	client, err := c.getConnection(namespace)
	if err != nil {
		return errors.Wrapf(err, "Unable to get connection to attach into container")
	}

	container, err := client.LoadContainer(ctx, name)
	if err != nil {
		return errors.Wrapf(err, "Cannot attach to container [%s] in namespace [%s]", name, namespace)
	}

	task, taskErr := container.Task(ctx, cio.NewAttach(
		cio.WithStreams(io.Stdin, io.Stdout, io.Stderr),
	))
	if taskErr != nil {
		return taskErr
	}

	status, err := task.Wait(ctx)
	if err != nil {
		return err
	}

	exitStatus := <-status
	return exitStatus.Error()
}
