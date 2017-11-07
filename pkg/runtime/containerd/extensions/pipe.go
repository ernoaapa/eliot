package extensions

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/typeurl"
	"github.com/gogo/protobuf/types"
)

var pipeSetExtensionName = "can.io.pipeset"

// PipeSet allows defining pipe from some source(s) to another container
type PipeSet struct {
	Stdout PipeFromStdout
}

// PipeFromStdout defines container stdout as source for the piping
type PipeFromStdout struct {
	Stdin PipeToStdin
}

// PipeToStdin defines container stdin as target for the piping
type PipeToStdin struct {
	Name string
}

// WithPipeExtension appends pipe extension data to the container object.
func WithPipeExtension(pipe PipeSet) containerd.NewContainerOpts {
	return func(ctx context.Context, client *containerd.Client, c *containers.Container) error {
		any, err := typeurl.MarshalAny(&pipe)
		if err != nil {
			return err
		}

		if c.Extensions == nil {
			c.Extensions = make(map[string]types.Any)
		}
		c.Extensions[pipeSetExtensionName] = *any
		return nil
	}
}

// GetPipeExtension returns PipeSet from container extensions or nil if not defined
func GetPipeExtension(container containers.Container) (*PipeSet, error) {
	extension, ok := container.Extensions[pipeSetExtensionName]
	if !ok {
		return nil, nil
	}

	decoded, err := typeurl.UnmarshalAny(&extension)
	if err != nil {
		return nil, err
	}

	io, ok := decoded.(*PipeSet)
	if !ok {
		return nil, fmt.Errorf("Failed to decode PipeSet from container [%s] extensions", container.ID)
	}

	return io, err
}
