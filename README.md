
![Eliot](docs/eliot-logo-medium.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/ernoaapa/eliot)](https://goreportcard.com/report/github.com/ernoaapa/eliot)[![Go Report Card](https://img.shields.io/badge/deb-packagecloud.io-844fec.svg)](https://packagecloud.io/ernoaapa/eliot)

> This is early alpha version! There's not all features yet implemented, not heavily tested with different devices and code might get large breaking changes until the first release.

Eliot is a open source system for managing containerized applications on top of the IoT device with an emphasis to usability, simplicity, security and stability. Eliot gives simplified app delivery, isolation and additional security to traditional installations.

[![asciicast](https://asciinema.org/a/vZcVZKEfAosSSrhWrJbmIqAd9.png)](https://asciinema.org/a/vZcVZKEfAosSSrhWrJbmIqAd9?autoplay=1&speed=2&t=4)

Docker and Kubernetes have inspired heavily and if you're familiar with those, you find really easy to get started with Eliot.

<sub>Built with ❤︎ by [Erno Aapa](https://github.com/ernoaapa) and [contributors](https://github.com/ernoaapa/eliot/contributors)</sub>

## Usage

- [Documentation](http://docs.eliot.run)
- [Binary releases](https://github.com/ernoaapa/eliot/releases)
- [Docker releases](https://hub.docker.com/r/ernoaapa/eliotd/tags)

Eliot is based on top of the [containerd](https://github.com/containerd/containerd) to provide simple, _Kubernetes like_ API for managing containers.

Eliot is built from following components
- `eli` - Command line tool for managing the device
- `eliotd` - Daemon for the device to manage containers

### Features
- Manage running containers in the device
- Attach to container process remotely for debugging
- Fast _develop-in-device_ development start

[Let us know](https://github.com/ernoaapa/eliot/issues/new) what would be the next awesome feature :)

## Getting started
[See the documentation](http://docs.eliot.run/getting_started.html) how to get started with Eliot.

**Rest of this document is about developing Eliot itself, not how to develop on top of the Eliot.**

## Development

### Prerequisites
- [Install Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Install Golang 1.11](https://golang.org/doc/install)
- [Install Docker](https://docs.docker.com/install/)
- [Install Linuxkit](https://github.com/linuxkit/linuxkit#build-the-linuxkit-tool)
- [Install goreleaser](https://goreleaser.com/#introduction.installing_goreleaser) (for building `eliotd`)
- Get Eliot source code `git clone https://github.com/ernoaapa/eliot && cd eliot`

### Developing `eli` cli
If you're making changes to the `eli` command line tool, you can just build and run the command
```shell
go run ./cmd/eli/* get nodes
```

### Developing `eliotd` daemon
To develop `eliotd` there's two different ways; latter is not tested
- run `eliotd` in EliotOS with Linuxkit
- run `eliotd` daemon locally

#### Run EliotOS locally
For development purpose, you can build and run the [EliotOS](https://github.com/ernoaapa/eliot-os) locally, but keep in mind that the **environment is amd64 not arm64** so container images what work in this environment might not work in RaspberryPI if the images are not multi-arch images.

1. Build `eliotd` binary and Docker image
   - `goreleaser --snapshot --rm-dist`
2. Get EliotOS linuxkit configuration
   - `curl https://raw.githubusercontent.com/ernoaapa/eliot-os/master/rpi3.yml > rpi3.yml`
3. Update `rpi3.yml`
   - Check from `goreleaser` the `amd64` container image name
   - Edit the `rpi3.yml` and update the `eliotd` image tag to match with the previous value
4. Build EliotOS image:
   - `linuxkit build rpi3.yml`
5. Start image:
   - MacOS: `sudo linuxkit run hyperkit -cpus 1 -mem "1048" -disk size=10G -networking vmnet moby`
6. Test connection
   - `eli get nodes`

#### Run `eliotd` locally
This is not tested, but should go roughly like this:

1. Install [runc](https://github.com/opencontainers/runc)
2. Install [containerd](https://github.com/containerd/containerd/blob/master/docs/getting-started.md#starting-containerd)
3. Run `go run ./cmd/eliotd/* --debug --grpc-api-listen 0.0.0.0:5000`
