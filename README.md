# Can
## Development
### Build
You can build binaries inside container so you don't need to install Go locally
```
make build
```

To build container
```
make container
```

### Test
You can run tests inside container so you don't need to install all tools locally
```
make test
```

### Run
Build container with `make container` and run (update the image name to match with latest):
```
docker run -it --rm -e MACHINE_ID=ernoaaapa ernoaapa/cand-amd64:3a56852-dirty --debug --labels foo=bar --manifest https://gist.githubusercontent.com/ernoaapa/9e0f8cc1945544182eaf9468fbb84ca8/raw/manifest.yaml
```

#### Run locally with loading manifest from file
```
go run ./cmd/cand/main.go --debug --labels foo=bar --manifest ./examples/hello-world.yml
```

#### Run locally with loading manifest from url
```
go run ./cmd/cand/main.go --debug --labels foo=bar --manifest https://gist.githubusercontent.com/ernoaapa/9e0f8cc1945544182eaf9468fbb84ca8/raw/manifest.yaml
```
