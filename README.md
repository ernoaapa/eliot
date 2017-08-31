# Layery


## Development
### Build
You can build layery inside container so you don't need to install Go locally
```
make build
```

### Test
You can run tests inside container so you don't need to install all tools locally
```
make test
```

### Run locally with loading manifest from file
```
go run ./layeryd.go --debug --labels foo=bar --manifest ./examples/hello-world.yml
```

### Run locally with loading manifest from url
```
go run ./layeryd.go --debug --labels foo=bar --manifest https://gist.githubusercontent.com/ernoaapa/9e0f8cc1945544182eaf9468fbb84ca8/raw/manifest.yaml
```
