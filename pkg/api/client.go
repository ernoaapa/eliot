package api

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

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
func (c *Client) Attach(containerID string, stdout, stderr io.Writer) error {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewPodsClient(conn)
	log.Debugf("Open stream connection to server to get logs")
	stream, err := client.Attach(context.Background(), &pb.AttachRequest{
		Namespace:   c.namespace,
		ContainerID: containerID,
	})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Debugln("Received EOF for log stream")
			err := stream.CloseSend()
			return err
		}
		if err != nil {
			return err
		}

		if resp.Stderr {
			fmt.Fprint(stderr, string(resp.Line))
		} else {
			fmt.Fprint(stdout, string(resp.Line))
		}
	}
}
