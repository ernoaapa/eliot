package runtime

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/containerd/containerd"
	namespaces "github.com/containerd/containerd/api/services/namespaces/v1"
	tasks "github.com/containerd/containerd/api/services/tasks/v1"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/plugin"
	"github.com/containerd/containerd/remotes"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/progress"
	opts "github.com/ernoaapa/can/pkg/runtime/containerd"
	"github.com/ernoaapa/can/pkg/runtime/containerd/mapping"
	imagespecs "github.com/opencontainers/image-spec/specs-go/v1"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

var (
	snapshotter = "overlayfs"
)

// ContainerdClient is containerd client wrapper
type ContainerdClient struct {
	context context.Context
	timeout time.Duration
	address string
}

// NewContainerdClient creates new containerd client with given timeout
func NewContainerdClient(context context.Context, timeout time.Duration, address string) *ContainerdClient {
	return &ContainerdClient{
		context: context,
		timeout: timeout,
		address: address,
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
		podName := mapping.GetPodName(container)
		if _, ok := pods[podName]; !ok {
			pod := model.NewPod(podName, namespace)
			pods[podName] = &pod
		}

		pods[podName].AppendContainer(
			mapping.MapContainerToInternalModel(container),
			mapping.MapContainerStatusToInternalModel(container, resolveContainerStatus(ctx, container)),
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

func getValues(podsByName map[string]*model.Pod) (result []model.Pod) {
	for _, pod := range podsByName {
		result = append(result, *pod)
	}
	return result
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
	return model.Pod{}, fmt.Errorf("Pod in namespace [%s] with name [%s] not found", namespace, podName)
}

// CreateContainer creates given container
func (c *ContainerdClient) CreateContainer(pod model.Pod, container model.Container) error {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connectionErr := c.getConnection(pod.Metadata.Namespace)
	if connectionErr != nil {
		return connectionErr
	}

	image, imageErr := client.GetImage(ctx, container.Image)
	if imageErr != nil {
		return imageErr
	}

	specOpts := []containerd.SpecOpts{
		containerd.WithImageConfig(image),
	}

	if len(container.Args) > 0 {
		specOpts = append(specOpts, containerd.WithProcessArgs(container.Args...))
	}

	if container.Tty {
		specOpts = append(specOpts, containerd.WithTTY)
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
			return errors.Wrapf(err, "Error while ensuring mount source directories exist")
		}
		log.Debugf("Adding %d mounts to container", len(container.Mounts))
		specOpts = append(specOpts, opts.WithMounts(container.Mounts))
	}

	if pod.Spec.HostNetwork {
		specOpts = append(specOpts,
			containerd.WithHostNamespace(specs.NetworkNamespace),
			containerd.WithHostHostsFile,
			containerd.WithHostResolvconf,
		)
	}

	if pod.Spec.HostPID {
		specOpts = append(specOpts, containerd.WithHostNamespace(specs.PIDNamespace))
	}

	log.Debugf("Create new container from image %s...", image.Name())
	_, err := client.NewContainer(ctx,
		container.Name,
		containerd.WithContainerLabels(mapping.NewLabels(pod, container)),
		containerd.WithNewSpec(specOpts...),
		containerd.WithSnapshotter(snapshotter),
		containerd.WithNewSnapshot(container.Name, image),
		containerd.WithRuntime(fmt.Sprintf("%s.%s", plugin.RuntimePlugin, "linux"), nil),
	)
	if err != nil {
		return errors.Wrapf(err, "Failed to create new container from image %s", image.Name())
	}
	return nil
}

// StartContainer starts the already created container
func (c *ContainerdClient) StartContainer(namespace, name string, tty bool) (result model.ContainerStatus, err error) {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connectionErr := c.getConnection(namespace)
	if connectionErr != nil {
		return result, connectionErr
	}

	container, err := client.LoadContainer(ctx, name)
	if err != nil {
		return result, errors.Wrapf(err, "Failed to load container [%s], cannot start it", name)
	}

	log.Debugf("Create task in container: %s", container.ID())
	io, err := containerd.NewDirectIO(ctx, tty)
	if err != nil {
		return result, errors.Wrapf(err, "Error while creating container task IO")
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

	return mapping.MapContainerStatusToInternalModel(container, resolveContainerStatus(ctx, container)), nil
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

	task, err := container.Task(ctx, nil)
	if err != nil {
		if !errdefs.IsNotFound(err) {
			return result, errors.Wrap(err, "Fetching container task returned unexpected error")
		}
	}

	if task != nil {
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
		ContainerID: container.ID(),
		Image:       container.Info().Image,
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

	_, err = client.Pull(
		ctx,
		ref,
		containerd.WithPullUnpack,
		containerd.WithSchema1Conversion,
		containerd.WithImageHandler(images.HandlerFunc(handler)),
	)
	if err != nil {
		return errors.Wrapf(err, "Error while pulling image [%s] to namespace [%s]", ref, namespace)
	}

	progress.AllDone()

	return nil
}

// GetNamespaces return all namespaces what cand manages
func (c *ContainerdClient) GetNamespaces() ([]string, error) {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connErr := c.getConnection(model.DefaultNamespace)
	if connErr != nil {
		return nil, connErr
	}

	resp, err := client.NamespaceService().List(ctx, &namespaces.ListNamespacesRequest{})
	if err != nil {
		return nil, err
	}

	return getNamespaces(resp.Namespaces), nil
}

func getNamespaces(namespaces []namespaces.Namespace) (result []string) {
	for _, namespace := range namespaces {
		if namespace.Name != "default" {
			result = append(result, namespace.Name)
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

// Attach returns pod logs
func (c *ContainerdClient) Attach(namespace, name string, io AttachIO) error {
	ctx, cancel := c.getContext()
	defer cancel()

	client, err := c.getConnection(namespace)
	if err != nil {
		return errors.Wrapf(err, "Unable to get connection for streaming logs")
	}

	container, err := client.LoadContainer(ctx, name)
	if err != nil {
		return errors.Wrapf(err, "Cannot return container logs for container [%s] in namespace [%s]", name, namespace)
	}

	task, taskErr := container.Task(ctx, containerd.WithAttach(io.Stdin, io.Stdout, io.Stderr))
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
