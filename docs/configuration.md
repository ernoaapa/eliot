# Configuration

## Pod Specification
When you need to deploy same Pods to multiple devices, you don't want to run multiple `eli create` commands with `--image` etc. flags. Easier way is to create Yaml file which describes all Pods and run `create -f <file.yml>

Here's simple yaml specification
```yml
metadata:
  name: "hello-world"
spec:
  containers:
    - name: "hello-world"
      image: "docker.io/eaapa/hello-world:latest"
```

If your container needs to listen some host port, give access to the host network with `hostNetwork` flag.
```yml
metadata:
  name: "with-host-network"
spec:
  hostNetwork: true
  containers:
    - name: "with-host-network"
      image: "docker.io/arm64v8/alpine:latest"
```

You can find more examples from [examples](https://github.com/ernoaapa/eliot/tree/master/examples) directory.

## Project Configuration
If you use `run` command to develop your software project in the device, you probably have specific container image, common bindings and other configurations and you don't want to define all of them with `eli run` flags. For this you can create `.eliot.yml` file in to the root of your project and define configurations in there.

```yml
name: some-custom-name
image: someother/image:version
syncs:
  .:/go/src/github.com/ernoaapa/eliot
binds:
  - /dev:/dev
```
