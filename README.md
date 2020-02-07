# go-librarios
Librarios API server

## Build
```bash
# install dependencies
$ make vendor

# build binary
$ make build
```

### Cross compile
```bash
# cross-compile for ARM64
$ make build-arm64

# cross-compile for linux x86_64
$ make build-linux

# cross-compile for osx x86_64
$ make build-osx

# cross-compile for windows x86_64
$ make build-windows
```

### Docker image
```bash
# create docker image
$ make docker-image

# push docker image
$ make docker-push
```
