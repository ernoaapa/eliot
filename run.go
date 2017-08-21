package main

import (
	"time"

	"github.com/containerd/containerd"
	"github.com/ernoaapa/layeryd/controller"
	"github.com/ernoaapa/layeryd/source"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Run the daemon",

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

		source := source.NewFileSource("./test.yml")

		for {
			if err := controller.Sync(ctx, client, source); err != nil {
				return err
			}
			time.Sleep(2 * time.Second)
		}
	},
}
