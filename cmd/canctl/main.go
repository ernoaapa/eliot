package main

import (
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "canctl"
	app.Usage = `commandline interface for managing can`
	app.UsageText = `canctl [global options] command [command options] [arguments...]

	 # Detect devices
	 canctl get devices

	 # Get running pods
	 canctl get pods

	# Get pods in device
	canctl --device hostname.local. get pods

	# See help of commands
	canctl run --help
	`
	app.Description = `The 'canctl' is tool for managing agent in the device.
	 With this tool, you can create, view and remove containers from the device.`
	app.Flags = append([]cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Client configuration",
			EnvVar: "CAN_CONFIG",
			Value:  "~/.can/config",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "Namespace to use with commands. By default reads from config.",
			EnvVar: "CAN_NAMESPACE",
		},
		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "Use specific device endpoint. E.g. '192.168.1.101:5000'",
			EnvVar: "CAN_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "device",
			Usage:  "Use specific device by name. E.g. 'somehost.local'",
			EnvVar: "CAN_DEVICE",
		},
	}, cmd.GlobalFlags...)
	app.Version = version.VERSION
	app.Before = cmd.GlobalBefore

	app.Commands = []cli.Command{
		getCommand,
		describeCommand,
		deleteCommand,
		attachCommand,
		runCommand,
		createCommand,
		configCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
