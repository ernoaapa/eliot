# hello-world
Simple container which just prints `Hello world!` into the stdout.

## How to build and use
> Replace `eaapa` with your Docker hub account!

```
# Build image
docker build -t eaapa/hello-world .
# Push to registry
docker push eaapa/hello-world

# ... update the deployment.json to have your username as image!

# Use the CLI to deploy image
go run ./cmd/can-cli/*.go --debug deployments create -f ./examples/hello-world/deployment.json

```