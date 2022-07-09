# TCP server with challenge-request DDoS protection

## 1. Prerequistites
- go
- docker
- makefile

## 2. Running

### Server:

```sh
make run
```

### Client (to test it out)
```sh
make run-client
```

### Running server in docker
```
$ make build-image
$ make run-container
```