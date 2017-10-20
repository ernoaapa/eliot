package stream

import (
	"bytes"
	"io"

	containers "github.com/ernoaapa/can/pkg/api/services/containers/v1"
	"github.com/pkg/errors"
)

// PipeStdout reads stdout from grpc stream and writes it to stdout/stderr
func PipeStdout(stream StdoutStreamClient, stdout, stderr io.Writer) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			err = stream.CloseSend()
			if err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "Received error while reading attach stream")
		}

		target := stdout
		if resp.Stderr {
			target = stderr
		}

		_, err = io.Copy(target, bytes.NewReader(resp.Output))
		if err != nil {
			return errors.Wrapf(err, "Error while copying data")
		}
	}
}

// PipeStdin reads input from Stdin and writes it to the grpc stream
func PipeStdin(stream StdinStreamClient, stdin io.Reader) error {
	for {
		buf := make([]byte, 1024)
		n, err := stdin.Read(buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "Error while reading stdin to buffer")
		}

		if err := stream.Send(&containers.StdinStreamRequest{Input: buf[:n]}); err != nil {
			return errors.Wrapf(err, "Sending to stream returned error")
		}
	}
}
