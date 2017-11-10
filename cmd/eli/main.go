package main

import (
	"os"

	"github.com/ernoaapa/elliot/cmd"
	"github.com/ernoaapa/elliot/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "eli"
	app.Usage = `commandline interface for managing elliot`
	app.UsageText = `eli [global options] command [command options] [arguments...]

	 # Detect devices
	 eli get devices

	 # Get running pods
	 eli get pods

	# Get pods in device
	eli --device hostname.local. get pods

	# See help of commands
	eli run --help
	`
	app.Description = `The 'eli' is tool for managing agent in the device.
	 With this tool, you can create, view and remove containers from the device.`
	app.Flags = append([]cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Client configuration",
			EnvVar: "ELLIOT_CONFIG",
			Value:  "~/.eli/config",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "Namespace to use with commands. By default reads from config.",
			EnvVar: "ELLIOT_NAMESPACE",
		},
		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "Use specific device endpoint. E.g. '192.168.1.101:5000'",
			EnvVar: "ELLIOT_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "device",
			Usage:  "Use specific device by name. E.g. 'somehost.local'",
			EnvVar: "ELLIOT_DEVICE",
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
