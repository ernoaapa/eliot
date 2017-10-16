package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/api/core"
	containers "github.com/ernoaapa/can/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/config"
	"github.com/ernoaapa/can/pkg/model"
	"github.com/ernoaapa/can/pkg/printers"
	"github.com/ernoaapa/can/pkg/resolve"
	"github.com/ernoaapa/can/pkg/sync"
	"github.com/ernoaapa/can/pkg/term"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:        "run",
	HelpName:    "run",
	Usage:       "Start container in the device",
	Description: "With run command, you can start new containers in the device",
	UsageText: `can-cli run [options] -- <command>

	 # Run code in current directory in the device
	 can-cli run

	 # Run 'build.sh' command in device with files in current directory
	 can-cli run -- ./build.sh
	 
	 # Run container image in the container
	 can-cli run --image docker.io/eaapa/hello-world:latest

	 # Run container with name in the device
	 can-cli run --image docker.io/eaapa/hello-world:latest --name my-pod
`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Name for the pod (default: current directory name)",
		},
		cli.StringFlag{
			Name:  "image",
			Usage: "The container image to start (default: resolve image based on project structure)",
		},
		cli.BoolFlag{
			Name:  "detach, d",
			Usage: "Run container in background and print container information (default: false)",
		},
		cli.StringSliceFlag{
			Name:  "mount",
			Usage: "Attach a filesystem mount to the container",
		},
		cli.StringSliceFlag{
			Name:  "bind",
			Usage: "Bind a directory in host to the container. Format: /source:/target:options, E.g. /var:/var:rshared",
		},
		cli.StringSliceFlag{
			Name:  "env, e",
			Usage: "Set environment variable into the container. E.g. --env FOO=bar",
		},
		cli.BoolTFlag{
			Name:  "rm",
			Usage: "Automatically remove the container when it exits (default: true)",
		},
		cli.BoolTFlag{
			Name:  "tty, t",
			Usage: "Allocate TTY for each container in the pod (default: true)",
		},
		cli.BoolTFlag{
			Name:  "stdin, i",
			Usage: "Keep stdin open on the container(s) in the pod, even if nothing is attached (default: true)",
		},
		cli.BoolFlag{
			Name:  "no-sync",
			Usage: "Do not sync any directory to the container (default: false)",
		},
		cli.StringSliceFlag{
			Name:  "sync",
			Usage: "Directory to sync to the target container (default: current directory)",
			Value: &cli.StringSlice{"."},
		},
		cli.StringFlag{
			Name:  "workdir, w",
			Usage: "Working directory inside the container",
		},
	},
	Action: func(clicontext *cli.Context) (err error) {

		var (
			projectConfig = config.ReadProjectConfig("./.can.yml")
			name          = projectConfig.NameOrElse(clicontext.String("name"))
			image         = projectConfig.ImageOrElse(clicontext.String("image"))
			detach        = clicontext.Bool("detach")
			rm            = clicontext.Bool("rm")
			tty           = clicontext.Bool("tty")
			env           = projectConfig.EnvWith(clicontext.StringSlice("env"))
			workdir       = clicontext.String("workdir")
			noSync        = clicontext.Bool("no-sync")
			syncDirs      = clicontext.StringSlice("sync")
			mounts        = cmd.GetMounts(clicontext)
			binds         = cmd.GetBinds(clicontext)
			args          = clicontext.Args()

			stdin  = os.Stdin
			stdout = os.Stdout
			stderr = os.Stderr
		)

		if image == "" {
			log.Println("No image defined, try to detect image for project...")
			image, err = resolve.Image(cmd.GetCurrentDirectory())
			if err != nil {
				log.Debugf("Unable to resolve automatically container image for project in directory [%s]. Error: %s", cmd.GetCurrentDirectory(), err)
				return fmt.Errorf("Unable to detect image for project. You must define target container image with --image option")
			}
			log.Printf("Auto-resolved image for project. Will use image: %s", image)
		}
		image = cmd.ExpandToFQIN(image)

		if name == "" {
			// Default to current directory name
			name = filepath.Base(cmd.GetCurrentDirectory())
		}

		if detach && rm {
			return fmt.Errorf("You cannot use --detach flag with --rm, it would remove right away after container started")
		}

		for _, variable := range env {
			if !model.IsValidEnvKeyValuePair(variable) {
				return fmt.Errorf("Invalid --env value [%s], must be in format KEY=value. E.g. --env FOO=bar", variable)
			}
		}

		config := cmd.GetConfigProvider(clicontext)
		client := cmd.GetClient(config)

		cont := []*containers.Container{
			&containers.Container{
				Name:       name,
				Image:      image,
				Tty:        tty,
				Args:       args,
				Env:        env,
				WorkingDir: workdir,
				Mounts:     append(mounts, binds...),
			},
		}

		if !noSync {
			workspaceMount, _ := cmd.ParseBindFlag(fmt.Sprintf("/var/lib/volumes/%s:%s:rw,rshared", name, projectConfig.Sync.Target))

			cont[0].Mounts = append(cont[0].Mounts, workspaceMount)
			if cont[0].WorkingDir == "" {
				cont[0].WorkingDir = projectConfig.Sync.Target
			}

			cont = append(cont, &containers.Container{
				Name:   fmt.Sprintf("rsync-%s", name),
				Image:  cmd.ExpandToFQIN(projectConfig.Sync.Image),
				Mounts: []*containers.Mount{workspaceMount},
			})
		}

		pod := &pods.Pod{
			Metadata: &core.ResourceMetadata{
				Name:      name,
				Namespace: config.GetNamespace(),
			},
			Spec: &pods.PodSpec{
				HostNetwork: true,
				Containers:  cont,
			},
		}

		if err := client.CreatePod(pod); err != nil {
			return errors.Wrapf(err, "Error in creating pod")
		}

		result, err := client.StartPod(name)
		if err != nil {
			return errors.Wrapf(err, "Error in starting pod")
		}

		if rm {
			defer func() {
				pod, err := client.DeletePod(pod)
				if err != nil {
					log.Errorf("Error while deleting pod [%s]: %s", pod.Metadata.Name, err)
				} else {
					log.Infof("Deleted pod [%s]", pod.Metadata.Name)
				}
			}()
		}

		if detach {
			writer := printers.GetNewTabWriter(os.Stdout)
			printer := cmd.GetPrinter(clicontext)
			return printer.PrintPodDetails(result, writer)
		}

		// TODO: Switch to created ContainerID when API exposes it
		attachContainerID := name

		if !noSync {
			done := make(chan struct{})
			destination := fmt.Sprintf("rsync://%s:%d/volume", config.GetEndpointHost(), 873)
			sync.StartRsync(done, syncDirs, destination, 1*time.Second)
			defer close(done)
		}

		term := term.TTY{
			Out: stdout,
		}

		if clicontext.Bool("stdin") {
			term.In = stdin
			term.Raw = true
		} else {
			sigc := cmd.ForwardAllSignals(func(signal syscall.Signal) error {
				return client.Signal(attachContainerID, signal)
			})
			defer stopCatch(sigc)
		}

		return term.Safe(func() error {
			log.Debugln("Attach to container [%s]", attachContainerID)
			return client.Attach(attachContainerID, term.In, term.Out, stderr)
		})
	},
}

func stopCatch(sigc chan os.Signal) {
	signal.Stop(sigc)
	close(sigc)
}
