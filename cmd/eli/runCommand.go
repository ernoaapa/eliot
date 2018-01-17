package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ernoaapa/eliot/cmd"
	"github.com/ernoaapa/eliot/pkg/api"
	"github.com/ernoaapa/eliot/pkg/api/core"
	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/eliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/eliot/pkg/cmd/ui"
	"github.com/ernoaapa/eliot/pkg/config"
	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/ernoaapa/eliot/pkg/printers"
	"github.com/ernoaapa/eliot/pkg/progress"
	"github.com/ernoaapa/eliot/pkg/resolve"
	"github.com/ernoaapa/eliot/pkg/sync"
	"github.com/ernoaapa/eliot/pkg/term"
	"github.com/ernoaapa/eliot/pkg/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:        "run",
	HelpName:    "run",
	Usage:       "Start container in the device",
	Description: "With run command, you can start new containers in the device",
	UsageText: `eli run [options] -- <command>

	 # Run code in current directory in the device
	 eli run

	 # Run 'build.sh' command in device with files in current directory
	 eli run -- ./build.sh
	 
	 # Run container image in the container
	 eli run --image docker.io/eaapa/hello-world:latest

	 # Run container with name in the device
	 eli run --image docker.io/eaapa/hello-world:latest --name my-pod
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
			projectConfig = config.ReadProjectConfig("./.eliot.yml")
			name          = cmd.First(clicontext.String("name"), projectConfig.Name)
			image         = cmd.First(clicontext.String("image"), projectConfig.Image)
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

		if detach && rm {
			return fmt.Errorf("You cannot use --detach flag with --rm, it would remove right away after container started")
		}

		for _, variable := range env {
			if !model.IsValidEnvKeyValuePair(variable) {
				return fmt.Errorf("Invalid --env value [%s], must be in format KEY=value. E.g. --env FOO=bar", variable)
			}
		}

		opts := []api.PodOpts{}

		if !noSync {
			syncTargetPath := projectConfig.Sync.Target

			opts = append(opts, api.WithContainer(&containers.Container{
				Name:  fmt.Sprintf("rsync-%s", name),
				Image: utils.ExpandToFQIN(projectConfig.Sync.Image),
				Env: []string{
					fmt.Sprintf("VOLUME=%s", syncTargetPath),
				},
			}))

			opts = append(opts, api.WithSharedMount(
				cmd.MustParseBindFlag(fmt.Sprintf("/var/lib/volumes/%s:%s:rw,rshared", name, syncTargetPath)),
			))

			if !clicontext.IsSet("workdir") && !clicontext.IsSet("w") {
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
		}()
		createErr := client.CreatePod(progressc, pod, opts...)
		close(progressc)
		if createErr != nil {
			return errors.Wrapf(createErr, "Error in creating pod")
		}

		result, err := client.StartPod(name)
		if err != nil {
			return errors.Wrapf(err, "Error in starting pod")
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

		if detach {
			writer := printers.GetNewTabWriter(os.Stdout)
			defer writer.Flush()
			printer := cmd.GetPrinter(clicontext)
			return printer.PrintPodDetails(result, writer)
		}

		attachContainerID, err := findRunningContainerID(result, name)
		if err != nil {
			return errors.Wrapf(err, "Cannot attach to container")
		}

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

		// Stop updating ui lines, let the std piping take the terminal
		ui.Stop()
		defer ui.Start()

		return term.Safe(func() error {
			return client.Attach(attachContainerID, api.NewAttachIO(term.In, term.Out, stderr), hooks...)
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
