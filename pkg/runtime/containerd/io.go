package containerd

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/containerd/containerd"
	"github.com/containerd/fifo"
	"github.com/pkg/errors"
)

func ensureDirsExist(paths ...string) error {
	for _, path := range paths {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return errors.Wrapf(err, "Failed to ensure file [%s] parent directory exist", path)
		}
	}

	return nil
}

// NewDirectIO returns an IO implementation that exposes the pipes directly
func NewDirectIO(ctx context.Context, stdin, stdout, stderr string, terminal bool) (f *DirectIO, err error) {
	if err = ensureDirsExist(stdin, stdout, stderr); err != nil {
		return nil, err
	}

	f = &DirectIO{
		stdin:    stdin,
		stdout:   stdout,
		stderr:   stderr,
		terminal: terminal,
	}

	if f.Stdin, err = fifo.OpenFifo(ctx, stdin, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil {
		return nil, errors.Wrapf(err, "Failed to open in FIFO [%s]", stdin)
	}
	if f.Stdout, err = fifo.OpenFifo(ctx, stdout, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil {
		f.Stdin.Close()
		return nil, errors.Wrapf(err, "Failed to open out FIFO [%s]", stdout)
	}
	if f.Stderr, err = fifo.OpenFifo(ctx, stderr, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil {
		f.Stdin.Close()
		f.Stdout.Close()
		return nil, errors.Wrapf(err, "Failed to open err FIFO [%s]", stderr)
	}
	return f, nil
}

// DirectIO allows task IO to be handled externally by the caller
type DirectIO struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser

	stdin    string
	stdout   string
	stderr   string
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
		Stdin:    f.stdin,
		Stdout:   f.stdout,
		Stderr:   f.stderr,
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
