package api

import (
	"bytes"
	"fmt"
	"io"
	"syscall"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ernoaapa/can/pkg/api/mapping"
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/progress"
)

// Client connects to RPC server
type Client struct {
	namespace  string
	serverAddr string
	ctx        context.Context
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
func (c *Client) GetPods() ([]*pb.Pod, error) {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewPodsClient(conn)
	resp, err := client.List(c.ctx, &pb.ListPodsRequest{
		Namespace: c.namespace,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetPods(), nil
}

// GetPod return Pod by name
func (c *Client) GetPod(podName string) (*pb.Pod, error) {
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

// CreatePod creates new pod to the target server
func (c *Client) CreatePod(pod *pb.Pod) error {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewPodsClient(conn)
	stream, err := client.Create(c.ctx, &pb.CreatePodRequest{
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
			return nil
		}
		if err != nil {
			return err
		}

		progress.Update(mapping.MapAPIModelToImageFetchProgress(resp.Images))
	}
}

// StartPod creates new pod to the target server
func (c *Client) StartPod(name string) (*pb.Pod, error) {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewPodsClient(conn)
	resp, err := client.Start(c.ctx, &pb.StartPodRequest{
		Namespace: c.namespace,
		Name:      name,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetPod(), nil
}

// DeletePod creates new pod to the target server
func (c *Client) DeletePod(pod *pb.Pod) (*pb.Pod, error) {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewPodsClient(conn)

	resp, err := client.Delete(c.ctx, &pb.DeletePodRequest{
		Namespace: pod.Metadata.Namespace,
		Name:      pod.Metadata.Name,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetPod(), nil
}

// Attach calls server and fetches pod logs
func (c *Client) Attach(containerID string, stdin io.Reader, stdout, stderr io.Writer) (err error) {
	done := make(chan struct{})
	md := metadata.Pairs(
		"namespace", c.namespace,
		"container", containerID,
	)
	ctx := metadata.NewOutgoingContext(c.ctx, md)
	conn, err := grpc.Dial(c.serverAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewPodsClient(conn)
	log.Debugf("Open stream connection to server to get logs")
	stream, err := client.Attach(ctx)
	if err != nil {
		return err
	}

	go func() {
		defer close(done)
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				err = stream.CloseSend()
				break
			}
			if err != nil {
				log.Debugf("Received error: %s", err)
				break
			}

			target := stdout
			if resp.Stderr {
				target = stderr
			}

			_, err = io.Copy(target, bytes.NewReader(resp.Output))
			if err != nil {
				log.Debugf("Error while copying data: %s", err)
				break
			}
		}
	}()

	if stdin != nil {
		go func() {
			defer close(done)

			for {
				buf := make([]byte, 1024)
				n, err := stdin.Read(buf)
				if err == io.EOF {
					// nothing else to pipe, kill this goroutine
					break
				}
				if err != nil {
					log.Debugf("Error while reading stdin to buffer: %s", err)
					break
				}

				err = stream.Send(&pb.StdinStreamRequest{
					Input: buf[:n],
				})
				if err != nil {
					log.Debugf("Sending to stream returned error %s", err)
					break
				}
			}
		}()
	}

	<-done
	return err
}

// Signal sends kill signal to container process
func (c *Client) Signal(containerID string, signal syscall.Signal) (err error) {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewPodsClient(conn)

	_, err = client.Signal(c.ctx, &pb.SignalRequest{
		Namespace:   c.namespace,
		ContainerID: containerID,
		Signal:      int32(signal),
	})

	return err
}
