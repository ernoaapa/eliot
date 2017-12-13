package main

import (
	"errors"
	"os"

	"github.com/ernoaapa/eliot/cmd"
	"github.com/ernoaapa/eliot/pkg/api"
	"github.com/ernoaapa/eliot/pkg/controller"
	"github.com/ernoaapa/eliot/pkg/device"
	"github.com/ernoaapa/eliot/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "eliotd"
	app.Usage = "Daemon which contains all Eliot, for example GRPC API for the CLI client"
	app.UsageText = `eliotd [arguments...]

	 # By default listen port 5000
	 eliotd

	 # Listen custom port
	 eliotd --gprc-listen 0.0.0.0:5001
	 
	 # Disable lifecycle controller and enable only the GRPC API
	 eliotd  --grpc=true --lifecycle-controller=false`
	app.Description = `API for create/update/delete the containers and a way to connect into the containers.`
	app.Flags = append([]cli.Flag{
		cli.StringFlag{
			Name:   "containerd",
			Usage:  "containerd socket path for containerd's GRPC server",
			EnvVar: "ELIOT_CONTAINERD",
			Value:  "/run/containerd/containerd.sock",
		},
		cli.BoolTFlag{
			Name:   "lifecycle-controller",
			Usage:  "Enable container lifecycle controller",
			EnvVar: "ELIOT_LIFECYCLE_CONTROLLER",
		},
		cli.BoolTFlag{
			Name:   "grpc-api",
			Usage:  "Enable GRPC API server",
			EnvVar: "ELIOT_GRPC_API",
		},
		cli.StringFlag{
			Name:   "grpc-api-listen",
			Usage:  "GRPC host:port what to listen for client connections",
			EnvVar: "ELIOT_GRPC_API_LISTEN",
			Value:  "localhost:5000",
		},
	}, cmd.GlobalFlags...)
	app.Version = version.VERSION
	app.Before = cmd.GlobalBefore

	app.Action = func(clicontext *cli.Context) error {
		resolver := device.NewResolver(cmd.GetLabels(clicontext))
		device := resolver.GetInfo()
		client := cmd.GetRuntimeClient(clicontext, device.Hostname)
		listen := clicontext.String("grpc-api-listen")

		supervisor := suture.NewSimple("eliotd")

		if clicontext.Bool("grpc-api") {
			log.Infoln("grpc-api enabled")
			supervisor.Add(api.NewServer(listen, client))
		}

		if clicontext.Bool("lifecycle-controller") {
			log.Infoln("lifecycle-controller enabled")
			supervisor.Add(controller.NewLifecycle(client))
		}

		if len(supervisor.Services()) == 0 {
			return errors.New("Nothing to run. You should enable one of [grpc-api, lifecycle-controller]")
		}

		supervisor.Serve()

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
