package stream

import (
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	log "github.com/sirupsen/logrus"
)

// LogsWriter is io.Writer implementation what writes bytes to stream
type LogsWriter struct {
	stream  pb.Pods_LogsServer
	logType pb.GetLogsResponse_Type
}

// NewLogsWriter creates new LogsWriter instance
func NewLogsWriter(stream pb.Pods_LogsServer, logType pb.GetLogsResponse_Type) *LogsWriter {
	return &LogsWriter{
		stream, logType,
	}
}

func (w *LogsWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	log.Debugf("Write %d bytes to stream", n)
	err = w.stream.Send(&pb.GetLogsResponse{
		Line: p,
		Type: w.logType,
	})
	log.Debugf("After send, received err: %v", err)
	if err != nil {
		return n, err
	}
	return n, nil
}
