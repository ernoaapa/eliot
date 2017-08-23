package controller

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/plugin"
	"github.com/ernoaapa/layeryd/model"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

func createContainers(ctx context.Context, client *containerd.Client, pod model.Pod, containers []model.Container) error {
	for _, container := range containers {
		if err := createContainer(ctx, client, pod, container); err != nil {
			return err
		}
	}
	return nil
}

func createContainer(ctx context.Context, client *containerd.Client, pod model.Pod, target model.Container) error {
	image, err := ensureImagePulled(ctx, client, target.Image)
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
		return err
	}

	log.Debugf("Create task in container: %s", container.ID())
	task, err := container.NewTask(ctx, containerd.NullIO)
	if err != nil {
		return errors.Wrapf(err, "Error while creating task for container [%s]", container.ID())
	}

	log.Debugln("Starting task...")
	err = task.Start(ctx)
	if err != nil {
		return err
	}
	log.Debugf("Task started (pid %d)", task.Pid())
	return nil
}

func stopContainers(ctx context.Context, containers []containerd.Container) error {
	for _, container := range containers {
		err := stopContainer(ctx, container)
		if err != nil {
			return err
		}
	}
	return nil
}

func stopContainer(ctx context.Context, container containerd.Container) error {
	task, err := container.Task(ctx, nil)
	if err == nil {
		task.Delete(ctx, containerd.WithProcessKill)
	}
	if err := container.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
		return err
	}
	return nil
}

func ensureImagePulled(ctx context.Context, client *containerd.Client, ref string) (image containerd.Image, err error) {
	image, err = client.Pull(ctx, ref)
	if err != nil {
		return image, errors.Wrapf(err, "Error pulling image [%s]", ref)
	}

	log.Debugf("Unpacking container image %s...", image.Target().Digest)
	err = image.Unpack(ctx, "")
	if err != nil {
		return image, err
	}

	return image, nil
}
