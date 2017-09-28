package stream

import (
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
)

// Writer is io.Writer implementation what writes stdout/stderr bytes to RPC stream
type Writer struct {
	stream StdoutStreamServer
	stderr bool
}

// StdoutStreamServer interface for the endpoint what returns stream of log lines
type StdoutStreamServer interface {
	Send(*pb.StdoutStreamResponse) error
}

// NewWriter creates new Writer instance
func NewWriter(stream StdoutStreamServer, stderr bool) *Writer {
	return &Writer{stream, stderr}
}

// Write writes bytes to given RPC stream
func (w *Writer) Write(p []byte) (n int, err error) {
	n = len(p)
	err = w.stream.Send(&pb.StdoutStreamResponse{
		Output: p[:],
		Stderr: w.stderr,
	})
	if err != nil {
		return n, err
	}
	return n, nil
}
