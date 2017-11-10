package main

import (
	"os"

	"github.com/ernoaapa/elliot/cmd"
	"github.com/ernoaapa/elliot/pkg/api"
	"github.com/ernoaapa/elliot/pkg/controller"
	"github.com/ernoaapa/elliot/pkg/device"
	"github.com/ernoaapa/elliot/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "elliotd"
	app.Usage = "Daemon which contains all Elliot, for example GRPC API for the CLI client"
	app.UsageText = `elliotd [arguments...]

	 # By default listen port 5000
	 elliotd
	
	 # Listen custom port
	 elliotd --listen 0.0.0.0:5001`
	app.Description = `API for create/update/delete the containers and a way to connect into the containers.`
	app.Flags = append([]cli.Flag{
		cli.StringFlag{
			Name:   "containerd",
			Usage:  "containerd socket path for containerd's GRPC server",
			EnvVar: "ELLIOT_CONTAINERD",
			Value:  "/run/containerd/containerd.sock",
		},
		cli.StringFlag{
			Name:   "listen",
			Usage:  "GRPC host:port what to listen for client connections",
			EnvVar: "ELLIOT_LISTEN",
			Value:  "localhost:5000",
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

		go func() {
			controller := controller.NewLifecycle(client)
			log.Infof("Start lifecycle controller...")
			err := controller.Run()
			log.Panicf("Lifecycle controller stopped with fatal error: %s", err)
		}()

		log.Infof("Start to listen %s....", listen)
		return server.Serve()
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
