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
