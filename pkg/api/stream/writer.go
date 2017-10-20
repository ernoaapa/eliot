package stream

import (
	containers "github.com/ernoaapa/can/pkg/api/services/containers/v1"
)

// Writer is io.Writer implementation what writes stdout/stderr bytes to RPC stream
type Writer struct {
	stream StdoutStreamServer
	stderr bool
}

// StdoutStreamServer interface for the endpoint what returns stream of log lines
type StdoutStreamServer interface {
	Send(*containers.StdoutStreamResponse) error
}

// StdoutStreamClient interface for the client what reads stream of log lines
type StdoutStreamClient interface {
	Recv() (*containers.StdoutStreamResponse, error)
	CloseSend() error
}

// NewWriter creates new Writer instance
func NewWriter(stream StdoutStreamServer, stderr bool) *Writer {
	return &Writer{stream, stderr}
}

// Write writes bytes to given RPC stream
func (w *Writer) Write(p []byte) (n int, err error) {
	n = len(p)
	err = w.stream.Send(&containers.StdoutStreamResponse{
		Output: p[:],
		Stderr: w.stderr,
	})
	if err != nil {
		return n, err
	}
	return n, nil
}
