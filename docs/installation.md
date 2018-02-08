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

`eli` provides [build command](client.md#eli-build-device) to build [EliotOS](eliotos.md) for RaspberryPI3 and install it to the sdcard. There's few different ways to do it, depending on your preferences.

#### a) Write image with Etcher
[Etcher](https://etcher.io/) is really handy desktop app for Win, Mac and Linux for writing img -files to sdcards.
1. Install [Etcher](https://etcher.io/)
2. Build image file `eli build device`
3. Start Etcher and select created `eliot-os.img` as the source and your sdcard drive as target

### b) Write iamge with `dd`
If you're for example writing the image in headless server and you don't want to use any gui, you can write the image with `dd`

**WARNING! With dd tool can overwrite any partition of your machine, for example your primary OS partition! Use only if you know what you do.**

1. Connect your sdcard and resolve what device is the card (in this example it's `disk3`)
2. Unmount the disk (but keep connected)
3. Build img-file `eli build device`
4. Write to card `sudo dd bs=1m if=eliot-os.img of=/dev/rdisk3 conv=sync`

#### c) Write disk manually
You can also build the disk manually, if none of above suite to you.

1. Format sdcard as FAT32 and name it for example `rpi3`
2. Mount it to for example `/Volumes/rpi3` (Mac does this automatically)
3. Build image and unpack it to the sdcard `eli build device --format tar | tar xv -C /Volumes/rpi3`
4. Unmount the disk

#### Boot up!
Now just connect RaspberryPI with ethernet cable to same network with your laptop and power on!
In less than 10s you should see the device with command `eli get devices`.

And that's it! ☺

Next step is to follow [getting started guide](getting_started.md#deploy-first-app) and deploy first app!

### Debian (Raspbian) installation
Eliot provides deb packages through packagecloud for Debian linux, for example Raspbian, to install Eliot and all dependencies.

<a href="https://packagecloud.io/ernoaapa/eliot"><img height="46" width="158" alt="Eliot DEB Repository · packagecloud" src="https://packagecloud.io/images/packagecloud-badge.png" /></a>

```shell
# Install eliot deb repository
curl -s https://packagecloud.io/install/repositories/ernoaapa/eliot/script.deb.sh | sudo bash

# Install Eliot and dependencies
sudo apt-get install -y eliot

# Start the services
sudo systemctl start containerd && sudo systemctl enable containerd
sudo systemctl start eliotd && sudo systemctl enable eliotd
```

That's it! Now try running `eli get nodes` in your local computer and you should see your device!

### Manual installation
- Build and install [runc](https://github.com/opencontainers/runc)
- Build and install and run [containerd](https://github.com/containerd/containerd)
- Build and install and run [eliotd](https://github.com/ernoaapa/eliot)
