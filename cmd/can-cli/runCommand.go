package main

import (
	"fmt"
	"os"

	"github.com/ernoaapa/can/cmd"
	pb "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/urfave/cli"
)

var runCommandHelp = `
	# Start new container
	can-cli run my-pod docker.io/ernoaapa/hello-world:latest
`

var runCommand = cli.Command{
	Name:        "run",
	HelpName:    "run",
	Usage:       "Start container in the device",
	Description: "With run command, you can start new containers in the device",
	UsageText: `can-cli run [options] NAME

	 # Start new container
	 can-cli run --image docker.io/eaapa/hello-world:latest my-pod
`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "image",
			Usage: "The container image to start",
		},
	},
	Action: func(clicontext *cli.Context) error {
		if clicontext.NArg() == 0 {
			return fmt.Errorf("You must give NAME parameter")
		}
		name := clicontext.Args().First()
		if name == "" {
			return fmt.Errorf("NAME argument cannot be empty")
		}

		if !clicontext.IsSet("image") || clicontext.String("image") == "" {
			return fmt.Errorf("You must define --image option")
		}
		image := clicontext.String("image")

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
