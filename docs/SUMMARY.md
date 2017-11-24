# Summary

* [x Introduction](README.md)
  * [x Motivation]()
  * [Use Cases]()
  * [Eliot vs. Other]()
* [x Getting Started]()
  * [x Core Concepts]()
  * [x Installation]()
  * [x Deploy first app]()
  * [x Development in device]()
* [Installation]()
  * [Manual installation]()
* [x eli client]()
  * [eli get devices]()
  * [eli create]()
  * [eli get pods]()
  * [eli describe pod]()
  * [eli delete pod]()
  * [eli attach]()
  * [eli run]()
* [x Configuration]()
  * [Yaml Pod Specification]()
  * [Project Configuration]()
* [EliotOS]()
* [Contributing]()
 * [Getting Started]()
 * [Debugging]()


## Introduction

Eliot is a open source system for managing containerized applications on top of the IoT device with an emphasis to usability, simplicity and security. Eliot gives simplified app delivery, isolation and additional security to connected device solutions.

In consideration of connected device limitations like unstable connection, limited computing resources, hardware connectivity, Eliot connects devices to single easy to use platform where you can manage devices and applications easily and safely.

Cloud Native technologies like Docker and Kubernetes have inspired heavily and if you're familiar with those, you find really easy to get started with Eliot.

### Motivation

I was building modern connected device product what customers are located around the world. I have over 10 years of software engineer experience with five years of DevOps and faced problem that there's no state-of-the-art solution for managing connected devices a way what is common in nowadays in cloud solutions. 
Most platforms and services focus heavily to the cloud connectivity, data processing and analysis, but I needed a solution to manage device Operating System and application deployment to build easy to use, modern service for our customers.

Key features needed:
- Quick realtime development
- Simple and fast application delivery
- Over-The-Air device management
- Resource allocation and restriction
- Security and software isolation
- Inter-process discovery and communication
- Take into account IoT limitations from ground up

And that's the day when Eliot were born ❤︎

### Eliot vs. Other

#### Docker



#### Kubernetes

Kubernetes is great platform to orchestrate containerized software in cloud environment with an emphasis to .


#### AWS, Azure, Google, IBM
All cloud IoT solutions base to the same practice, you use SDK to implement software what collects data from sensor and send it to the cloud where data gets processed and analysed. Analysis result can send message back to the device to trigger some action.

Eliot don't try to provide this kind of features at all, actually you can use any cloud service with Eliot.
Eliot provides a easy way to deliver your cloud integration to the device and gives you a way to update the software across thousands of devices safely and easily.

Even better, might be that you don't need to code anything! There might be available open source implementation made by someone in Docker community or you can share your code to the thousands of Docker users around the world with single command.

#### Ansible,Chef,SaltStack

Not for managing thousands or hundred of thousands of devices
Need direct connect
Not meant for devices what connects "harvoin"


## Getting Started
In this guide you will learn basic concepts and learn how you can deploy containers to the device with Eliot.

### Core Concepts
To work with Eliot, you describe the _desired state_ and Eliot Agent works to make the device _current state_ match with _desired state_.

With Eliot, you run all your processes in containers. Containers gives easy software delivery, network, disk and resources isolation.
You can use existing tooling, like [Docker command line tool](), to build container images and any Container Registry, for example [Docker public hub]().

You define what software containers you would like to run in the device, by deploying [Pods]() in to the device.

#### Pods

Pod is a basic building block of Eliot and is the main model object what you deploy and manage.
If you're familiar with Kubernetes definition of _Pod_, in Eliot it's exactly same except some features are not yet implemented.

A Pod wraps application container (or, in some cases, multiple containers) and gives configuration options like host network, restart policy, etc. All containers in same Pod are meant to get started, deployed, and deleted together.

### Installation

Easiest way to test run Eliot is to run EliotOS in RaspberryPI 3b.

EliotOS is minimal Linux Operating System which contains only minimal components to run Eliot which are Linux kernel, runc, containerd and `eliotd` daemon. Check the [EliotOS]() section for more info.

1. Install [Etcher CLI](https://etcher.io/cli/)
2. Download EliotOS image
3. Flash the EliotOS image to SD card
```shell
etcher eliotos.img
```
4. Connect RaspberryPI to network and power on!

For more details and other installation options, see the [Installation]() section.

### Deploy first app
When you have connected RaspberryPI to network and power it on, first step is to check that Eliot can discover and connect to the device.
```shell
eli get devices
```
List should contain your RaspberryPI device.

There should not be anything running, you can verify it by listing all created Pods in the device
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

This was the quick start how to deploy containers to the Eliot. Next step is to learn [how you can develop your software in the device in real time]()


### Development in device
When developing IoT solution, you usually need to develop software what reads data from hardware sensor and send it to the cloud. For this you need to have access to the device, connect to the sensor and be able to develop your software.
For this Eliot offer really easy way, `run` -command.

With `run` command you can connect to the device, start container what contains all required development tools and synchronize your local files to the device in real time.

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

## `eli` client
To see full documentation about commands and options, type `eli --help` and `eli <commnand> --help`.

### `eli get devices`
Eliot can search devices automatically from network with mDNS protocol. You can get list of devices what Eliot finds with `get devices` command.

```shell
eli get devices
```

### `eli create`
It's a good practice to store Pod definitions in version control and deploy exactly same deployment to each device.
You can write definition in `yaml` file which follows the [yaml specification]() and use `create` command to create all resources.

```shell
eli create -f <file.yml>
```

### `eli get pods`
You can get list of all running Pods with `get pods`.

```shell
eli get pods
```

If you have multiple devices, you get list of all Pods in all devices. If you want to get list of Pods in specific device, give `--device` flag.

```shell
eli get pods --device <device name>
```

### `eli describe pod`
To view Pod details like container image(s), statuses, etc., use command `describe pod <pod name>`.

```shell
eli describe pod <pod name>
```

### `eli delete pod`
To stop and clean up Pod from device give Pod name to `delete pod <pod name>` command.

```shell
eli delete pod <pod name>
```
After this, Eliot will stop and remove all container(s) from the device and free the used resources.

### `eli attach`
Sometimes you want to view output of your process, you can give pod name to the `attach` command.

```shell
eli attach <pod name>
```

If Pod contains multiple containers, you must pass containerID with `--container` flag.

```shell
eli attach --container <containerID> <pod name>
```

### `eli run`
When you develop your software, often you need to have access to the device to connect your software to some hardware sensor. To make development as easy as possible, Eliot have `run` command.

`run` command will: 
1. Detect what type of project you have in current directory and selects container image (you can use `--image` flag to override image)
2. Starts required containers
3. Mounts `/dev` directory into the container
4. Run `/bin/bash` session in container
5. Move your terminal session to the container
6. Start syncing your local files to the container

```shell
eli run --bind /dev:/dev /bin/bash
```

## Configuration
###  Yaml Pod Specification
When you need to deploy same Pods to multiple devices, you don't want to run multiple `eli create` commands with `--image` etc. flags. Easier way is to create Yaml file which describes all Pods and run `create -f <file.yml>

Here's example yaml specification
```yml
metadata:
  name: "hello-world"
spec:
  containers:
    - name: "hello-world"
      image: "docker.io/eaapa/hello-world:latest"
```

You can find more examples from `examples` directory.

###  Project Configuration
If you use `run` command to develop your software project in the device, you probably have specific container image, common bindings and other configurations and you don't want to define all of them with `eli run` flags. For this you can create `.eli.yml` file in to the root of your project and define configurations in there.

```yml
name: some-custom-name
image: someother/image:versin
sync:
  target: /go/src/github.com/ernoaapa/eliot
binds:
  - /dev:/dev
```

## EliotOS
EliotOS is minimal Linux Operating System, built with [Linuxkit](https://github.com/linuxkit/linuxkit), which contains only minimal components to run Eliot.

- `Linux kernel`
- `containerd` - Container runtime
- `runc` - Run containers based on OCI specification
- `eliotd` - Primary node agent which manages containers based on specs

There's more info in [project GitHub page](https://github.com/ernoaapa/eliot-os) and pre-built images available in [releases page](https://github.com/ernoaapa/eliot-os/releases).