package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/discovery"
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
	app.Name = "discovery"
	app.Usage = "Server to discover device with Bonjour (mDNS,DNS-SD) protocol"
	app.Version = VersionString
	app.Flags = append([]cli.Flag{
		cli.IntFlag{
			Name:  "expose",
			Usage: "The port number to expose. Must match with can-api listen port",
			Value: 5000,
		}}, cmd.GlobalFlags...)
	app.Before = cmd.GlobalBefore

	app.Action = func(clicontext *cli.Context) error {
		resolver := device.NewResolver(cmd.GetLabels(clicontext))
		device := resolver.GetInfo()

		server := discovery.NewServer(
			device.Hostname,
			clicontext.Int("expose"),
		)
		err := server.Start()
		if err != nil {
			panic(err)
		}
		defer server.Stop()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		log.Println("Shutting down.")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
