# Development getting started
## Prerequisites
- Go 1.9.x

## Development
Depends on are you going to develop the CLI client or the API how you should proceed.
### canctl
Easiest way is to just run with `go run`
```
go run ./cmd/canctl/* <command>
```
### can-api
To run fully functioning `can-api`, you need filesystem access for example to create FIFO files for container logs.
You can develop some of the features by tunneling the `containerd` socket connection.

```
# leave open
ssh <user@ip> -L /run/containerd/containerd.sock:/run/containerd/containerd.sock

# In another window
go run ./cmd/can-api/* 
```

## Test
You can run tests inside container so you don't need to install all tools locally
```
make test
```