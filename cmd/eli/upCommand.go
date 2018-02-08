package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	"github.com/ernoaapa/eliot/pkg/progress"
	"github.com/ernoaapa/eliot/pkg/resolve"
	"github.com/ernoaapa/eliot/pkg/sync"
	"github.com/ernoaapa/eliot/pkg/term"
	"github.com/ernoaapa/eliot/pkg/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var upCommand = cli.Command{
	Name:        "up",
	HelpName:    "up",
	Usage:       "Start development session in the device",
	Description: "With up command, you can start new development session in the device",
	UsageText: `eli up [options] -- <command>

	 # Run code in current directory in the device
	 eli up

	 # Run 'build.sh' command in device with files in current directory
	 eli up -- ./build.sh
	 
	 # Run container image in the container
	 eli up --image docker.io/eaapa/hello-world:latest

	 # Run container with name in the device
	 eli up --image docker.io/eaapa/hello-world:latest --name my-pod
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
		cli.StringSliceFlag{
			Name:  "sync",
			Usage: "Directory to sync to the target container (default: current directory)",
		},
		cli.StringFlag{
			Name:  "workdir, w",
			Usage: "Working directory inside the container",
		},
	},
	Action: func(clicontext *cli.Context) (err error) {

		var (
			projectConfig = config.ReadProjectConfig("./.eliot.yml")
			name          = cmd.First(clicontext.String("name"), projectConfig.Name, filepath.Base(cmd.GetCurrentDirectory()))
			image         = cmd.First(clicontext.String("image"), projectConfig.Image)
			rm            = clicontext.Bool("rm")
			tty           = clicontext.Bool("tty")
			env           = projectConfig.EnvWith(clicontext.StringSlice("env"))
			workdir       = cmd.First(clicontext.String("workdir"), projectConfig.WorkDir)
			syncs         = cmd.MustParseSyncs(append(projectConfig.Syncs, clicontext.StringSlice("sync")...))
			mounts        = cmd.MustParseMounts(append(projectConfig.Mounts, clicontext.StringSlice("mount")...))
			binds         = cmd.MustParseBinds(append(projectConfig.Binds, clicontext.StringSlice("bind")...))
			args          = append(projectConfig.Command, cmd.DropDoubleDash(clicontext.Args())...)

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
			log := ui.NewLine().Loading("Resolve image for the project...")
			info, err := client.GetInfo()
			if err != nil {
				log.Fatalf("Unable to resolve image for the project. Failed to get target device architecture: %s", err)
			}

			var projectType string
			projectType, image, err = resolve.Image(info.Arch, cmd.GetCurrentDirectory())
			if err != nil {
				log.Fatal("Unable to automatically resolve image for the project. You must define target container image with --image option")
			}
			log.Donef("Detected %s project, use image: %s (arch %s)", projectType, image, info.Arch)
		}

		image = utils.ExpandToFQIN(image)

		for _, variable := range env {
			if !model.IsValidEnvKeyValuePair(variable) {
				return fmt.Errorf("Invalid --env value [%s], must be in format KEY=value. E.g. --env FOO=bar", variable)
			}
		}

		opts := []api.PodOpts{}

		if len(syncs) > 0 && projectConfig.SyncContainer != nil {
			syncDestinations := []string{}
			mounts := []*containers.Mount{}

			for _, sync := range syncs {
				syncDestinations = append(syncDestinations, sync.Destination)
				mounts = append(mounts, cmd.MustParseBindFlag(fmt.Sprintf("/var/lib/volumes/%s/%s:%s:rw,rshared", name, strings.Replace(sync.Destination, "/", "_", -1), sync.Destination)))
			}

			projectConfig.SyncContainer.Env = append(projectConfig.SyncContainer.Env, fmt.Sprintf("VOLUMES=%s", strings.Join(syncDestinations, " ")))

			opts = append(opts, api.WithContainer(projectConfig.SyncContainer))

			opts = append(opts, api.WithSharedMount(mounts...))

			if workdir == "" && len(syncDestinations) == 1 {
				opts = append(opts, api.WithWorkingDir(syncDestinations[0]))
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
					{
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

		progressc := make(chan []*progress.ImageFetch)
		go cmd.ShowDownloadProgress(progressc)

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
				log := ui.NewLine().Loadingf("Delete pod %s", pod.Metadata.Name)
				_, err := client.DeletePod(pod)
				if err != nil {
					log.Errorf("Error while deleting pod [%s]: %s", pod.Metadata.Name, err)
				} else {
					log.Donef("Deleted pod [%s]", pod.Metadata.Name)
				}
			}()
		}

		attachContainerID, err := cmd.FindRunningContainerID(result, name)
		if err != nil {
			return errors.Wrapf(err, "Cannot attach to container")
		}

		hooks := []api.AttachHooks{}
		if len(syncs) > 0 {
			hooks = append(hooks, func(endpoint config.Endpoint, done <-chan struct{}) {
				sync.StartRsync(done, endpoint.GetHost(), 873, syncs, 1*time.Second)
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
			defer cmd.StopCatch(sigc)
		}

		// Stop updating ui lines, let the std piping take the terminal
		ui.Stop()
		defer ui.Start()

		return term.Safe(func() error {
			return client.Attach(attachContainerID, api.NewAttachIO(term.In, term.Out, stderr), hooks...)
		})
	},
}
