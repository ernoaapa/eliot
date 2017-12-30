# Installation
At the moment Eliot is tested only on RaspberryPI 3b with [EliotOS](eliotos.md) but installing it to other environments should be pretty straight forward.

If you happen to test in some other device, [please let us know!](https://github.com/ernoaapa/eliot/issues/new)

## Install Eliot-OS
By far the easiest way to use Eliot is by using [EliotOS](eliotos.md) and `eli` provides easy command to build and install it to sdcard.

> Note: This is tested only with RaspberryPI 3b
1. Format sdcard as you would normally
2. Mount it to for example `/Volumes/rpi3`
3. Build image and unpack it to the directory `eli build device | tar xv -C /Volumes/rpi3`
4. Unmount the disk
5. Connect RaspberryPI with ethernet cable to same network with your laptop and power on!
6. In less than 10s you should see the device with command `eli get devices`
7. Follow [Getting started](getting_started.md#deploy-first-app) to deploy first app

And that's it! ☺

## Manual installation

- Install [runc](https://github.com/opencontainers/runc)
- Install [containerd](https://github.com/containerd/containerd)
- Install [eliotd](https://github.com/ernoaapa/eliot)
