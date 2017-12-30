# Getting Started
In this guide you will learn basic concepts and learn how you can deploy containers to the device with Eliot.

## Core Concepts
To work with Eliot, you describe the _desired state_ and `eliotd` Eliot Agent works to make the device _current state_ match with _desired state_.

With Eliot, you run all your processes in containers. Containers gives easy software delivery, network, disk and resources isolation.
You can use existing tooling, like [Docker command line tool](https://docs.docker.com/engine/reference/commandline/cli/), to build container images and any Container Registry, for example [Docker public hub](https://hub.docker.com).

You define what software containers you would like to run in the device, by deploying [Pods](getting_started.md#pods) in to the device.

### Pods

Pod is a basic building block of Eliot and is the main model object what you deploy and manage.
If you're familiar with [Kubernetes definition of _Pod_](https://kubernetes.io/docs/concepts/workloads/pods/pod/), in Eliot it's exactly same except some features are not yet implemented.

A Pod wraps application container (or, in some cases, multiple containers) and gives configuration options like host network, restart policy, etc. All containers in same Pod are meant to get started, deployed, and deleted together.

## Installation

Easiest way to get started with Eliot is to run [EliotOS](eliotos.md) in RaspberryPI 3b.

[EliotOS](eliotos.md) is minimal Linux Operating System, built with [linuxkit](https://github.com/linuxkit/linuxkit), which contains only minimal components to run Eliot which are Linux kernel, `runc`, `containerd` and `eliotd` daemon. Check the [EliotOS](eliotos.md) section for more info.

1. Install [Etcher CLI](https://etcher.io/cli/)
1. Format sdcard as you would normally
2. Mount it to for example `/Volumes/rpi3`
3. Build image and unpack it to the directory `eli build device | tar xv -C /Volumes/rpi3`
4. Unmount the disk
5. Install sdcard to RaspberryPI, connect with ethernet cable to same network with your laptop and power on!
6. In less than 10s you should see the device with command `eli get devices`

For more details and other installation options, see the [Installation](installation.md) section.

## Deploy first app
When you have connected RaspberryPI to network and power it on, first step is to check that Eliot can discover and connect to the device.
```shell
eli get devices
```
List should contain your RaspberryPI device.

There should not be anything running, you can verify it by listing all created _Pods_ in the device
```shell
eli get pods
```
Pod listing should be empty.

Now let's deploy [eaapa/hello-world](https://hub.docker.com/eaapa/hello-world) image, what just prints _Hello World_ text.
```shell
eli create --image eaapa/hello-world
```

When create completes, you can check the output by attaching to the container.
> Note: we pass `--stdin=false` here because we don't want that when we hit ^C, we send kill signal to the container. We only want to connect stdout to our terminal session
```shell
eli attach hello-world
```
You should see now _Hello World_ text printing once in a second to the terminal. Hit `CTRL/CMD + C` to close the connection.

Next you want to remove the Pod from the device.
```shell
eli delete pod hello-world
```

Now the device is back to the initial state, clean and no any services running.

This was the quick start how to deploy containers to the Eliot. 
Next step is to learn [how you can develop your software in the device in real time](getting_started.md#development-in-device)


## Development in device
When developing IoT solution, you usually need to develop software what reads data from hardware sensor and send it to the cloud. For this you need to have access to the device, connect to the sensor and be able to develop your software.
For this Eliot offer really easy way, `eli run` -command.

With `eli run` command you can connect to the device, start container what contains all required development tools and synchronize your local files to the device in real time.

Following command will: 
1. Detect what type of project you have and selects container image
3. Mounts `/dev` directory into the container
4. Run `/bin/bash` session in container
5. Move your terminal session to the container
6. Start syncing your local files to the container

```shell
eli run --bind /dev:/dev /bin/bash
```

Now you should be in the container and you can make changes in your local computer, compile and run your project in the device. 
If you type `exit` your terminal session comes back to your local computer and Eliot removes all created containers.

When you're done with development, you can create container image with your preferred tools (e.g. `docker`) and use the `create` command to deploy your software to the device.

Easy right! :)
