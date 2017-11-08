package extensions

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/typeurl"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
)

var lifecycleExtensionName = "can.io.lifecycle"

// RestartPolicy describes policy for restarting th econtainer
type RestartPolicy int

const (
	// Always is default restart policy what means that every time when container stops, it should be restarted.
	Always = iota
	// OnFailure means that only if process fails (non zero exit code) the container should be restarted
	OnFailure
)

func (p RestartPolicy) String() string {
	switch p {
	case Always:
		return "always"
	case OnFailure:
		return "onfailure"
	default:
		return "unknown"
	}
}

// ContainerLifecycle contains all lifecycle related information like restart counter and restart policy.
type ContainerLifecycle struct {
	// StartCount gets incremented on every time when container get started
	// If value is zero, assumed that it's not yet created
	StartCount    int
	RestartPolicy RestartPolicy
}

// WithLifecycleExtension is containerd.NewContainerOpts implementation what add lifecycle extension data to the container object.
func WithLifecycleExtension(ctx context.Context, client *containerd.Client, c *containers.Container) error {
	return updateLifecycleExtension(c, ContainerLifecycle{})
}

func updateLifecycleExtension(c *containers.Container, lifecycle ContainerLifecycle) error {
	any, err := typeurl.MarshalAny(&lifecycle)
	if err != nil {
		return err
	}
	extensions := c.Extensions
	if extensions == nil {
		extensions = make(map[string]types.Any)
	}
	extensions[lifecycleExtensionName] = *any

	c.Extensions = extensions
	return nil
}

// IncrementRestart is containerd.UpdateContainerOpts implementation what increments restart counter
func IncrementRestart(ctx context.Context, client *containerd.Client, c *containers.Container) error {
	lifecycle, err := GetLifecycleExtension(*c)
	if err != nil {
		return errors.Wrapf(err, "Cannot increment container restart counter")
	}
	lifecycle.StartCount++

	return updateLifecycleExtension(c, lifecycle)
}

// GetLifecycleExtension returns ContainerLifecycle from container extensions or nil if not defined
func GetLifecycleExtension(c containers.Container) (ContainerLifecycle, error) {
	extension, ok := c.Extensions[lifecycleExtensionName]
	if !ok {
		return ContainerLifecycle{}, fmt.Errorf("ContainerLifecycle extension not found in container [%s]", c.ID)
	}

	decoded, err := typeurl.UnmarshalAny(&extension)
	if err != nil {
		return ContainerLifecycle{}, errors.Wrapf(err, "Error while unmarshalling ContainerLifecycle of container [%s]", c.ID)
	}

	lifecycle, ok := decoded.(*ContainerLifecycle)
	if !ok {
		return ContainerLifecycle{}, fmt.Errorf("Failed to decode ContainerLifecycle from container [%s] extensions", c.ID)
	}

	return *lifecycle, err
}
