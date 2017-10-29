package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/device"
	"github.com/ernoaapa/can/pkg/discovery"
	"github.com/ernoaapa/can/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = version.VERSION
	app.Name = "can-discovery"
	app.Usage = `Lightweight server to expose service over Bonjour (mDNS,DNS-SD) protocol`
	app.UsageText = `can-discovery [arguments...]
	
	 # If you use default 5000 port in the 'can-api' just
	 can-discovery

	 # If you customised the API port number, you need to add --expose flag
	 can-discovery --expose 5001`
	app.Description = `To make the 'can-api' discoverable over Bonjour protocol in the network, 
	 this lightweight server will listen Bonjour protocol and exposes the configured
	 port.`
	app.Flags = append([]cli.Flag{
		cli.IntFlag{
			Name:  "expose",
			Usage: "The port number to expose. Must match with can-api listen port",
			Value: 5000,
		}}, cmd.GlobalFlags...)
	app.Version = version.VERSION
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
