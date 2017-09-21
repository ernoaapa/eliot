package main

import (
	"github.com/ernoaapa/can/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var podsCommand = cli.Command{
	Name: "pods",
	Subcommands: []cli.Command{
		{
			Name:  "ls",
			Usage: "List pods",
			Action: func(clicontext *cli.Context) error {
				client := api.NewClient("localhost:5000")

				pods, err := client.GetPods()

				log.Printf("Created: %s", pods)
				return err
			},
		},
	},
}
