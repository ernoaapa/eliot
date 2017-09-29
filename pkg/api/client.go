package api

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
)

// Client connects to RPC server
type Client struct {
	namespace  string
	serverAddr string
}

// NewClient creates new RPC server client
func NewClient(namespace, serverAddr string) *Client {
	return &Client{
		namespace,
		serverAddr,
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
	resp, err := client.List(context.Background(), &pb.ListPodsRequest{
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
func (c *Client) CreatePod(pod *pb.Pod) (*pb.Pod, error) {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewPodsClient(conn)
	resp, err := client.Create(context.Background(), &pb.CreatePodRequest{
		Pod: pod,
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

	resp, err := client.Delete(context.Background(), &pb.DeletePodRequest{
		Namespace: pod.Metadata.Namespace,
		Name:      pod.Metadata.Name,
	})
	return resp.GetPod(), nil
}

// Attach calls server and fetches pod logs
func (c *Client) Attach(containerID string, stdin io.Reader, stdout, stderr io.Writer) (err error) {
	wg := &sync.WaitGroup{}
	md := metadata.Pairs(
		"namespace", c.namespace,
		"container", containerID,
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
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

	wg.Add(1)
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Debugln("Received EOF for log stream")
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
				break
			}
		}
		wg.Done()
	}()

	if stdin != nil {
		wg.Add(1)
		go func() {
			for {
				target := []byte{1}
				_, err := stdin.Read(target)
				if err != nil {
					break
				}

				err = stream.Send(&pb.StdinStreamRequest{
					Input: target,
				})
				if err != nil {
					break
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	return err
}
