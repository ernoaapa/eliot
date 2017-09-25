package runtime

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/containerd/containerd"
	namespaces "github.com/containerd/containerd/api/services/namespaces/v1"
	tasks "github.com/containerd/containerd/api/services/tasks/v1"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/plugin"
	"github.com/ernoaapa/can/pkg/model"
	mapper "github.com/ernoaapa/can/pkg/runtime/containerd"
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

// GetContainers return all containers active in containerd grouped by pod name
func (c *ContainerdClient) GetContainers(namespace string) (map[string][]model.Container, error) {
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
	return mapper.MapToModelByPodNames(containers), nil
}

// CreateContainer creates given container
func (c *ContainerdClient) CreateContainer(pod model.Pod, container model.Container) error {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connectionErr := c.getConnection(pod.GetNamespace())
	if connectionErr != nil {
		return connectionErr
	}

	image, pullErr := c.ensureImagePulled(pod.GetNamespace(), container.Image)
	if pullErr != nil {
		return pullErr
	}

	spec, specErr := containerd.GenerateSpec(ctx, client, nil, containerd.WithImageConfig(image))
	if specErr != nil {
		return specErr
	}

	log.Debugf("Create new container from image %s...", image.Name())
	_, err := client.NewContainer(ctx,
		container.ID,
		containerd.WithContainerLabels(mapper.NewContainerLabels(pod, container)),
		containerd.WithSpec(spec),
		containerd.WithSnapshotter(snapshotter),
		containerd.WithNewSnapshotView(container.ID, image),
		containerd.WithRuntime(fmt.Sprintf("%s.%s", plugin.RuntimePlugin, "linux"), nil),
	)
	if err != nil {
		return errors.Wrapf(err, "Failed to create new container from image %s", image.Name())
	}
	return nil
}

// StartContainer starts the already created container
func (c *ContainerdClient) StartContainer(containerID string) error {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connectionErr := c.getConnection(model.DefaultNamespace)
	if connectionErr != nil {
		return connectionErr
	}

	container, err := client.LoadContainer(ctx, containerID)
	if err != nil {
		return errors.Wrapf(err, "Failed to load container with id %s, cannot start it", containerID)
	}

	log.Debugf("Create task in container: %s", container.ID())
	task, err := container.NewTask(ctx, containerd.NullIO)
	if err != nil {
		return errors.Wrapf(err, "Error while creating task for container [%s]", container.ID())
	}

	log.Debugln("Starting task...")
	err = task.Start(ctx)
	if err != nil {
		return errors.Wrapf(err, "Failed to start task in container", container.ID())
	}
	log.Debugf("Task started (pid %d)", task.Pid())
	return nil
}

// StopContainer stops given container
func (c *ContainerdClient) StopContainer(containerID string) error {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connectionErr := c.getConnection(model.DefaultNamespace)
	if connectionErr != nil {
		return connectionErr
	}

	container, err := client.LoadContainer(ctx, containerID)
	if err != nil {
		return errors.Wrapf(err, "Failed to load container with id %s, cannot stop it", containerID)
	}

	task, err := container.Task(ctx, nil)
	if err == nil {
		task.Delete(ctx, containerd.WithProcessKill)
	}
	if err := container.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
		return errors.Wrapf(err, "Failed to delete container %s", container.ID())
	}
	return nil
}

func (c *ContainerdClient) ensureImagePulled(namespace, ref string) (image containerd.Image, err error) {
	ctx, cancel := c.getContext()
	defer cancel()

	client, err := c.getConnection(namespace)
	if err != nil {
		return image, err
	}

	image, err = client.Pull(ctx, ref)
	if err != nil {
		return image, errors.Wrapf(err, "Error pulling image [%s]", ref)
	}

	log.Debugf("Unpacking container image [%s]...", image.Target().Digest)
	err = image.Unpack(ctx, snapshotter)
	if err != nil {
		return image, errors.Wrapf(err, "Error while unpacking image [%s]", image.Target().Digest)
	}

	return image, nil
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
func (c *ContainerdClient) IsContainerRunning(containerID string) (bool, error) {
	ctx, cancel := c.getContext()
	defer cancel()

	client, connErr := c.getConnection(model.DefaultNamespace)
	if connErr != nil {
		return false, connErr
	}

	container, loadErr := client.LoadContainer(ctx, containerID)
	if loadErr != nil {
		return false, errors.Wrapf(loadErr, "Failed to load container with id %s, cannot resolve running state", containerID)
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
func (c *ContainerdClient) GetContainerTaskStatus(containerID string) string {
	ctx, cancel := c.getContext()
	defer cancel()

	client, err := c.getConnection(model.DefaultNamespace)
	if err != nil {
		log.Warnf("Unable to get connection for resolving task status for containerID %s", containerID)
		return "UNKNOWN"
	}

	resp, err := client.TaskService().Get(ctx, &tasks.GetRequest{
		ContainerID: containerID,
	})
	if err != nil {
		log.Warnf("Unable to resolve Container task status: %s", err)
		return "UNKNOWN"
	}

	return resp.Process.Status.String()
}

// GetLogs returns pod logs
func (c *ContainerdClient) GetLogs(namespace, containerID string, stdin io.Reader, stdout, stderr io.Writer) error {
	ctx, cancel := c.getContext()
	defer cancel()

	client, err := c.getConnection(namespace)
	if err != nil {
		return errors.Wrapf(err, "Unable to get connection for streaming logs")
	}

	container, err := client.LoadContainer(ctx, containerID)
	if err != nil {
		return errors.Wrapf(err, "Cannot return container logs for containerID [%s] in namespace [%s]", containerID, namespace)
	}

	task, taskErr := container.Task(ctx, containerd.WithAttach(stdin, stdout, stderr))
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
