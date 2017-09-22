package main

import (
	"os"
	"time"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/controller"
	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/model"
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
	app.Name = "cand"
	app.Usage = "Can daemon"
	app.Version = VersionString
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
			Name:  "labels",
			Usage: "Labels to add to the device info.  Labels must be key=value pairs separated by ','.",
		},

		cli.StringFlag{
			Name:  "manifest",
			Usage: "url path to manifest file. E.g. file:///some/path/to/file.yaml",
		},

		cli.StringFlag{
			Name:  "report",
			Usage: "Where to send pod status. E.g. 'console' or 'http://foo.bar.com'",
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
		resolver := device.NewResolver(cmd.GetLabels(clicontext))
		client := cmd.GetRuntimeClient(clicontext)

		sourceUpdates := make(chan []model.Pod)

		source, err := cmd.GetManifestSource(clicontext, resolver, sourceUpdates)
		if err != nil {
			return err
		}
		go source.Start()

		reporter, err := cmd.GetStateReporter(clicontext, resolver, client)
		if err != nil {
			return err
		}
		go reporter.Start()

		controller := controller.New(client)

		log.Infoln("Started, start waiting for changes in source")

		for {
			err := controller.Sync(<-sourceUpdates)
			if err != nil {
				log.Warnf("Failed to update state to containerd: %s", err)
			}
		}
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
