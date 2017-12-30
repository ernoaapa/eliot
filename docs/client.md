# `eli` client
To see full documentation about commands and options, type `eli --help` and `eli <command> --help`.

## `eli get devices`
Eliot can search devices automatically from network with mDNS protocol. You can get list of devices what Eliot finds with `get devices` command.

```shell
eli get devices
```

## `eli create`
It's a good practice to store Pod definitions in version control and deploy exactly same deployment to each device.
You can write definition in `yaml` file which follows the [yaml specification](configuration.md#pod-specification) and use `create` command to create all resources.

```shell
eli create -f <file.yml>
```

## `eli get pods`
You can get list of all running Pods with `get pods`.

```shell
eli get pods
```

If you have multiple devices, you get list of all Pods in all devices. If you want to get list of Pods in specific device, give `--device` flag.

```shell
eli get pods --device <device name>
```

## `eli describe pod`
To view Pod details like container image(s), statuses, etc., use command `describe pod <pod name>`.

```shell
eli describe pod <pod name>
```

## `eli delete pod`
To stop and clean up Pod from device give Pod name to `delete pod <pod name>` command.

```shell
eli delete pod <pod name>
```
After this, Eliot will stop and remove all container(s) from the device and free the used resources.

## `eli attach`
Sometimes you want to view output of your process, you can give pod name to the `attach` command.

```shell
eli attach <pod name>
```

If Pod contains multiple containers, you must pass containerID with `--container` flag.

```shell
eli attach --container <containerID> <pod name>
```

## `eli run`
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

## `eli build device`
Easiest way to run Eliot in your device is to use [EliotOS](https://github.com/ernoaapa/eliot-os) which is minimal Operating System where's just minimal components installed to run Eliot and everything else run on top of the Eliot in containers.

With `eli build device` command you can build EliotOS image what you can just unpack to your device sdcard.

> Note: At the moment we support only RaspberryPi 3b, for other devices [see installation guide](installation.md)

```shell
# Build default EliotOS
eli build device > my-image.tar

# You can view the underlying Linuxkit configuration
eli build device --dry-run
```

If you want to customise the Linuxkit configuration before building

```shell
eli build device --dry-run > custom-linuxkit.yml

# Edit the my-custom-linuxkit.yml -file...

# Build from the custom file
eli build device custom-linuxkit.yml > custom-image.tar
```

The build command supports shell piping; you can pipe-in the Linuxkit config and pipe-out the result image tar file, to some other command. This is really handy specially if you wan't to make updating the device easy.

#### Piping example
For example, you want to:
- Change the hostname in Linuxkit config to include creation timestamp
- Build image for RaspberryPi3
- Unpack the package to sdcard in path `/Volumes/raspberrypi3`

The `custom-linuxkit.yml` includes:
```
# ... snip ...

files:
  - path: /etc/hostname
    contents: MY-HOSTNAME
  - path: /etc/issue

# ... snip ...
```

```shell
# Tested on OS X...
sed -e "s/\MY-HOSTNAME/eliot-$(date +%s)/" custom-linuxkit.yml \
  | eli build device \
  | tar xv -C /Volumes/raspberrypi3
```

There's plenty of options how you can template the linuxkit file, for example [envtpl](https://github.com/subfuzion/envtpl). Use your  favourite.

_Pretty handy, ain't it? :D_

### How it works?
Linuxkit doesn't support building arm images on x86, but RaspberryPi is arm based computer.
For building images, Eliot hosts [Linuxkit build server](https://github.com/ernoaapa/linuxkit-server) and when you execute `eli build device`, it sends the config to `build.eliot.run` server, which builds the image on arm server and send it back as tar package.

If you want to host and use your own build server, see the [Linuxkit build server documentation](https://github.com/ernoaapa/linuxkit-server) and pass `--build-server http://my-custom-build-server.com` flag to build the image in your own server.