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
		address := clicontext.GlobalString("address")
		namespace := clicontext.GlobalString("namespace")

		ctx, cancel := appContext(clicontext)
		defer cancel()

		client, err := containerd.New(address, containerd.WithDefaultNamespace(namespace))
		if err != nil {
			log.Warnf("Error from new client: %s", err)
			return err
		}

		source, err := getSource(clicontext)
		if err != nil {
			return err
		}

		updates := source.GetUpdates(model.NodeInfo{})
		for {
			if err := controller.Sync(ctx, client, <-updates); err != nil {
				return err
			}
		}
	},
}
