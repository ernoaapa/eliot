package api

import (
	"bytes"
	"fmt"
	"io"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ernoaapa/can/pkg/api/mapping"
	containers "github.com/ernoaapa/can/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/api/stream"
	"github.com/ernoaapa/can/pkg/progress"
)

// Client connects to RPC server
type Client struct {
	Namespace string
	Endpoint  string
	ctx       context.Context
}

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

// NewClient creates new RPC server client
func NewClient(namespace, serverAddr string) *Client {
	return &Client{
		namespace,
		serverAddr,
		context.Background(),
	}
}

// GetPods calls server and fetches all pods information
func (c *Client) GetPods() ([]*pods.Pod, error) {
	conn, err := grpc.Dial(c.Endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pods.NewPodsClient(conn)
	resp, err := client.List(c.ctx, &pods.ListPodsRequest{
		Namespace: c.Namespace,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetPods(), nil
}

// GetPod return Pod by name
func (c *Client) GetPod(podName string) (*pods.Pod, error) {
	pods, err := c.GetPods()
	if err != nil {
		return nil, err
	}

	for _, pod := range pods {
		if pod.Metadata.Name == podName {
			return pod, nil
		}
	}
	return nil, fmt.Errorf("No pod found with name [%s]", podName)
}

// PodOpts adds more information to the Pod going to be created
type PodOpts func(pod *pods.Pod) error

// CreatePod creates new pod to the target server
func (c *Client) CreatePod(pod *pods.Pod, opts ...PodOpts) error {
	for _, o := range opts {
		err := o(pod)
		if err != nil {
			return err
		}
	}

	conn, err := grpc.Dial(c.Endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pods.NewPodsClient(conn)
	stream, err := client.Create(c.ctx, &pods.CreatePodRequest{
		Pod: pod,
	})
	if err != nil {
		return err
	}

	progress := progress.NewRenderer()
	defer progress.Stop()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			err = stream.CloseSend()
			progress.Done()
			return err
		}
		if err != nil {
			return err
		}

		progress.Update(mapping.MapAPIModelToImageFetchProgress(resp.Images))
	}
}

// StartPod creates new pod to the target server
func (c *Client) StartPod(name string) (*pods.Pod, error) {
	conn, err := grpc.Dial(c.Endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pods.NewPodsClient(conn)
	resp, err := client.Start(c.ctx, &pods.StartPodRequest{
		Namespace: c.Namespace,
		Name:      name,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetPod(), nil
}

// DeletePod creates new pod to the target server
func (c *Client) DeletePod(pod *pods.Pod) (*pods.Pod, error) {
	conn, err := grpc.Dial(c.Endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pods.NewPodsClient(conn)

	resp, err := client.Delete(c.ctx, &pods.DeletePodRequest{
		Namespace: pod.Metadata.Namespace,
		Name:      pod.Metadata.Name,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetPod(), nil
}

type AttachHooks func(client *Client, done <-chan struct{})

// Attach calls server and fetches pod logs
func (c *Client) Attach(containerID string, attachIO AttachIO, hooks ...AttachHooks) (err error) {
	done := make(chan struct{})
	errc := make(chan error)

	md := metadata.Pairs(
		"namespace", c.Namespace,
		"container", containerID,
	)
	ctx, cancel := context.WithCancel(metadata.NewOutgoingContext(c.ctx, md))
	defer cancel()

	conn, err := grpc.Dial(c.Endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := containers.NewContainersClient(conn)
	log.Debugf("Open stream connection to server to get logs")
	stream, err := client.Attach(ctx)
	if err != nil {
		return err
	}

	go func() {
		errc <- PipeStdout(stream, attachIO.Stdout, attachIO.Stderr)
	}()

	if attachIO.Stdin != nil {
		go func() {
			errc <- PipeStdin(stream, attachIO.Stdin)
		}()
	}

	for _, hook := range hooks {
		go hook(c, done)
	}

	for {
		err := <-errc
		close(done)
		return err
	}
}

// PipeStdout reads stdout from grpc stream and writes it to stdout/stderr
func PipeStdout(s stream.StdoutStreamClient, stdout, stderr io.Writer) error {
	for {
		resp, err := s.Recv()
		if err == io.EOF {
			err = s.CloseSend()
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
func PipeStdin(s stream.StdinStreamClient, stdin io.Reader) error {
	for {
		buf := make([]byte, 1024)
		n, err := stdin.Read(buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "Error while reading stdin to buffer")
		}

		if err := s.Send(&containers.StdinStreamRequest{Input: buf[:n]}); err != nil {
			return errors.Wrapf(err, "Sending to stream returned error")
		}
	}
}

// Signal sends kill signal to container process
func (c *Client) Signal(containerID string, signal syscall.Signal) (err error) {
	conn, err := grpc.Dial(c.Endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := containers.NewContainersClient(conn)

	_, err = client.Signal(c.ctx, &containers.SignalRequest{
		Namespace:   c.Namespace,
		ContainerID: containerID,
		Signal:      int32(signal),
	})

	return err
}
