package mapping

import (
	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/ernoaapa/eliot/pkg/runtime/containerd/extensions"
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

// MapPipeToContainerdModel maps model.PipeSet to containerd extension PipeSet
func MapPipeToContainerdModel(pipe model.PipeSet) extensions.PipeSet {
	return extensions.PipeSet{
		Stdout: extensions.PipeFromStdout{
			Stdin: extensions.PipeToStdin{
				Name: pipe.Stdout.Stdin.Name,
			},
		},
	}
}
