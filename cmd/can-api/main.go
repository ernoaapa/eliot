package main

import (
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/api"
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
	app.Name = "can-api"
	app.Usage = "Can API server"
	app.Version = VersionString
	app.Flags = append([]cli.Flag{
		cli.StringFlag{
			Name:   "containerd",
			Usage:  "containerd socket path for containerd's GRPC server",
			EnvVar: "CAND_CONTAINERD",
			Value:  "/run/containerd/containerd.sock",
		},
		cli.StringFlag{
			Name:  "listen",
			Usage: "GRPC host:port what to listen for client connections",
			Value: "localhost:5000",
		},
	}, cmd.GlobalFlags...)

	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	app.Action = func(clicontext *cli.Context) error {
		client := cmd.GetRuntimeClient(clicontext)
		listen := clicontext.String("listen")
		server := api.NewServer(listen, client)

		log.Infof("Start to listen %s....", listen)
		return server.Serve()
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
