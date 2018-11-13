package mapping

import (
	"encoding/json"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/ernoaapa/eliot/pkg/runtime/containerd/extensions"
)

// GetPodName resolves pod name where the container belongs
func GetPodName(container containers.Container) string {
	labels := ContainerLabels(container.Labels)
	podName := labels.getPodName()
	if podName == "" {
		// container is not eliot managed container so add it under 'system' pod in namespace 'default'
		podName = "system"
	}
	return podName
}

// InitialisePodModel creates new Pod struct with name and namespace metadata
func InitialisePodModel(container containers.Container, namespace, name, hostname string) model.Pod {
	return model.Pod{
		Metadata: model.NewMetadata(namespace, name),
		Spec: model.PodSpec{
			Containers:    []model.Container{},
			HostNetwork:   !haveNamespace(container, specs.NetworkNamespace),
			HostPID:       !haveNamespace(container, specs.PIDNamespace),
			RestartPolicy: getRestartPolicy(container),
		},
		Status: model.PodStatus{
			Hostname:          hostname,
			ContainerStatuses: []model.ContainerStatus{},
		},
	}
}

// MapContainersToInternalModel maps containerd models to internal model
func MapContainersToInternalModel(containers []containers.Container) (result []model.Container) {
	for _, container := range containers {
		result = append(result, MapContainerToInternalModel(container))
	}
	return result
}

// MapContainerToInternalModel maps containerd model to internal model
func MapContainerToInternalModel(container containers.Container) model.Container {
	labels := ContainerLabels(container.Labels)
	return model.Container{
		Name:       labels.getContainerName(),
		Image:      container.Image,
		Tty:        RequireTty(container),
		Args:       processArgs(container),
		Env:        processEnv(container),
		WorkingDir: processWorkingDir(container),
		Pipe:       mapPipeToInternalModel(container),
		Mounts:     mapMountsToInternalModel(container),
		Devices:	mapDevicesToInternalModel(container),
	}
}

// RequireTty find out is the container configured to create TTY
func RequireTty(container containers.Container) bool {
	spec, err := getSpec(container)
	if err != nil {
		log.Fatalf("Cannot read container spec to resolve process TTY value: %s", err)
		return false
	}
	return spec.Process.Terminal
}

// Spec returns the current OCI specification for the container
func getSpec(container containers.Container) (*specs.Spec, error) {
	var s specs.Spec
	if err := json.Unmarshal(container.Spec.Value, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func mapPipeToInternalModel(container containers.Container) *model.PipeSet {
	pipe, err := extensions.GetPipeExtension(container)
	if err != nil {
		log.Errorf("Failed to read Pipe extension from container [%s]: %s", container.ID, err)
	}
	if pipe == nil {
		return nil
	}

	return &model.PipeSet{
		Stdout: &model.PipeFromStdout{
			Stdin: &model.PipeToStdin{
				Name: pipe.Stdout.Stdin.Name,
			},
		},
	}
}

func processArgs(container containers.Container) []string {
	spec, err := getSpec(container)
	if err != nil {
		log.Fatalf("Cannot read container spec to resolve process args: %s", err)
		return nil
	}

	return spec.Process.Args
}

func processEnv(container containers.Container) []string {
	spec, err := getSpec(container)
	if err != nil {
		log.Fatalf("Cannot read container spec to resolve process environment variables: %s", err)
		return nil
	}

	return spec.Process.Env
}

func processWorkingDir(container containers.Container) string {
	spec, err := getSpec(container)
	if err != nil {
		log.Fatalf("Cannot read container spec to resolve process current working directory: %s", err)
		return ""
	}

	return spec.Process.Cwd
}

func mapMountsToInternalModel(container containers.Container) (result []model.Mount) {
	spec, err := getSpec(container)
	if err != nil {
		log.Fatalf("Cannot read container spec to resolve container mounts: %s", err)
		return result
	}

	for _, mount := range spec.Mounts {
		result = append(result, model.Mount{
			Type:        mount.Type,
			Source:      mount.Source,
			Destination: mount.Destination,
			Options:     mount.Options,
		})
	}
	return result
}


func mapDevicesToInternalModel(container containers.Container) (result []model.Device) {
	spec, err := getSpec(container)
	if err != nil {
		log.Fatalf("Cannot read container spec to resolve container devices: %s", err)
		return result
	}

	ptrint32 := func(i *int64) uint32 {
		if (i == nil) {
			return uint32(0);
		}
		return uint32(*i)
	}

	for _, device := range spec.Linux.Resources.Devices {
		//log.Debugf("%p", device.Major)
		result = append(result, model.Device{
			DeviceType:  device.Type,
			MajorId:     ptrint32(device.Major),
			MinorId:     ptrint32(device.Minor),
		})
	}
	return result
}

// MapContainerStatusToInternalModel maps containerd model to internal container status model
func MapContainerStatusToInternalModel(container containers.Container, status containerd.Status) model.ContainerStatus {
	labels := ContainerLabels(container.Labels)
	return model.ContainerStatus{
		ContainerID:  container.ID,
		Name:         labels.getContainerName(),
		Image:        container.Image,
		State:        mapContainerStatus(status),
		RestartCount: getRestartCount(container),
	}
}

func getRestartCount(container containers.Container) int {
	lifecycle, err := extensions.GetLifecycleExtension(container)
	if err != nil && !extensions.IsNotFound(err) {
		log.Warnf("Error while resolving container restart count, fallback to zero: %s", err)
	}

	if lifecycle.StartCount <= 1 {
		return 0
	}
	return lifecycle.StartCount - 1
}

func haveNamespace(container containers.Container, namespace specs.LinuxNamespaceType) bool {
	spec, err := getSpec(container)
	if err != nil {
		log.Fatalf("Cannot read container spec to resolve namespace %s: %s", namespace, err)
		return false
	}

	for _, ns := range spec.Linux.Namespaces {
		if ns.Type == namespace {
			return true
		}
	}
	return false
}

func getRestartPolicy(container containers.Container) string {
	lifecycle, err := extensions.GetLifecycleExtension(container)
	if err != nil && !extensions.IsNotFound(err) {
		log.Warnf("Error while resolving container restart policy, fallback to default: %s", err)
	}

	return lifecycle.RestartPolicy.String()
}

func mapContainerStatus(status containerd.Status) string {
	if status.Status == "" {
		return string(containerd.Unknown)
	}
	return string(status.Status)
}
