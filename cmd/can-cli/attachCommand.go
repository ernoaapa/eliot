package main

import (
	"fmt"
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/urfave/cli"
)

var attachCommand = cli.Command{
	Name:        "attach",
	HelpName:    "attach",
	Usage:       "Attach to container stdout and stderr output",
	Description: "You can use this command to get connection to container process and receive stdout and stderr output",
	UsageText: `can-cli attach [options] POD_NAME

	 # View pod attach
	 can-cli attach my-pod

	 # If pod contains multiple containers, you must define container id
	 can-cli attach --container some-id my-pod
`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "container, c",
			Usage: "Print logs of this container",
		},
	},
	Action: func(clicontext *cli.Context) error {
		client := cmd.GetClient(clicontext)

		if clicontext.NArg() == 0 || clicontext.Args().First() == "" {
			return fmt.Errorf("You must give Pod name as first argument")
		}
		podName := clicontext.Args().First()
		containerName := clicontext.String("container")

		pod, err := client.GetPod(podName)
		if err != nil {
			return err
		}

		containerCount := len(pod.Spec.Containers)
		if containerCount == 0 {
			return fmt.Errorf("Pod [%s] don't have any containers", podName)
		} else if containerCount == 1 {
			if containerName == "" {
				containerName = pod.Spec.Containers[0].Name
			}
		}
		if containerName == "" {
			return fmt.Errorf("Pod [%s] contains %d containers, you must define --container flag", podName, containerCount)
		}

		return client.Attach(containerName, os.Stdout, os.Stderr)
	},
}
