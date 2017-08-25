package main

import (
	"os"
	"time"

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
	app.Name = "layeryd"
	app.Usage = "Layery daemon"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		cli.StringFlag{
			Name:  "address, a",
			Usage: "address for containerd's GRPC server",
			Value: "/run/containerd/containerd.sock",
		},
		cli.DurationFlag{
			Name:  "timeout",
			Usage: "total timeout for ctr commands",
		},
		cli.DurationFlag{
			Name:  "status-interval",
			Usage: "interval for sending statuses",
			Value: 5 * time.Second,
		},
		cli.StringFlag{
			Name:   "namespace, n",
			Usage:  "namespace to use with commands",
			Value:  "layery",
			EnvVar: "LAYERYD_NAMESPACE",
		},
	}
	app.Commands = append([]cli.Command{
		runCommand,
	}, extraCmds...)

	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
