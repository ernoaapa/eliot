package main

import (
	"fmt"
	"os"

	"github.com/ernoaapa/can/cmd"
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: "List pods",
	Action: func(clicontext *cli.Context) error {
		if clicontext.NArg() < 2 {
			return fmt.Errorf("You must give two parameters, NAME and IMAGE")
		}
		name := clicontext.Args().Get(0)
		image := clicontext.Args().Get(1)

		if name == "" {
			return fmt.Errorf("NAME argument cannot be empty")
		}
		if image == "" {
			return fmt.Errorf("IMAGE argument cannot be empty")
		}

		config := cmd.GetConfig(clicontext)
		client := cmd.GetClient(clicontext)

		pod := &pb.Pod{
			Metadata: &pb.ResourceMetadata{
				Name:      name,
				Namespace: config.GetCurrentContext().Namespace,
			},
			Spec: &pb.PodSpec{
				Containers: []*pb.Container{
					&pb.Container{
						Name:  name,
						Image: image,
					},
				},
			},
		}

		result, err := client.CreatePod(pod)
		if err != nil {
			return err
		}

		writer := printers.GetNewTabWriter(os.Stdout)
		printer := cmd.GetPrinter(clicontext)
		return printer.PrintPodDetails(result, writer)
	},
}
