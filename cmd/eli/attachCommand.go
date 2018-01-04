package main

import (
	"fmt"
	"os"

	"github.com/ernoaapa/eliot/cmd"
	"github.com/ernoaapa/eliot/pkg/api"
	"github.com/ernoaapa/eliot/pkg/cmd/log"
	"github.com/ernoaapa/eliot/pkg/term"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var attachCommand = cli.Command{
	Name:        "attach",
	HelpName:    "attach",
	Usage:       "Attach to container stdout and stderr output",
	Description: "You can use this command to get connection to container process and receive stdout and stderr output",
	UsageText: `eli attach [options] POD_NAME

	 # View pod attach
	 eli attach my-pod

	 # If pod contains multiple containers, you must define container id
	 eli attach --container some-id my-pod
`,
	Flags: []cli.Flag{
		cli.BoolTFlag{
			Name:  "stdin, i",
			Usage: "Keep stdin open on the container(s) in the pod, even if nothing is attached (default: true)",
		},
		cli.StringFlag{
			Name:  "container, c",
			Usage: "Target container in the pod",
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

		containerID, err := cmd.ResolveContainerID(pod.Status.ContainerStatuses, containerName)
		if err != nil {
			return errors.Wrapf(err, "Failed to resolve containerID for pod [%s]", podName)
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
