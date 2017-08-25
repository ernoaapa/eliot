package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/plugin"
	"github.com/ernoaapa/layeryd/model"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// ContainerdClient is containerd client wrapper
type ContainerdClient struct {
	client    *containerd.Client
	context   context.Context
	timeout   time.Duration
	address   string
	namespace string
}

// NewContainerdClient creates new containerd client with given timeout and namespace
func NewContainerdClient(context context.Context, timeout time.Duration, address, namespace string) *ContainerdClient {
	return &ContainerdClient{
		context:   context,
		timeout:   timeout,
		address:   address,
		namespace: namespace,
	}
}

func (c *ContainerdClient) getContext() (context.Context, context.CancelFunc) {
	var (
		ctx    = c.context
		cancel context.CancelFunc
	)

	ctx = namespaces.WithNamespace(ctx, c.namespace)

	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	return ctx, cancel
}

func (c *ContainerdClient) getConnection() (*containerd.Client, error) {
	if c.client == nil {
		log.Debugf("Try to establish connection to containerd %s", c.address)
		client, err := containerd.New(c.address, containerd.WithDefaultNamespace(c.namespace))
		if err != nil {
			return client, errors.Wrapf(err, "Unable to create connection to containerd")
		}
		c.client = client
	}
	return c.client, nil
}

func (c *ContainerdClient) resetConnection() {
	c.client = nil
}

// GetContainers return all containerd containers
func (c *ContainerdClient) GetContainers() (containers []containerd.Container, err error) {
	ctx, cancel := c.getContext()
	defer cancel()

	client, err := c.getConnection()
	if err != nil {
		return containers, err
	}
	containers, err = client.Containers(ctx)
	if err != nil {
		c.resetConnection()
		return containers, errors.Wrap(err, "Error while getting list of containers")
	}
	return containers, nil
}

// CreateContainers create all given container definitions
func (c *ContainerdClient) CreateContainers(pod model.Pod, containers []model.Container) error {
	for _, container := range containers {
		if err := c.CreateContainer(pod, container); err != nil {
			return err
		}
	}
	return nil
}

// CreateContainer creates given container
func (c *ContainerdClient) CreateContainer(pod model.Pod, target model.Container) error {
	ctx, cancel := c.getContext()
	defer cancel()

	client, err := c.getConnection()
	if err != nil {
		return err
	}

	image, err := c.EnsureImagePulled(target.Image)
	if err != nil {
		return err
	}

	spec, err := containerd.GenerateSpec(containerd.WithImageConfig(ctx, image))
	if err != nil {
		return err
	}

	log.Debugf("Create new container from image %s...", image.Name())
	container, err := client.NewContainer(ctx, target.BuildID(pod.GetName()),
		containerd.WithSpec(spec),
		containerd.WithNewSnapshotView(target.BuildID(pod.GetName()), image),
		containerd.WithRuntime(fmt.Sprintf("%s.%s", plugin.RuntimePlugin, "linux")),
	)
	if err != nil {
		c.resetConnection()
		return errors.Wrapf(err, "Failed to create new container from image %s", image.Name())
	}

	log.Debugf("Create task in container: %s", container.ID())
	task, err := container.NewTask(ctx, containerd.NullIO)
	if err != nil {
		c.resetConnection()
		return errors.Wrapf(err, "Error while creating task for container [%s]", container.ID())
	}

	log.Debugln("Starting task...")
	err = task.Start(ctx)
	if err != nil {
		c.resetConnection()
		return errors.Wrapf(err, "Failed to start task in container", container.ID())
	}
	log.Debugf("Task started (pid %d)", task.Pid())
	return nil
}

// StopContainers stop all given containers
func (c *ContainerdClient) StopContainers(containers []containerd.Container) error {
	for _, container := range containers {
		err := c.StopContainer(container)
		if err != nil {
			return err
		}
	}
	return nil
}

// StopContainer stops given container
func (c *ContainerdClient) StopContainer(container containerd.Container) error {
	ctx, cancel := c.getContext()
	defer cancel()

	task, err := container.Task(ctx, nil)
	if err == nil {
		task.Delete(ctx, containerd.WithProcessKill)
	}
	if err := container.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
		c.resetConnection()
		return errors.Wrapf(err, "Failed to delete container %s", container.ID())
	}
	return nil
}

// EnsureImagePulled pulls the image reference to ensure image is fetched
func (c *ContainerdClient) EnsureImagePulled(ref string) (image containerd.Image, err error) {
	ctx, cancel := c.getContext()
	defer cancel()

	client, err := c.getConnection()
	if err != nil {
		return image, err
	}

	image, err = client.Pull(ctx, ref)
	if err != nil {
		c.resetConnection()
		return image, errors.Wrapf(err, "Error pulling image [%s]", ref)
	}

	log.Debugf("Unpacking container image [%s]...", image.Target().Digest)
	err = image.Unpack(ctx, "")
	if err != nil {
		c.resetConnection()
		return image, errors.Wrapf(err, "Error while unpacking image [%s]", image.Target().Digest)
	}

	return image, nil
}
