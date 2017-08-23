package main

import (
	"github.com/containerd/containerd"
	"github.com/ernoaapa/layeryd/controller"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Run the daemon",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Usage: "Read pod info from file",
		},
		cli.StringFlag{
			Name:  "interval",
			Usage: "Interval to fetch updates",
			Value: "10s",
		},
	},
	Action: func(clicontext *cli.Context) error {
		deviceInfo := getDeviceInfo()
		ctx, cancel := appContext(clicontext)
		defer cancel()

		source, err := getSource(clicontext)
		if err != nil {
			return err
		}

		reporter, err := getReporter(clicontext)
		if err != nil {
			return err
		}

		updates := source.GetUpdates(deviceInfo)
		log.Infoln("Started, start waiting for changes in source")

		for {
			var client *containerd.Client
			if client == nil {
				client, err = getContainerdClient(clicontext)
				if err != nil {
					log.Warnf("Failed to create containerd client: %s", err)
				}
			}

			if client != nil {
				state, err := controller.Sync(ctx, client, <-updates)
				if err != nil {
					log.Warnf("Failed to update state to containerd: %s", err)
					client = nil
				}

				if state != nil {
					err := reporter.Report(deviceInfo, *state)
					if err != nil {
						log.Warnf("Error while reporting current device state: %s", err)
					}
				}
			}
		}
	},
}
