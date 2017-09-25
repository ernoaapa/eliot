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

// GetLogs calls server and fetches pod logs
func (c *Client) GetLogs(containerID string, stdout, stderr io.Writer) error {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewPodsClient(conn)
	log.Debugf("Open stream connection to server to get logs")
	stream, err := client.Logs(context.Background(), &pb.GetLogsRequest{
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

		switch t := resp.Type; t {
		case pb.GetLogsResponse_STDOUT:
			fmt.Fprint(stdout, string(resp.Line))
		case pb.GetLogsResponse_STDERR:
			fmt.Fprint(stderr, string(resp.Line))
		default:
			return fmt.Errorf("Received unknown GetLogsResponse.Type, [%s]. Client version not matching with server?", t)
		}
	}
}
