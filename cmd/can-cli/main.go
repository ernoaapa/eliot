package main

import (
	"os"

	"github.com/coreos/etcd/version"
	"github.com/ernoaapa/can/pkg/client"
	"github.com/ernoaapa/can/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Version string to be set at compile time via command line (-ldflags "-X main.VersionString=1.2.3")
var (
	VersionString string
	extraCmds     = []cli.Command{}
)

func main() {
	app := cli.NewApp()
	app.Name = "can"
	app.Usage = "Can CLI"
	app.Version = version.Version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Can server API",
			Value: "/Users/ernoaapa/.can/config",
		},
	}

	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	app.Commands = []cli.Command{
		deploymentCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getClient(clicontext *cli.Context) *client.Client {
	configPath := clicontext.GlobalString("config")
	config, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Error while reading configuration file [%s]: %s", configPath, err)
	}
	return client.NewClient(
		config.GetCurrentEndpoint().URL,
		config.GetCurrentUser().Token,
	)
}
