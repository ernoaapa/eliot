package main

import (
	"os"
	"time"

	"github.com/ernoaapa/layeryd/controller"
	"github.com/ernoaapa/layeryd/device"
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
			Usage: "total timeout for containerd requests",
		},

		cli.StringFlag{
			Name:  "manifest",
			Usage: "url path to manifest file. E.g. file:///some/path/to/file.yaml",
		},

		cli.StringFlag{
			Name:  "manifest-update-interval",
			Usage: "Interval to update desired state",
			Value: "10s",
		},

		cli.DurationFlag{
			Name:  "state-update-interval",
			Usage: "interval for reporting current state",
			Value: 5 * time.Second,
		},
	}

	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	app.Action = func(clicontext *cli.Context) error {
		deviceInfo := device.GetInfo()
		client := getRuntimeClient(clicontext)

		source, err := getManifestSource(clicontext)
		if err != nil {
			return err
		}

		reporter, err := getStateReporter(clicontext, deviceInfo, client)
		if err != nil {
			return err
		}
		go reporter.Start()

		changes := source.GetUpdates()
		log.Infoln("Started, start waiting for changes in source")

		for {
			err := controller.Sync(client, <-changes)
			if err != nil {
				log.Warnf("Failed to update state to containerd: %s", err)
			}
		}
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
