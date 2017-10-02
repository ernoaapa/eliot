package mapping

import (
	"github.com/ernoaapa/can/pkg/model"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// MapMountToContainerdModel maps model.Mount to containerd spec struct
func MapMountToContainerdModel(mount model.Mount) specs.Mount {
	return specs.Mount{
		Type:        mount.Type,
		Source:      mount.Source,
		Destination: mount.Destination,
		Options:     mount.Options,
	}
}
