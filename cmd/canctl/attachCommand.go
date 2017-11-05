package main

import (
	"fmt"
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/api"
	"github.com/ernoaapa/can/pkg/cmd/log"
	"github.com/ernoaapa/can/pkg/term"
	"github.com/urfave/cli"
)

var attachCommand = cli.Command{
	Name:        "attach",
	HelpName:    "attach",
	Usage:       "Attach to container stdout and stderr output",
	Description: "You can use this command to get connection to container process and receive stdout and stderr output",
	UsageText: `canctl attach [options] POD_NAME

	 # View pod attach
	 canctl attach my-pod

	 # If pod contains multiple containers, you must define container id
	 canctl attach --container some-id my-pod
`,
	Flags: []cli.Flag{
		cli.BoolTFlag{
			Name:  "stdin, i",
			Usage: "Keep stdin open on the container(s) in the pod, even if nothing is attached (default: true)",
		},
		cli.StringFlag{
			Name:  "container, c",
			Usage: "Print logs of this container",
		},
	},
	Action: func(clicontext *cli.Context) error {
		var (
			stdin  = os.Stdin
			stdout = os.Stdout
			stderr = os.Stderr
		)

		config := cmd.GetConfigProvider(clicontext)
		client := cmd.GetClient(config)

		if clicontext.NArg() == 0 || clicontext.Args().First() == "" {
			return fmt.Errorf("You must give Pod name as first argument")
		}
		podName := clicontext.Args().First()
		containerName := clicontext.String("container")

		pod, err := client.GetPod(podName)
		if err != nil {
			return err
		}

		containerID := ""
		containerCount := len(pod.Status.ContainerStatuses)
		if containerCount == 0 {
			return fmt.Errorf("Pod [%s] don't have any containers", podName)
		} else if containerCount == 1 {
			containerID = pod.Status.ContainerStatuses[0].ContainerID
		} else {
			if containerName == "" {
				return fmt.Errorf("Pod [%s] contains %d containers, you must define --container flag", podName, containerCount)
			}

			for _, status := range pod.Status.ContainerStatuses {
				if status.Name == containerName {
					containerID = status.ContainerID
				}
			}
		}
		if containerID == "" {
			return fmt.Errorf("Pod [%s] contains %d containers, you must define --container flag", podName, containerCount)
		}

		term := term.TTY{
			Out: stdout,
		}

		if clicontext.Bool("stdin") {
			term.In = stdin
			term.Raw = true
		}

		// Stop updating log lines, let the std piping take the terminal
		log.Stop()
		defer log.Start()

		return term.Safe(func() error {
			return client.Attach(containerID, api.NewAttachIO(term.In, term.Out, stderr))
		})
	},
}
