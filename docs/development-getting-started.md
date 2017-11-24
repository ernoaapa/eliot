# Development getting started
## Prerequisites
- Go 1.9.x

## Development
Depends on are you going to develop the CLI client or the API how you should proceed.

### eli
Easiest way is to just run with `go run`
```
go run ./cmd/eli/* <command>
```

### eliotd
To run fully functioning `eliotd`, you need filesystem access for example to create FIFO files for container logs.
You can develop some of the features by tunneling the `containerd` socket connection from the device to your local computer.

```
# leave open
ssh <user@device-ip> -L /run/containerd/containerd.sock:/run/containerd/containerd.sock

# In another window
go run ./cmd/eliotd/* 
```

### Documentation
If you want to improve the documentation, you can run locally the gitbook 

## Test
You can run tests inside container so you don't need to install all tools locally
```
make test
```