package mapping

import (
	"net"

	core "github.com/ernoaapa/eliot/pkg/api/core"
	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	node "github.com/ernoaapa/eliot/pkg/api/services/node/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/model"
)

// MapInfoToAPIModel maps internal node info model to API model
func MapInfoToAPIModel(info *model.NodeInfo) *node.Info {
	return &node.Info{
		Uptime:      info.Uptime,
		Labels:      mapLabelsToAPIModel(info.Labels),
		Hostname:    info.Hostname,
		Addresses:   addressesToString(info.Addresses),
		GrpcPort:    int64(info.GrpcPort),
		MachineID:   info.MachineID,
		SystemUUID:  info.SystemUUID,
		BootID:      info.BootID,
		Arch:        info.Arch,
		Os:          info.OS,
		Version:     info.Version,
		Filesystems: mapFilesystemsToAPIModel(info.Filesystems),
	}
}

func mapLabelsToAPIModel(labels map[string]string) (result []*node.Label) {
	for key, value := range labels {
		result = append(result, &node.Label{Key: key, Value: value})
	}
	return result
}

func addressesToString(addresses []net.IP) (result []string) {
	for _, ip := range addresses {
		result = append(result, ip.String())
	}
	return result
}

// MapPodsToAPIModel maps list of internal pod models to API model
func MapPodsToAPIModel(pods []model.Pod) (result []*pods.Pod) {
	for _, pod := range pods {
		result = append(result, MapPodToAPIModel(pod))
	}
	return result
}

// MapPodToAPIModel maps internal Pod model to API model
func MapPodToAPIModel(pod model.Pod) *pods.Pod {
	return &pods.Pod{
		Metadata: &core.ResourceMetadata{
			Name:      pod.Metadata.Name,
			Namespace: pod.Metadata.Namespace,
		},
		Spec: &pods.PodSpec{
			Containers:    MapContainersToAPIModel(pod.Spec.Containers),
			HostNetwork:   pod.Spec.HostNetwork,
			HostPID:       pod.Spec.HostPID,
			RestartPolicy: pod.Spec.RestartPolicy,
		},
		Status: &pods.PodStatus{
			Hostname:          pod.Status.Hostname,
			ContainerStatuses: MapContainerStatusesToAPIModel(pod.Status.ContainerStatuses),
		},
	}
}

// MapContainersToAPIModel maps list of internal Container models to API model
func MapContainersToAPIModel(source []model.Container) (result []*containers.Container) {
	for _, container := range source {
		result = append(result, &containers.Container{
			Name:       container.Name,
			Image:      container.Image,
			WorkingDir: container.WorkingDir,
			Args:       container.Args,
			Env:        container.Env,
			Devices:    mapDevicesToAPIModel(container.Devices),
			Mounts:     mapMountsToAPIModel(container.Mounts),
			Pipe:       mapPipeToAPIModel(container.Pipe),
		})
	}
	return result
}

func mapDevicesToAPIModel(devices []model.Device) (result []*containers.Device) {
	for _, device := range devices {
		result = append(result, &containers.Device{
			DeviceType:  device.DeviceType,
			Minorid:     device.MinorId,
			Majorid:     device.MajorId,
		})
	}
	return result
}

func mapMountsToAPIModel(mounts []model.Mount) (result []*containers.Mount) {
	for _, mount := range mounts {
		result = append(result, &containers.Mount{
			Type:        mount.Type,
			Source:      mount.Source,
			Destination: mount.Destination,
			Options:     mount.Options,
		})
	}
	return result
}

func mapPipeToAPIModel(pipe *model.PipeSet) *containers.PipeSet {
	if pipe == nil {
		return nil
	}
	return &containers.PipeSet{
		Stdout: &containers.PipeFromStdout{
			Stdin: &containers.PipeToStdin{
				Name: pipe.Stdout.Stdin.Name,
			},
		},
	}
}

// MapContainerStatusesToAPIModel maps list of internal ContainerStatus models to API model
func MapContainerStatusesToAPIModel(statuses []model.ContainerStatus) (result []*containers.ContainerStatus) {
	for _, status := range statuses {
		result = append(result, &containers.ContainerStatus{
			ContainerID:  status.ContainerID,
			Name:         status.Name,
			Image:        status.Image,
			State:        status.State,
			RestartCount: int32(status.RestartCount),
		})
	}
	return result
}

func mapFilesystemsToAPIModel(disks []model.Filesystem) (result []*node.Filesystem) {
	for _, disk := range disks {
		result = append(result, &node.Filesystem{
			Filesystem: disk.Filesystem,
			TypeName:   disk.TypeName,
			MountDir:   disk.MountDir,
			Total:      disk.Total,
			Free:       disk.Free,
			Available:  disk.Available,
		})
	}
	return result
}
