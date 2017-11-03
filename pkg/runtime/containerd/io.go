package containerd

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/containerd/containerd"
	"github.com/containerd/fifo"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/pkg/errors"
)

// NewFifos returns a new set of fifos for the task
func NewFifos(id string) (*containerd.FIFOSet, error) {
	root := "/run/containerd/fifo"
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, err
	}
	dir, err := ioutil.TempDir(root, "")
	if err != nil {
		return nil, err
	}
	return &containerd.FIFOSet{
		Dir: dir,
		In:  filepath.Join(dir, id+"-stdin"),
		Out: filepath.Join(dir, id+"-stdout"),
		Err: filepath.Join(dir, id+"-stderr"),
	}, nil
}

func ensureDirsExist(set model.IOSet) error {
	if err := os.MkdirAll(filepath.Dir(set.In), 0700); err != nil {
		return errors.Wrapf(err, "Failed to create IOSet.In [%s]", set.In)
	}
	if err := os.MkdirAll(filepath.Dir(set.Out), 0700); err != nil {
		return errors.Wrapf(err, "Failed to create IOSet.Out [%s]", set.Out)
	}
	if err := os.MkdirAll(filepath.Dir(set.Err), 0700); err != nil {
		return errors.Wrapf(err, "Failed to create IOSet.Err [%s]", set.Err)
	}
	return nil
}

// NewDirectIO returns an IO implementation that exposes the pipes directly
func NewDirectIO(ctx context.Context, ioSet model.IOSet, terminal bool) (f *DirectIO, err error) {
	err = ensureDirsExist(ioSet)
	if err != nil {
		return nil, err
	}

	f = &DirectIO{
		ioSet:    ioSet,
		terminal: terminal,
	}

	if f.Stdin, err = fifo.OpenFifo(ctx, ioSet.In, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil {
		return nil, errors.Wrapf(err, "Failed to open in FIFO [%s]", ioSet.In)
	}
	if f.Stdout, err = fifo.OpenFifo(ctx, ioSet.Out, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil {
		f.Stdin.Close()
		return nil, errors.Wrapf(err, "Failed to open out FIFO [%s]", ioSet.Out)
	}
	if f.Stderr, err = fifo.OpenFifo(ctx, ioSet.Err, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil {
		f.Stdin.Close()
		f.Stdout.Close()
		return nil, errors.Wrapf(err, "Failed to open err FIFO [%s]", ioSet.Err)
	}
	return f, nil
}

// DirectIO allows task IO to be handled externally by the caller
type DirectIO struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser

	ioSet    model.IOSet
	terminal bool
}

// IOCreate returns IO avaliable for use with task creation
func (f *DirectIO) IOCreate(id string) (containerd.IO, error) {
	return f, nil
}

// Config returns the IOConfig
func (f *DirectIO) Config() containerd.IOConfig {
	return containerd.IOConfig{
		Terminal: f.terminal,
		Stdin:    f.ioSet.In,
		Stdout:   f.ioSet.Out,
		Stderr:   f.ioSet.Err,
	}
}

// Cancel stops any IO copy operations
//
// Not applicable for DirectIO
func (f *DirectIO) Cancel() {
	// nothing to cancel as all operations are handled externally
}

// Wait on any IO copy operations
//
// Not applicable for DirectIO
func (f *DirectIO) Wait() {
	// nothing to wait on as all operations are handled externally
}

// Close closes all open fds
func (f *DirectIO) Close() error {
	err := f.Stdin.Close()
	if err2 := f.Stdout.Close(); err == nil {
		err = err2
	}
	if err2 := f.Stderr.Close(); err == nil {
		err = err2
	}
	return err
}

// Delete removes the underlying directory containing fifos
func (f *DirectIO) Delete() error {
	return nil
}
