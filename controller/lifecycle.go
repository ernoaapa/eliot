package controller

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/plugin"
	"github.com/ernoaapa/layeryd/model"

	log "github.com/sirupsen/logrus"
)

func createContainers(ctx context.Context, client *containerd.Client, containers []model.Container) error {
	for _, container := range containers {
		if err := createContainer(ctx, client, container); err != nil {
			return err
		}
	}
	return nil
}

func createContainer(ctx context.Context, client *containerd.Client, target model.Container) error {
	image, err := ensureImagePulled(ctx, client, target.Image)
	if err != nil {
		log.Warnf("Error getting image: %v", err)
		return err
	}

	spec, err := containerd.GenerateSpec(containerd.WithImageConfig(ctx, image))
	if err != nil {
		return err
	}

	log.Infof("Create new container from image %s...", image.Name())
	container, err := client.NewContainer(ctx, target.Name,
		containerd.WithSpec(spec),
		containerd.WithNewSnapshotView(target.Name, image),
		containerd.WithRuntime(fmt.Sprintf("%s.%s", plugin.RuntimePlugin, "linux")),
	)
	if err != nil {
		return err
	}

	log.Printf("Create task in container: %s", container.ID())
	task, err := container.NewTask(ctx, containerd.NullIO)
	if err != nil {
		log.Infof("Error in newtask: %s", err)
		return err
	}

	log.Println("Starting task...")
	err = task.Start(ctx)
	if err != nil {
		return err
	}
	log.Printf("Task started (%d)", task.Pid())
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
		log.Warnf("Error from pull: %v", err)
		return image, err
	}

	log.Infof("unpacking %s...", image.Target().Digest)
	err = image.Unpack(ctx, "")
	if err != nil {
		return image, err
	}

	return image, err
}
