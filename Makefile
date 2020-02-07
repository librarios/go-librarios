include .env
-include .env.override

BASE_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

ENV_GOMOD_ON=GO111MODULE=on
GOBUILD_OPT=-mod=vendor -v
GOTEST_OPT=-mod=vendor -v

BINARY=librarios
BINARY_WINDOWS=lebrarios.exe

all: build

# Build
build:
	@$(ENV_GOMOD_ON) $(GOBUILD) $(GOBUILD_OPT) -o $(BINARY)
build-arm64:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(ENV_GOMOD_ON) $(GOBUILD) $(GOBUILD_OPT) -o $(BINARY)
build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(ENV_GOMOD_ON) $(GOBUILD) $(GOBUILD_OPT) -o $(BINARY)
build-osx:
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(ENV_GOMOD_ON) $(GOBUILD) $(GOBUILD_OPT) -o $(BINARY)
build-windows:
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(ENV_GOMOD_ON) $(GOBUILD) $(GOBUILD_OPT) -o $(BINARY_WINDOWS)

# Clean
clean:
	@$(GOCLEAN)
	@rm -f $(BINARY)

# Run
run: build
	@./$(BINARY)

# Install dependencies to vendor/
vendor:
	@$(GOMOD) vendor
vendor-update:
	@$(GOGET) -u

# -----------------------
# Docker
# -----------------------
DOCKERFILE=Dockerfile
BUILD_DOCKERFILE=build-Dockerfile

# Build docker image
docker-image: build-linux
	@docker rmi $(DOCKER_IMAGE_NAME):$(VERSION)
	@DOCKER_BUILDKIT=0 docker build --no-cache -f $(DOCKERFILE) -t $(DOCKER_IMAGE_NAME):$(VERSION) .

# Build docker image within docker
docker-image-in-docker:
	@DOCKER_BUILDKIT=0 docker build --no-cache -f $(BUILD_DOCKERFILE) -t $(DOCKER_IMAGE_NAME):$(VERSION) .

docker-push:
	@docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):latest
	@docker push $(DOCKER_IMAGE_NAME):$(VERSION)
	@docker push $(DOCKER_IMAGE_NAME):latest
	@docker rmi $(DOCKER_IMAGE_NAME):latest

docker-rmi:
	@docker rmi $(DOCKER_IMAGE_NAME):$(VERSION)
