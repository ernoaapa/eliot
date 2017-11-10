package main

import (
	"github.com/ernoaapa/elliot/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var configSetCommand = cli.Command{
	Name:      "set",
	Usage:     "Set client config parameter",
	UsageText: "eli config set NAME VALUE",
	Action: func(clicontext *cli.Context) error {
		if clicontext.NArg() != 2 {
			log.Fatalf("You must give two parameters, NAME and VALUE")
		}

		var (
			name  = clicontext.Args()[0]
			value = clicontext.Args()[1]
		)

		config := cmd.GetConfig(clicontext)
		config.Set(name, value)
		return cmd.UpdateConfig(clicontext, config)
	},
}
