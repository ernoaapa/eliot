package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ernoaapa/eliot/cmd"
	"github.com/ernoaapa/eliot/pkg/api"
	"github.com/ernoaapa/eliot/pkg/controller"
	"github.com/ernoaapa/eliot/pkg/discovery"
	"github.com/ernoaapa/eliot/pkg/node"
	"github.com/ernoaapa/eliot/pkg/profile"
	log "github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
	"github.com/urfave/cli"
)

// Get overrided at build time
var version = "master"
var commit = "unknown"
var date = time.Now().Format("2006-01-02_15:04:05")

func main() {
	app := cli.NewApp()
	app.Name = "eliotd"
	app.Usage = "Daemon for the node to enable Eliot"
	app.UsageText = `eliotd [arguments...]

	 # By default listen port 5000
	 eliotd

	 # Listen custom port
	 eliotd --grpc-api-listen 0.0.0.0:5001
	 
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
		cli.StringFlag{
			Name:   "containerd-snapshotter",
			Usage:  "containerd snapshotter to use",
			EnvVar: "ELIOT_CONTAINERD_SNAPSHOTTER",
			Value:  "overlayfs",
		},
		cli.DurationFlag{
			Name:   "timeout, t",
			Usage:  "total timeout for runtime requests",
			EnvVar: "ELIOT_TIMEOUT",
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
		cli.BoolTFlag{
			Name:   "discovery",
			Usage:  "Enable discover GRPC server over zeroconf",
			EnvVar: "ELIOT_DISCOVERY",
		},
		cli.BoolFlag{
			Name:   "profile",
			Usage:  "Turn on pprof profiling",
			EnvVar: "ELIOT_PROFILE",
		},
		cli.StringFlag{
			Name:   "profile-address",
			Usage:  "The http address for the pprof server",
			EnvVar: "ELIOT_PROFILE_ADDRESS",
			Value:  "0.0.0.0:8000",
		},
		cli.StringFlag{
			Name:   "labels",
			Usage:  "Comma separated list of node labels. E.g. --labels node=rpi3,location=home,environment=testing",
			EnvVar: "ELIOT_LABELS",
		},
	}, cmd.GlobalFlags...)
	app.Version = fmt.Sprintf("Version: %s, Commit: %s, Build at: %s", version, commit, date)
	app.Before = cmd.GlobalBefore

	app.Action = func(clicontext *cli.Context) error {
		var (
			grpcListen = clicontext.String("grpc-api-listen")
			grpcPort   = parseGrpcPort(grpcListen)
		)

		resolver := node.NewResolver(grpcPort, version, cmd.GetLabels(clicontext))
		node := resolver.GetInfo()
		client := cmd.GetRuntimeClient(clicontext, node.Hostname)

		supervisor := suture.NewSimple("eliotd")
		serviceCount := 0

		if clicontext.BoolT("profile") {
			profileAddr := clicontext.String("profile-address")
			log.Infof("profiling enabled, address: %s", profileAddr)
			supervisor.Add(profile.NewServer(profileAddr))
			serviceCount++
		}

		if clicontext.Bool("grpc-api") {
			log.Infoln("grpc-api enabled")
			supervisor.Add(api.NewServer(grpcListen, client, resolver))
			serviceCount++
		}

		if clicontext.Bool("lifecycle-controller") {
			log.Infoln("lifecycle-controller enabled")
			supervisor.Add(controller.NewLifecycle(client))
			serviceCount++
		}

		if clicontext.Bool("grpc-api") && clicontext.Bool("discovery") {
			log.Infoln("grpc discovery over zeroconf enabled")
			supervisor.Add(discovery.NewServer(node.Hostname, grpcPort, version))
			serviceCount++
		}

		if serviceCount == 0 {
			return errors.New("Nothing to run. You should enable one of [grpc-api, lifecycle-controller, discovery]")
		}

		supervisor.Serve()

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func parseGrpcPort(addr string) int {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		log.Panicf("Invalid formated grpc address [%s]", addr)
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Panicf("Unable to parse grpc port: %s", err)
	}
	return port
}
