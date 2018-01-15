package main

import (
	"fmt"
	"os"

	"github.com/ernoaapa/eliot/cmd"
	"github.com/ernoaapa/eliot/pkg/api"
	"github.com/ernoaapa/eliot/pkg/cmd/ui"
	"github.com/ernoaapa/eliot/pkg/term"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var execCommand = cli.Command{
	Name:        "exec",
	HelpName:    "exec",
	Usage:       "Execute a command in a running container",
	Description: "You can use this command to run command in container process",
	UsageText: `eli exec [options] POD_NAME

	 # Run 'date' in my-pod
	 eli exec my-pod date

	 # If pod contains multiple containers, you must define container id
	 eli exec --container some-id my-pod date

	 # If you have parameters with command, add double dash (--) to separate
	 # command from the eli command
	 eli exec --container some-id my-pod -- ls -lt /usr
`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "tty, t",
			Usage: "Allocate TTY for each container in the pod",
		},
		cli.BoolFlag{
			Name:  "stdin, i",
			Usage: "Keep stdin open on the container(s) in the pod, even if nothing is attached",
		},
		cli.StringFlag{
			Name:  "container, c",
			Usage: "Container name. If omitted, the first container in the pod will be chosen",
		},
	},
	Action: func(clicontext *cli.Context) error {
		var (
			stdin  = os.Stdin
			stdout = os.Stdout
			stderr = os.Stderr

			tty           = clicontext.Bool("tty")
			podName       = clicontext.Args().First()
			containerName = clicontext.String("container")
			args          = cmd.DropDoubleDash(clicontext.Args().Tail())
		)

		config := cmd.GetConfigProvider(clicontext)
		client := cmd.GetClient(config)

		if clicontext.NArg() == 0 || clicontext.Args().First() == "" {
			return fmt.Errorf("You must give Pod name as first argument")
		}

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

		// Stop updating ui lines, let the std piping take the terminal
		ui.Stop()
		defer ui.Start()

		return term.Safe(func() error {
			return client.Exec(containerID, args, tty, api.NewAttachIO(term.In, term.Out, stderr))
		})
	},
}
