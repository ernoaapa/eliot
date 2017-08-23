package main

import (
	"github.com/containerd/containerd"
	"github.com/ernoaapa/layeryd/controller"
	"github.com/ernoaapa/layeryd/model"
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
		ctx, cancel := appContext(clicontext)
		defer cancel()

		source, err := getSource(clicontext)
		if err != nil {
			return err
		}

		updates := source.GetUpdates(model.DeviceInfo{})
		for {
			var client *containerd.Client
			if client == nil {
				client, err = getContainerdClient(clicontext)
				if err != nil {
					log.Warnf("Failed to create containerd client", err.Error())
				}
			}

			if client != nil {
				if err := controller.Sync(ctx, client, <-updates); err != nil {
					log.Warnf("Failed to update state to containerd: %s", err)
					client = nil
				}
			}
		}
	},
}
