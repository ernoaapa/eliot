# Installation
> At the moment Eliot is tested only on RaspberryPI 3b with [EliotOS](eliotos.md) but installing it to other environments should be pretty straight forward. If you happen to test in some other device, [please let us know!](https://github.com/ernoaapa/eliot/issues/new)

## Install Eliot client
Eliot client, called `eli`, is just a binary what you can download from [Eliot releases](https://github.com/ernoaapa/eliot/releases).

### MacOS
1. `brew install ernoaapa/eliot/eli`
2. Test `eli --version`

### Linux
1. Download `eli` binary from [releases](https://github.com/ernoaapa/eliot/releases)
2. Place it into your $PATH
3. Test `eli --version`

## Install EliotOS
By far the easiest and most secure way to use Eliot is by using [EliotOS](eliotos.md). EliotOS is minimal Linux Operating System, built with [linuxkit](https://github.com/linuxkit/linuxkit), which contains only minimal components to run Eliot which are Linux kernel, `runc`, `containerd` and `eliotd` daemon. Check the [EliotOS](eliotos.md) section for more info.


### RaspberryPI 3
`eli` provides [build command](client.md#eli-build-device) to build [EliotOS](eliotos.md) for RaspberryPI3 and install it to the sdcard.

1. Format sdcard as you would normally
2. Mount it to for example `/Volumes/rpi3`
3. Build image and unpack it to the directory `eli build device | tar xv -C /Volumes/rpi3`
4. Unmount the disk
5. Connect RaspberryPI with ethernet cable to same network with your laptop and power on!
6. In less than 10s you should see the device with command `eli get devices`
7. And that's it! ☺

Next step is to follow [getting started guide](getting_started.md#deploy-first-app) and deploy first app!

## Manual installation
- Install and run [runc](https://github.com/opencontainers/runc)
- Install and run [containerd](https://github.com/containerd/containerd)
- Install and run [eliotd](https://github.com/ernoaapa/eliot)
