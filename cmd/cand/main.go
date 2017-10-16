package main

import (
	"os"
	"sync"
	"time"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	// Version string to be set at compile time via command line (-ldflags "-X main.VersionString=1.2.3")
	VersionString string
)

func main() {
	app := cli.NewApp()
	app.Name = "cand"
	app.Usage = `The primary "device agent"

	The 'cand' is the primary "device agent" that runs on each device.
	The agent takes a list of Pod specifications and that are provided through
	various mechanisms and ensures that the containers described in those
	specifications are running and healthy.
	`
	app.UsageText = `
	# Defaults usually is enough
	cand

	# If containerd socket is stored somewhere else
	cand --containerd /some/where/else/containerd.sock

	# To get debug output
	cand --debug
	`
	app.Version = VersionString
	app.Before = cmd.GlobalBefore
	app.Flags = append([]cli.Flag{
		cli.StringFlag{
			Name:   "containerd",
			Usage:  "containerd for containerd's GRPC server",
			EnvVar: "CAND_CONTAINERD",
			Value:  "/run/containerd/containerd.sock",
		},
		cli.DurationFlag{
			Name:   "timeout",
			EnvVar: "CAND_TIMEOUT",
			Usage:  "Timeout for containerd requests",
		},

		cli.StringFlag{
			Name:   "labels",
			Usage:  "Labels to add to the device info.  Labels must be key=value pairs separated by ','",
			EnvVar: "CAND_LABELS",
		},

		cli.StringFlag{
			Name:  "manifest",
			Usage: "url path to manifest file. E.g. file:///some/path/to/file.yaml",
		},

		cli.StringFlag{
			Name:  "report",
			Usage: "Where to send pod status. E.g. 'console' or 'http://foo.bar.com'",
		},

		cli.DurationFlag{
			Name:  "manifest-update-interval",
			Usage: "Interval to update desired state",
			Value: 10 * time.Second,
		},

		cli.DurationFlag{
			Name:  "update-interval",
			Usage: "Interval for updating state",
			Value: 1 * time.Second,
		},
	}, cmd.GlobalFlags...)

	app.Action = func(clicontext *cli.Context) error {
		var wg sync.WaitGroup
		resolver := device.NewResolver(cmd.GetLabels(clicontext))

		sourceUpdates := make(chan []model.Pod)
		stateUpdates := make(chan []model.Pod)

		source, err := cmd.GetManifestSource(clicontext, resolver, sourceUpdates)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			source.Start()
		}()

		reporter, err := cmd.GetStateReporter(clicontext, resolver, stateUpdates)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			reporter.Start()
		}()

		controller := cmd.GetController(clicontext, sourceUpdates, stateUpdates)
		wg.Add(1)
		go func() {
			defer wg.Done()
			controller.Start()
		}()

		log.Infoln("Started!")
		wg.Wait()
		log.Infoln("Source, reporter and controller died. Shutting down")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
