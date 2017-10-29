package main

import (
	"os"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/api"
	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "can-api"
	app.Usage = "Provides GRPC API for the CLI client"
	app.UsageText = `can-api [arguments...]

	 # By default listen port 5000
	 can-api
	
	 # Listen custom port
	 can-api --listen 0.0.0.0:5001`
	app.Description = `API for create/update/delete the containers and a way to connect into the containers.`
	app.Flags = append([]cli.Flag{
		cli.StringFlag{
			Name:   "containerd",
			Usage:  "containerd socket path for containerd's GRPC server",
			EnvVar: "CAN_CONTAINERD",
			Value:  "/run/containerd/containerd.sock",
		},
		cli.StringFlag{
			Name:  "listen",
			Usage: "GRPC host:port what to listen for client connections",
			Value: "localhost:5000",
		},
	}, cmd.GlobalFlags...)
	app.Version = version.VERSION
	app.Before = cmd.GlobalBefore

	app.Action = func(clicontext *cli.Context) error {
		resolver := device.NewResolver(cmd.GetLabels(clicontext))
		device := resolver.GetInfo()
		client := cmd.GetRuntimeClient(clicontext, device.Hostname)
		listen := clicontext.String("listen")
		server := api.NewServer(listen, client)

		log.Infof("Start to listen %s....", listen)
		return server.Serve()
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
