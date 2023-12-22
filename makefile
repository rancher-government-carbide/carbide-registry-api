BINARY_NAME=carbide-images-api
CONTAINERTAG=rancher-government-carbide/$(BINARY_NAME)
CONTAINERFILE=./Containerfile
SRC=./cmd
PKG=./pkg/*
VERSION=0.1.0
COMMIT_HASH=$(shell git rev-parse HEAD)
GOENV=CGO_ENABLED=0
BUILD_FLAGS=-ldflags="-X 'main.Version=$(VERSION)'"
TEST_FLAGS=-v -cover -count 1
# change to docker if not using rancher desktop
CLI=nerdctl

# Build the binary
$(BINARY_NAME):
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(SRC)

.PHONY: check
check: test lint

# Test the binary
.PHONY: test
test:
	go test -v $(SRC) $(PKG)

# Run linters
.PHONY: lint
lint:
	go vet $(SRC) $(PKG)
	staticcheck $(SRC) $(PKG)

# Build the container image
.PHONY: container
container:
	$(CLI) build -t $(CONTAINERTAG):$(COMMIT_HASH) -f $(CONTAINERFILE) . && $(CLI) image tag $(CONTAINERTAG):$(COMMIT_HASH) $(CONTAINERTAG):latest
	
# Push the binary
.PHONY: container-push
container-push: container
	$(CLI) push $(CONTAINER_NAME):$(COMMIT_HASH) && $(CLI) push $(CONTAINER_NAME):latest

# Ensure dependencies are available
.PHONY: dependencies
dependencies:
	go mod tidy && go get -v -d ./...

# Clean the binary
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

# Show help
.PHONY: help
help:
	@printf "Available targets:\n"
	@printf "  $(BINARY_NAME) 		Build the binary (default)\n"
	@printf "  test 			Build and test the binary\n"
	@printf "  lint 			Run go vet and staticcheck\n"
	@printf "  check 		Test and lint the binary\n"
	@printf "  container 		Build the container\n"
	@printf "  container-push 	Build and push the container\n"
	@printf "  dependencies 		Ensure dependencies are available\n"
	@printf "  clean 		Clean build results\n"
	@printf "  help 			Show help\n"
