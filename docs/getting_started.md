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

## Deploy first app
Assuming that you already complete installation of `eli` client and `eliot` agent. If not, check the[installation guide](installation.md) and come back after finishing it.  

When you have connected RaspberryPI to network and power it on, first step is to check that Eliot can discover and connect to the device.
```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli get devices]
  ✓ Discovered 1 devices from network

HOSTNAME                       ENDPOINT
linuxkit-96165e7f48d7.local.   192.168.64.79:5000
```
List should contain your RaspberryPI device.

There should not be anything running, you can verify it by listing all created _Pods_ in the device
```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli get pods]
  ✓ Discovered 1 device(s) from network
  • Connect to linuxkit-96165e7f48d7.local. (192.168.64.79:5000)

NAMESPACE   NAME          CONTAINERS   STATUS
```
Pod listing should be empty.

Now let's deploy [eaapa/hello-world](https://hub.docker.com/eaapa/hello-world) image, what just prints _Hello World_ text.
```shell
eli create pod --image eaapa/hello-world testing
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli create pod --image eaapa/hello-world testing]
  ✓ Discovered 1 device(s) from network
  • Connect to linuxkit-96165e7f48d7.local. (192.168.64.79:5000)
  ⠸ Download docker.io/eaapa/hello-world:latest
Name:             testing
Namespace:        eliot
Device:           linuxkit-96165e7f48d7
State:            running(1)
Restart Policy:   always
Host Network:     false
Host PID:         false
Containers:
          eaapa-hello-world:
                    Image:           docker.io/eaapa/hello-world:latest
                    ContainerID:     b97cq4r744405e9hmscg
                    State:           running
                    Restart Count:   0
                    Working Dir:     /
                    Args:
                              - /bin/sh
                    Env:
                              - PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
                              - PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
                    Mounts:
                              - type=proc,source=proc,destination=/proc,options=
                              - type=tmpfs,source=tmpfs,destination=/dev,options=nosuid:strictatime:mode=755:size=65536k
                              - type=devpts,source=devpts,destination=/dev/pts,options=nosuid:noexec:newinstance:ptmxmode=0666:mode=0620:gid=5
                              - type=tmpfs,source=shm,destination=/dev/shm,options=nosuid:noexec:nodev:mode=1777:size=65536k
                              - type=mqueue,source=mqueue,destination=/dev/mqueue,options=nosuid:noexec:nodev
                              - type=sysfs,source=sysfs,destination=/sys,options=nosuid:noexec:nodev:ro
                              - type=tmpfs,source=tmpfs,destination=/run,options=nosuid:strictatime:mode=755:size=65536k
```

When create completes, you can check the output by attaching to the container
```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli attach testing]
Hello world!
Hello world!
Hello world!
^C
```
You should see now _Hello World_ text printing once in a second to the terminal. Hit `CTRL + C` to close the connection.

Next you want to remove the Pod from the device.
```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli delete pod testing]
  ✓ Discovered 1 device(s) from network
  • Connect to linuxkit-96165e7f48d7.local. (192.168.64.79:5000)
  ✓ Fetched pods
  ✓ Deleted pod testing
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
**[terminal]
**[prompt ernoaapa@mac]**[path ~/go/src/github.com/ernoaapa/eliot]**[delimiter  $ ]**[command eli run --bind /dev:/dev /bin/bash]
root@linuxkit-96165e7f48d7:/go/src/github.com/ernoaapa/eliot#
```

Now you should be in the container and you can make changes in your local computer, compile and run your project in the device. 
If you type `exit` your terminal session comes back to your local computer and Eliot removes all created containers.

When you're done with development, you can create container image with your preferred tools (e.g. `docker`) and use the `create` command to deploy your software to the device permanently.

Easy right! :)
