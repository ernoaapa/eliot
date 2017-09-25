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
	Action: func(clicontext *cli.Context) error {
		client := cmd.GetClient(clicontext)

		if clicontext.NArg() == 0 || clicontext.Args().First() == "" {
			return fmt.Errorf("You must give CONTAINERID argument")
		}
		containerID := clicontext.Args().First()

		return client.GetLogs(containerID, os.Stdout, os.Stderr)
	},
}
