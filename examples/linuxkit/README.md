# Linuxkit
Example minimum distribution to run containers and `cand`

## Prerequisites
- Moby
- Linuxkit

## Build
```
moby build linuxkit.yml
```

## Run
```
sudo linuxkit run hyperkit -networking vmnet -ip 192.168.64.10 linuxkit
```