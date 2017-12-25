package main

import (
	"fmt"
	"io/ioutil"

	"github.com/ernoaapa/eliot/pkg/cmd"
	"github.com/ernoaapa/eliot/pkg/cmd/log"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var deviceBuildCommand = cli.Command{
	Name:    "build",
	Aliases: []string{"b"},
	Usage:   "Build device image",
	UsageText: `eli device build [options]
	
	 # Build device image
	 eli device build
	 
	 # Create Linuxkit file but don't build it
	 eli device build --dry-run
	 eli device build --dry-run > linuxkit.yml`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Usage: "Linuxkit build source file",
		},
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Print the final Linuxkit config and don't actually build it",
		},
	},
	Action: func(clicontext *cli.Context) (err error) {
		log := log.NewLine().Loading("Get Linuxkit config...")
		var (
			file     = clicontext.String("file")
			dryRun   = clicontext.Bool("dry-run")
			linuxkit []byte
		)

		if file != "" {
			linuxkit, err = ioutil.ReadFile(file)
			if err != nil {
				log.Errorf("Failed to read Linuxkit file: %s", err)
				return err
			}
		} else if cmd.IsPipingIn() {
			linuxkit, err = cmd.ReadAllStdin()
			if err != nil {
				log.Errorf("Failed to read Linuxkit config from stdin: %s", err)
			}
		} else {
			log.Errorf("You must define --file and give path to Linuxkit config file!")
			return errors.New("No Linuxkit config defined")
		}

		if len(linuxkit) == 0 {
			log.Errorf("Invalid Linuxkit config!")
		}

		log.Donef("Resolved Linuxkit config!")

		if dryRun {
			fmt.Println(string(linuxkit))
			return nil
		}
		return errors.New("Not implemented")
	},
}
