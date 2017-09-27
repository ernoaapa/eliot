package stream

import (
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	log "github.com/sirupsen/logrus"
)

// Writer is io.Writer implementation what writes stdin/stdout bytes to RPC stream
type Writer struct {
	stream StdOutputStreamServer
	stderr bool
}

// StdOutputStreamServer interface for the endpoint what returns stream of log lines
type StdOutputStreamServer interface {
	Send(*pb.StdOutputStreamResponse) error
}

// NewWriter creates new Writer instance
func NewWriter(stream StdOutputStreamServer, stderr bool) *Writer {
	return &Writer{stream, stderr}
}

// Write writes bytes to given RPC stream
func (w *Writer) Write(p []byte) (n int, err error) {
	n = len(p)
	log.Debugf("Write %d bytes to stream", n)
	err = w.stream.Send(&pb.StdOutputStreamResponse{
		Line:   p,
		Stderr: w.stderr,
	})
	log.Debugf("After send, received err: %v", err)
	if err != nil {
		return n, err
	}
	return n, nil
}
