# Installation
> At the moment Eliot is tested only on RaspberryPI 3b, but installing it to other environments should be pretty straight forward. If you happen to test in some other device, [please let us know!](https://github.com/ernoaapa/eliot/issues/new)

## Install CLI
Eliot client, called `eli`, is just a binary what you can download from [Eliot releases](https://github.com/ernoaapa/eliot/releases).

### MacOS
1. `brew install ernoaapa/eliot/eli`
2. Test `eli --version`

### Linux
1. Download `eli` binary from [releases](https://github.com/ernoaapa/eliot/releases)
2. Place it into your $PATH
3. Test `eli --version`

## Install device
There's three options for device installation

1. [Use EliotOS on RaspberryPI 3](/installation.html#eliotos-on-raspberrypi3)
2. [Use Debian linux (e.g. Raspbian) and use deb packages](/installation.html#debian-raspbian-installation)
3. [Manual installation to any linux](/installation.html#manual-installation)

### EliotOS on RaspberryPI3
By far the easiest and most secure way to use Eliot is by using [EliotOS](eliotos.md). EliotOS is minimal Linux Operating System, built with [linuxkit](https://github.com/linuxkit/linuxkit), which contains only minimal components to run Eliot which are Linux kernel, `runc`, `containerd` and `eliotd` daemon. Check the [EliotOS](eliotos.md) section for more info.

`eli` provides [build command](client.md#eli-build-device) to build [EliotOS](eliotos.md) for RaspberryPI3 and install it to the sdcard.

1. Format sdcard as you would normally
2. Mount it to for example `/Volumes/rpi3`
3. Build image and unpack it to the directory `eli build device | tar xv -C /Volumes/rpi3`
4. Unmount the disk
5. Connect RaspberryPI with ethernet cable to same network with your laptop and power on!
6. In less than 10s you should see the device with command `eli get devices`
7. And that's it! ☺

Next step is to follow [getting started guide](getting_started.md#deploy-first-app) and deploy first app!

### Debian (Raspbian) installation
Eliot provides deb packages through packagecloud for Debian linux, for example Raspbian, to install Eliot and all dependencies.

<a href="https://packagecloud.io/ernoaapa/eliot"><img height="46" width="158" alt="Eliot DEB Repository · packagecloud" src="https://packagecloud.io/images/packagecloud-badge.png" /></a>

```shell
# Install eliot deb repository
curl -s https://packagecloud.io/install/repositories/ernoaapa/eliot/script.deb.sh | sudo bash

# Install Eliot and dependencies
apt-get update && apt-get install -y eliot

# Start the services
systemctl start containerd && systemctl enable containerd
systemctl start eliotd && systemctl enable eliotd
```

That's it! Now try running `eli get nodes` and you should see your device!

### Manual installation
- Build and install [runc](https://github.com/opencontainers/runc)
- Build and install and run [containerd](https://github.com/containerd/containerd)
- Build and install and run [eliotd](https://github.com/ernoaapa/eliot)
