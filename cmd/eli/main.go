package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ernoaapa/eliot/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Get overrided at build time
var version = "master"
var commit = "unknown"
var date = time.Now().Format("2006-01-02_15:04:05")

func main() {
	app := cli.NewApp()
	app.Name = "eli"
	app.Usage = `commandline interface for managing eliot`
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
			EnvVar: "ELIOT_CONFIG",
			Value:  "~/.eli/config",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "Namespace to use with commands. By default reads from config.",
			EnvVar: "ELIOT_NAMESPACE",
		},
		cli.StringFlag{
			Name:   "endpoint",
			Usage:  "Use specific device endpoint. E.g. '192.168.1.101:5000'",
			EnvVar: "ELIOT_ENDPOINT",
		},
		cli.StringFlag{
			Name:   "device",
			Usage:  "Use specific device by name. E.g. 'somehost.local'",
			EnvVar: "ELIOT_DEVICE",
		},
	}, cmd.GlobalFlags...)
	app.Version = fmt.Sprintf("Version: %s, Commit: %s, Build at: %s", version, commit, date)
	app.Before = cmd.GlobalBefore

	app.Commands = []cli.Command{
		getCommand,
		describeCommand,
		deleteCommand,
		attachCommand,
		runCommand,
		upCommand,
		execCommand,
		createCommand,
		configCommand,
		buildCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
