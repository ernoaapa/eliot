# `eli` client
To see full documentation about commands and options, type `eli --help` and `eli <command> --help`.

## `eli get devices`
Eliot can search devices automatically from network with mDNS protocol.

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli get devices]
  ✓ Discovered 1 devices from network

HOSTNAME                       ENDPOINT
linuxkit-96165e7f48d7.local.   192.168.64.79:5000
```

If you have multiple devices, you need to give `--device` flag for all other commands to specify the target device.

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli --device linuxkit-96165e7f48d7.local. get pods]
```

## `eli run [-i -t] <image> [command]`
Like `docker run`, `eli run` start container, but start it in the device, not in your local computer.
With `run` command you can quickly run some container in the device, and after you complete, (by default) eliot removes the container and leaves the device clean.

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~/go/src/github.com/ernoaapa/eliot]**[delimiter  $ ]**[command eli run alpine -- /bin/sh]
root@linuxkit-96165e7f48d7:/# uname -a
Linux raspberrypi-e2ccbe63f23d 4.9.72-linuxkit #1 SMP Thu Dec 28 19:08:26 UTC 2017 x86_64 Linux
root@linuxkit-96165e7f48d7:/# exit
```

## `eli up -- <command>`
When you develop your software, often you need to have access to the device to read some hardware sensor from your software. Ideal place for development would be in the device, but it's always too slow and clumsy way to code software. 
To make development as easy as possible, Eliot have `up` command.

`up` command will: 
1. Detect what type of project you have in current directory and selects container image (you can use `--image` flag to override image)
2. Start required containers
3. Run default command or `<command>` in container
4. Move your terminal session to the container
5. Start syncing your local files to the container

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~/go/src/github.com/ernoaapa/eliot]**[delimiter  $ ]**[command eli up -- /bin/bash]
  ✓ Detected golang project, use image: docker.io/library/golang:latest (arch amd64)
root@linuxkit-96165e7f48d7:/go/src/github.com/ernoaapa/eliot# ls -l
total 268
-rw-------  1 nobody 65533    304 Nov 20 17:53 Dockerfile.in
-rw-------  1 nobody 65533    305 Dec 31 08:07 Dockerfile.tmpl
-rw-------  1 nobody 65533   1048 Nov 15 06:13 LICENSE
-rw-------  1 nobody 65533   6549 Jan  2 05:09 Makefile
-rw-------  1 nobody 65533   1689 Dec 31 03:28 README.md
drwxr-xr-x  3 nobody 65533   4096 Jan  5 01:16 _book
drwxr-xr-x  5 nobody 65533   4096 Jan  5 01:16 bin
-rw-------  1 nobody 65533    327 Jan  5 00:48 book.json
drwxr-xr-x  2 nobody 65533   4096 Dec 19 02:49 build
drwxr-xr-x  4 nobody 65533   4096 Jan  5 01:16 cmd
drwxr-xr-x  2 nobody 65533   4096 Jan  4 17:48 docs
drwxr-xr-x  2 nobody 65533   4096 Jan  5 00:53 examples
drwxr-xr-x 42 nobody 65533   4096 Jan  5 01:16 node_modules
-rw-------  1 nobody 65533 209760 Jan  5 00:47 package-lock.json
-rw-------  1 nobody 65533   2015 Jan  4 07:27 vendor.conf
root@linuxkit-96165e7f48d7:/go/src/github.com/ernoaapa/eliot# exit
  ✓ Deleted pod eliot
```

You can override defaults with flags (see `eli up --help`) or you can create `.eliot.yml` project configuration. See [configuration](configuration.md#project-configuration) for more info.

## `eli create -f <file.yml>`
It's a good practice to store _Pod_ definitions in version control and deploy exactly same deployment to each device.
You can write definition in `yaml` file which follows the [yaml specification](configuration.md#pod-specification) and use `create` command to create all resources.

**pods.yml**
```yaml
metadata:
  name: "hello-world"
spec:
  containers:
    - name: "hello-world"
      image: "docker.io/eaapa/hello-world:latest"
```

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli create -f pods.yml]
  ✓ Discovered 1 device(s) from network
  • Connect to linuxkit-96165e7f48d7.local. (192.168.64.79:5000)
  ⠏ Download docker.io/eaapa/hello-world:latest [==================================================================>-]
Name:             hello-world
Namespace:        eliot
Device:           linuxkit-96165e7f48d7
State:            running(1)
Restart Policy:   always
Host Network:     false
Host PID:         false
Containers:
          hello-world:
                    Image:           docker.io/eaapa/hello-world:latest
                    ContainerID:     b97cqo3744405e9hmsd0
                    State:           running
                    Restart Count:   0
                    Working Dir:     /
                    Args:
                              - /print-hello-world.sh
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

## `eli create pod --image <image ref> <pod name>`
Sometimes you want to create a _Pod_ and making [yaml specification](configuration.md#pod-specification) is just overhead, you can use `eli create pod` to create a _Pod_ to the device.

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli create pod --image alpine testing]
  ✓ Discovered 1 device(s) from network
  • Connect to linuxkit-96165e7f48d7.local. (192.168.64.79:5000)
  ⠸ Download docker.io/library/alpine:latest
Name:             testing
Namespace:        eliot
Device:           linuxkit-96165e7f48d7
State:            running(1)
Restart Policy:   always
Host Network:     false
Host PID:         false
Containers:
          library-alpine:
                    Image:           docker.io/library/alpine:latest
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
You can have `--image` multiple times to add multiple containers into the _Pod_.

## `eli get pods`
You can get list of all running Pods with `get pods`.

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli get pods]
  ✓ Discovered 1 device(s) from network
  • Connect to linuxkit-96165e7f48d7.local. (192.168.64.79:5000)

NAMESPACE   NAME          CONTAINERS   STATUS
eliot       testing       1            running(1)
eliot       hello-world   1            running(1)
```

## `eli describe pod <pod name>`
To view _Pod_ details like container image(s), statuses, etc., use command `describe pod <pod name>`.

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli describe pod hello-world]
  ✓ Discovered 1 device(s) from network
  • Connect to linuxkit-96165e7f48d7.local. (192.168.64.79:5000)
Name:             hello-world
Namespace:        eliot
Device:           linuxkit-96165e7f48d7
State:            running(1)
Restart Policy:   always
Host Network:     false
Host PID:         false
Containers:
          hello-world:
                    Image:           docker.io/eaapa/hello-world:latest
                    ContainerID:     b97cqo3744405e9hmsd0
                    State:           running
                    Restart Count:   0
                    Working Dir:     /
                    Args:
                              - /print-hello-world.sh
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

## `eli delete pod <pod name>`
To stop and clean up _Pod_ from device give _Pod_ name to `delete pod <pod name>` command.

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli delete pod hello-world]
  ✓ Discovered 1 device(s) from network
  • Connect to linuxkit-96165e7f48d7.local. (192.168.64.79:5000)
  ✓ Fetched pods
  ✓ Deleted pod hello-world
```
After this, Eliot will stop and remove all container(s) from the device and free the used resources.

## `eli exec [--container id] <pod name> -- <command>`
Sometimes you want to execute command inside the container to for example to debug some problem.
If the _Pod_ contains multiple containers, you need to give target container id with `--container` flag.
> Note: separate the eli command from target command with `--`

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli exec testing -- date]
Fri Jan  5 01:03:45 UTC 2018
```

With `eli exec` you can also open terminal session and enter into the container:
```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli exec -i -t testing -- /bin/sh]
/ #
/ # ls -la
total 56
drwxr-xr-x    1 root     root          4096 Jan  5 00:55 .
drwxr-xr-x    1 root     root          4096 Jan  5 00:55 ..
drwxr-xr-x    2 root     root          4096 Dec  1 16:32 bin
drwxr-xr-x    5 root     root           340 Jan  5 00:55 dev
drwxr-xr-x   15 root     root          4096 Dec  1 16:32 etc
drwxr-xr-x    2 root     root          4096 Dec  1 16:32 home
drwxr-xr-x    5 root     root          4096 Dec  1 16:32 lib
drwxr-xr-x    5 root     root          4096 Dec  1 16:32 media
drwxr-xr-x    2 root     root          4096 Dec  1 16:32 mnt
dr-xr-xr-x  149 root     root             0 Jan  5 00:55 proc
drwx------    1 root     root          4096 Jan  5 01:06 root
drwxr-xr-x    2 root     root            40 Jan  5 00:55 run
drwxr-xr-x    2 root     root          4096 Dec  1 16:32 sbin
drwxr-xr-x    2 root     root          4096 Dec  1 16:32 srv
dr-xr-xr-x   13 root     root             0 Jan  5 00:55 sys
drwxrwxrwt    2 root     root          4096 Dec  1 16:32 tmp
drwxr-xr-x    7 root     root          4096 Dec  1 16:32 usr
drwxr-xr-x   11 root     root          4096 Dec  1 16:32 var
/ # exit
```
> Note: If you have minimal container, it might not include the /bin/sh and you get error `/bin/sh: no such file or directory`

## `eli attach [-i] [--container id] <pod name>`
Sometimes you want to hook up your current terminal session to the container process stdin/stdout.
If _Pod_ contains multiple containers, you must pass containerID with `--container` flag.

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli attach hello-world]
Hello world!
Hello world!
Hello world!
Hello world!
Hello world!
Hello world!
^C
```

You can also give `-i` flag to hook up your stdin into the container, but watch out, if you for example press ^C (ctrl+c) to exit, you actually send kill signal to the process in the container which will stop the container.

## `eli build device [--format img | tar]`
Easiest way to run Eliot in your device is to use [EliotOS](https://github.com/ernoaapa/eliot-os) which is minimal Operating System where's just minimal components installed to run Eliot and everything else run on top of the Eliot in containers.

With `eli build device` command you can build EliotOS image that you can just write to your device sdcard.

> Note: At the moment we support only RaspberryPi 3b, for other devices [see installation guide](installation.md)

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli build device > my-image.img]
```

[EliotOS](https://github.com/ernoaapa/eliot-os) is built with [Linuxkit]([EliotOS](https://github.com/ernoaapa/eliot-os) and you can view the Linuxkit configuration with `--dry-run` flag.

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli build device --dry-run]
  ✓ Resolved Linuxkit config!
  ✓ Resolved output: image.img!
kernel:
  image: linuxkit/kernel:4.9.72
  cmdline: "console=tty0 console=ttyS0 console=ttyAMA0"
init:
  - linuxkit/init:9250948d0de494df8a811edb3242b4584057cfe4
  - linuxkit/runc:abc3f292653e64a2fd488e9675ace19a55ec7023
  - linuxkit/containerd:e58a382c33bb509ba3e0e8170dfaa5a100504c5b
  - linuxkit/ca-certificates:de21b84d9b055ad9dcecc57965b654a7a24ef8e0
onboot:
  - name: sysctl
    image: linuxkit/sysctl:ce3bde5118a41092f1b7048c85d14fb35237ed45
  - name: netdev
    image: linuxkit/modprobe:1a192d168adadec47afa860e3fc874fbc2a823ff
    # https://github.com/linuxkit/linuxkit/blob/master/docs/platform-rpi3.md#networking
    command: ["modprobe", "smsc95xx"]
  - name: dhcpcd
    image: linuxkit/dhcpcd:0d59a6cc03412289ef4313f2491ec666c1715cc9
    # Halts until dhcpcd can resolve ip address
    command: ["/sbin/dhcpcd", "--nobackground", "-f", "/dhcpcd.conf", "-1"]
  - name: format
    image: linuxkit/format:e945016ec780a788a71dcddc81497d54d3b14bc7
  - name: mount-lib
    image: linuxkit/mount:b346ec277b7074e5c9986128a879c10a1d18742b
    command: ["/usr/bin/mountie", "/var/lib"]
    # Mount /var/log to the first found disk device
  - name: mount-log
    image: linuxkit/mount:b346ec277b7074e5c9986128a879c10a1d18742b
    command: ["/usr/bin/mountie", "/var/log"]

services:
  - name: getty
    image: linuxkit/getty:22e27189b6b354e1d5d38fc0536a5af3f2adb79f
    env:
    # Makes the terminal open without password prompt
     - INSECURE=true
  - name: ntpd
    image: linuxkit/openntpd:536e5947607c9e6a6771957c2ff817230cba0d3c
  - name: dhcpcd
    image: linuxkit/dhcpcd:0d59a6cc03412289ef4313f2491ec666c1715cc9

  - name: eliotd
    image: ernoaapa/eliotd:v0.2.2
    command: ["/eliotd", "--debug", "--grpc-api-listen", "0.0.0.0:5000"]
    capabilities:
      - all
    net: host
    pid: host
    runtime:
      mkdir: ["/var/lib/volumes"]
    binds:
      - /containers:/containers
      - /var/lib/volumes:/var/lib/volumes
      - /var/lib/containerd:/var/lib/containerd
      - /run/containerd:/run/containerd
      - /etc/resolv.conf:/etc/resolv.conf
      - /etc/machine-id:/etc/machine-id
      - /var/log:/var/log # To be able to serve default containers logs through api
      - /tmp:/tmp # To be able to read temporary fifo log files
files:
  - path: /etc/issue
    contents: "welcome to EliotOS"
  - path: /etc/machine-id
    contents: "todo-generate"
    mode: "0600"
trust:
  org:
    - linuxkit
```

If you want to customise the Linuxkit configuration before building

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli build device --dry-run > custom-linuxkit.yml]
# Edit the my-custom-linuxkit.yml -file...
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command eli build device custom-linuxkit.yml > custom-image.img]
```

#### Shell piping
The build command supports shell piping; you can pipe-in the Linuxkit config and pipe-out the result tar image, to some other command. This is really handy specially if you wan't to make updating the device easy.

For example, you want to:
- Change the hostname in Linuxkit config to include creation timestamp
- Build image for RaspberryPi3
- Unpack the package to sdcard in path `/Volumes/raspberrypi3`

The `custom-linuxkit.yml` includes:
```yaml
# ... snip ...

files:
  - path: /etc/hostname
    contents: MY-HOSTNAME
  - path: /etc/issue

# ... snip ...
```

```shell
**[terminal]
**[prompt ernoaapa@mac]**[path ~]**[delimiter  $ ]**[command sed -e "s/\MY-HOSTNAME/eliot-$(date +%s)/" custom-linuxkit.yml \
  | eli build device --format tar \
  | tar xv -C /Volumes/raspberrypi3]
```

There's plenty of options how you can template the linuxkit file, for example [envtpl](https://github.com/subfuzion/envtpl). Use your  favourite.

_Pretty handy, ain't it? :D_

#### How it works?
Linuxkit doesn't support building arm images on x86, but RaspberryPi is arm based computer.
For building images, Eliot hosts [Linuxkit build server](https://github.com/ernoaapa/linuxkit-server) and when you execute `eli build device`, it sends the config to `build.eliot.run` server, which builds the image on arm server and send it back as either disk image (.img) or tar package (.tar).

If you want to host and use your own build server, see the [Linuxkit build server documentation](https://github.com/ernoaapa/linuxkit-server) and pass `--build-server http://my-custom-build-server.com` flag to build the image in your own server.