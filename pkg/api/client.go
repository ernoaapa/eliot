package api

import (
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
