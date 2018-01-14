package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ernoaapa/eliot/cmd"
	"github.com/ernoaapa/eliot/pkg/api"
	"github.com/ernoaapa/eliot/pkg/api/core"
	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/cmd/ui"
	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/ernoaapa/eliot/pkg/progress"
	"github.com/ernoaapa/eliot/pkg/resolve"
	"github.com/ernoaapa/eliot/pkg/term"
	"github.com/ernoaapa/eliot/pkg/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:        "run",
	HelpName:    "run",
	Usage:       "Run commmand in new container in the device",
	Description: "With run command, you can run command in a new container in the device",
	UsageText: `eli run [options] <image> -- <command> <args>

	 # Run shell session in 'eaapa/hello-world' container
	 eli run -i -t eaapa/hello-world -- /bin/sh
`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Name for the pod (default: current directory name)",
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
		cli.StringFlag{
			Name:  "workdir, w",
			Usage: "Working directory inside the container",
		},
	},
	Action: func(clicontext *cli.Context) (err error) {

		var (
			name    = clicontext.String("name")
			image   = clicontext.Args().First()
			rm      = clicontext.Bool("rm")
			tty     = clicontext.Bool("tty")
			env     = clicontext.StringSlice("env")
			workdir = clicontext.String("workdir")
			mounts  = cmd.GetMounts(clicontext)
			binds   = cmd.GetBinds(clicontext)
			args    = cmd.DropDoubleDash(clicontext.Args().Tail())

			stdin  = os.Stdin
			stdout = os.Stdout
			stderr = os.Stderr
		)

		if clicontext.IsSet("detach") || clicontext.IsSet("d") {
			return errors.New("Detach not available. If you want to run container in background, use 'eli create pod' -command")
		}

		conf := cmd.GetConfigProvider(clicontext)
		client := cmd.GetClient(conf)

		if image == "" {
			uiline := ui.NewLine().Loading("Resolve image for the project...")
			info, err := client.GetInfo()
			if err != nil {
				uiline.Fatalf("Unable to resolve image for the project. Failed to get target device architecture: %s", err)
			}

			var projectType string
			projectType, image, err = resolve.Image(info.Arch, cmd.GetCurrentDirectory())
			if err != nil {
				uiline.Fatal("Unable to automatically resolve image for the project. You must define target container image with --image option")
			}
			uiline.Donef("Detected %s project, use image: %s (arch %s)", projectType, image, info.Arch)
		}

		image = utils.ExpandToFQIN(image)

		if name == "" {
			// Default to current directory name
			name = filepath.Base(cmd.GetCurrentDirectory())
		}

		for _, variable := range env {
			if !model.IsValidEnvKeyValuePair(variable) {
				return fmt.Errorf("Invalid --env value [%s], must be in format KEY=value. E.g. --env FOO=bar", variable)
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

		lines := map[string]ui.Line{}
		progressc := make(chan []*progress.ImageFetch)

		go func() {
			for fetches := range progressc {
				for _, fetch := range fetches {
					if _, ok := lines[fetch.Image]; !ok {
						lines[fetch.Image] = ui.NewLine().Loadingf("Download %s", fetch.Image)
					}

					if fetch.IsDone() {
						if fetch.Failed {
							lines[fetch.Image].Errorf("Failed %s", fetch.Image)
						} else {
							lines[fetch.Image].Donef("Downloaded %s", fetch.Image)
						}
					} else {
						current, total := fetch.GetProgress()
						lines[fetch.Image].WithProgress(current, total)
					}
				}
			}

			for image, line := range lines {
				line.Donef("Downloaded %s", image)
			}
		}()
		createErr := client.CreatePod(progressc, pod)
		close(progressc)
		if createErr != nil {
			return errors.Wrapf(createErr, "Error in creating pod")
		}

		if rm {
			defer func() {
				uiline := ui.NewLine().Loadingf("Delete pod %s", pod.Metadata.Name)
				_, err := client.DeletePod(pod)
				if err != nil {
					uiline.Errorf("Error while deleting pod [%s]: %s", pod.Metadata.Name, err)
				} else {
					uiline.Donef("Deleted pod [%s]", pod.Metadata.Name)
				}
			}()
		}

		result, err := client.StartPod(name)
		if err != nil {
			return errors.Wrapf(err, "Error in starting pod")
		}

		attachContainerID, err := findRunningContainerID(result, name)
		if err != nil {
			return errors.Wrapf(err, "Cannot attach to container")
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

		// Stop updating ui lines, let the std piping take the terminal
		ui.Stop()
		defer ui.Start()

		return term.Safe(func() error {
			return client.Attach(attachContainerID, api.NewAttachIO(term.In, term.Out, stderr))
		})
	},
}

func findRunningContainerID(pod *pods.Pod, name string) (string, error) {
	if pod.Status != nil && len(pod.Status.ContainerStatuses) > 0 {
		for _, status := range pod.Status.ContainerStatuses {
			if status.Name == name {
				return status.ContainerID, nil
			}
		}
	}

	return "", fmt.Errorf("Cannot find ContainerID with name %s", name)
}

func stopCatch(sigc chan os.Signal) {
	signal.Stop(sigc)
	close(sigc)
}
