package containerd

import (
	"context"
	"io"
	"syscall"

	"github.com/containerd/containerd"
	"github.com/containerd/fifo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// DetachedFifoIO creates a new fifo files but don't attach it to os.Stdin/os.Stdout.
// Later on we can attach to it and read the log lines.
func DetachedFifoIO(id string) (_ containerd.IO, err error) {
	paths, err := createNewFifos(id)
	if err != nil {
		return nil, err
	}
	log.Printf("%v", paths)
	return DetachedIO{
		config: containerd.IOConfig{
			Terminal: false,
			Stdout:   paths.Out,
			Stderr:   paths.Err,
			Stdin:    paths.In,
		},
	}, nil
}

func createNewFifos(id string) (*containerd.FIFOSet, error) {
	paths, err := containerd.NewFifos(id)
	if err != nil {
		return paths, err
	}
	return paths, touchFifos([]string{
		paths.Out,
		paths.Err,
		paths.In,
	})
}

func touchFifos(paths []string) error {
	var (
		f   io.ReadWriteCloser
		err error
		ctx = context.Background()
	)
	for _, path := range paths {
		if f, err = fifo.OpenFifo(ctx, path, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil {
			return errors.Wrapf(err, "Error while trying to touch (=create) new fifo empty file [%s]", path)
		}

		if err := f.Close(); err != nil {
			return errors.Wrap(err, "Error while trying to close touched fifo file")
		}
		log.Debugf("Created fifo %s", path)
	}
	return nil
}

// DetachedIO is containerd.IO implementation what
// is collecting output to fifo files, but is "detached"
// i.e. you cannot cancel copying etc.
type DetachedIO struct {
	config containerd.IOConfig
}

// Config returns the IO configuration.
func (d DetachedIO) Config() containerd.IOConfig {
	return d.config
}

// Cancel aborts all current io operations
func (DetachedIO) Cancel() {
}

// Wait blocks until all io copy operations have completed
func (DetachedIO) Wait() {
}

// Close cleans up all open io resources
func (DetachedIO) Close() error {
	return nil
}
