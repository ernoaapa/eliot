package main

import (
	"fmt"
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/urfave/cli"
)

var logsCommand = cli.Command{
	Name:  "logs",
	Usage: "View pod logs",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "container, c",
			Usage: "Print logs of this container ID",
		},
	},
	Action: func(clicontext *cli.Context) error {
		client := cmd.GetClient(clicontext)

		if clicontext.NArg() == 0 || clicontext.Args().First() == "" {
			return fmt.Errorf("You must give PODNAME argument")
		}
		podName := clicontext.Args().First()
		containerID := clicontext.String("container")

		pod, err := client.GetPod(podName)
		if err != nil {
			return err
		}

		containerCount := len(pod.Spec.Containers)
		if containerCount == 0 {
			return fmt.Errorf("Pod [%s] don't have any containers", podName)
		} else if containerCount == 1 {
			if containerID == "" {
				containerID = pod.Spec.Containers[0].ID
			}
		}
		if containerID == "" {
			return fmt.Errorf("Pod [%s] contains %d containers, you must define --container flag", podName, containerCount)
		}

		return client.GetLogs(containerID, os.Stdout, os.Stderr)
	},
}
