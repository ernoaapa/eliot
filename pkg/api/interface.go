package api

import (
	"io"

	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/config"
)

// PodOpts adds more information to the Pod going to be created
type PodOpts func(pod *pods.Pod) error

// AttachHooks is additional process what runs when is attached to container
type AttachHooks func(endpoint config.Endpoint, done <-chan struct{})

// AttachIO wraps stdin/stdout for attach
type AttachIO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewAttachIO is wrapper for stdin, stdout and stderr
func NewAttachIO(stdin io.Reader, stdout, stderr io.Writer) AttachIO {
	return AttachIO{stdin, stdout, stderr}
}
