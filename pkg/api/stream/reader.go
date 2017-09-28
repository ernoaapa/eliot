package stream

import (
	"bytes"

	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	log "github.com/sirupsen/logrus"
)

// Reader is io.Reader implementation what reads bytes from RPC stream
type Reader struct {
	buffer bytes.Buffer
	stream StdinStreamServer
}

// StdinStreamServer interface for the endpoint what takes stdin stream in
type StdinStreamServer interface {
	Recv() (*pb.StdinStreamRequest, error)
}

// NewReader creates new Reader instance
func NewReader(stream StdinStreamServer) *Reader {
	return &Reader{stream: stream}
}

// Write writes bytes to given RPC stream
func (w *Reader) Read(p []byte) (n int, err error) {
	if w.buffer.Len() == 0 {
		log.Debugf("Nothing in the buffer, start waiting for stream input")
		req, err := w.stream.Recv()
		if err != nil {
			return 0, err
		}
		log.Debugf("Received input %s", req.GetInput())
		w.buffer.Write(req.GetInput())
	}
	return w.buffer.Read(p)
}
