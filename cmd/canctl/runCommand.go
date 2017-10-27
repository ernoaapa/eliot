package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ernoaapa/can/cmd"
	"github.com/ernoaapa/can/pkg/api"
	"github.com/ernoaapa/can/pkg/api/core"
	containers "github.com/ernoaapa/can/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/config"
	"github.com/ernoaapa/can/pkg/display"
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
	UsageText: `canctl run [options] -- <command>

	 # Run code in current directory in the device
	 canctl run

	 # Run 'build.sh' command in device with files in current directory
	 canctl run -- ./build.sh
	 
	 # Run container image in the container
	 canctl run --image docker.io/eaapa/hello-world:latest

	 # Run container with name in the device
	 canctl run --image docker.io/eaapa/hello-world:latest --name my-pod
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
			binds         = cmd.GetBinds(clicontext, projectConfig.Binds...)
			args          = clicontext.Args()

			stdin  = os.Stdin
			stdout = os.Stdout
			stderr = os.Stderr
		)

		if image == "" {
			display := display.NewLine()
			display.Active("Resolve image for the project...")
			projectType, image, err := resolve.Image(cmd.GetCurrentDirectory())
			if err != nil {
				display.Error("Unable to automatically resolve container image for the project. You must define target container image with --image option")
				return fmt.Errorf("Unable to detect image for project. You must define target container image with --image option")
			}
			display.Donef("Detected %s project, use image: %s", projectType, image)
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

		conf := cmd.GetConfigProvider(clicontext)
		client, err := cmd.GetClient(conf)
		if err != nil {
			return err
		}

		opts := []api.PodOpts{}

		if !noSync {
			syncTargetPath := projectConfig.Sync.Target

			opts = append(opts, api.WithContainer(&containers.Container{
				Name:  fmt.Sprintf("rsync-%s", name),
				Image: cmd.ExpandToFQIN(projectConfig.Sync.Image),
				Env: []string{
					fmt.Sprintf("VOLUME=%s", syncTargetPath),
				},
			}))

			opts = append(opts, api.WithSharedMount(
				cmd.MustParseBindFlag(fmt.Sprintf("/var/lib/volumes/%s:%s:rw,rshared", name, syncTargetPath)),
			))

			if !clicontext.IsSet("workdir") {
				opts = append(opts, api.WithWorkingDir(syncTargetPath))
			}
		}

		pod := &pods.Pod{
			Metadata: &core.ResourceMetadata{
				Name:      name,
				Namespace: conf.GetNamespace(),
			},
			Spec: &pods.PodSpec{
				HostNetwork: true,
				Containers: []*containers.Container{
					&containers.Container{
						Name:       name,
						Image:      image,
						Tty:        tty,
						Args:       args,
						Env:        env,
						WorkingDir: workdir,
						Mounts:     append(mounts, binds...),
					},
				},
			},
		}

		if err := client.CreatePod(pod, opts...); err != nil {
			return errors.Wrapf(err, "Error in creating pod")
		}

		result, err := client.StartPod(name)
		if err != nil {
			return errors.Wrapf(err, "Error in starting pod")
		}

		if rm {
			defer func() {
				display := display.NewLine()
				display.Activef("Delete pod %s", pod.Metadata.Name)
				_, err := client.DeletePod(pod)
				if err != nil {
					display.Errorf("Error while deleting pod [%s]: %s", pod.Metadata.Name, err)
				} else {
					display.Donef("Deleted pod [%s]", pod.Metadata.Name)
				}
			}()
		}

		if detach {
			writer := printers.GetNewTabWriter(os.Stdout)
			defer writer.Flush()
			printer := cmd.GetPrinter(clicontext)
			return printer.PrintPodDetails(result, writer)
		}

		// TODO: Switch to created ContainerID when API exposes it
		attachContainerID := name

		hooks := []api.AttachHooks{}
		if !noSync {
			hooks = append(hooks, func(endpoint config.Endpoint, done <-chan struct{}) {
				destination := fmt.Sprintf("rsync://%s:%d/volume", endpoint.GetHost(), 873)
				sync.StartRsync(done, syncDirs, destination, 1*time.Second)
			})
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

		// Stop updating display output, let the std piping take the terminal
		display.Stop()
		defer display.Start()

		return term.Safe(func() error {
			log.Debugln("Attach to container [%s]", attachContainerID)
			return client.Attach(attachContainerID, api.NewAttachIO(term.In, term.Out, stderr), hooks...)
		})
	},
}

func stopCatch(sigc chan os.Signal) {
	signal.Stop(sigc)
	close(sigc)
}
