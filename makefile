.PHONY: dependencies test container container-push linux darwin windows run clean help 

BINARY_NAME=carbide-api
CONTAINERTAG=rancher-government-carbide/$(BINARY_NAME)
SRC=./cmd
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

check: test lint

# Test the binary
test: $(BINARY_NAME)
	go test $(SRC)

# Run linters
lint:
	go vet $(SRC)
	staticcheck $(SRC)

# Build the container image
container:
	$(CLI) build -t $(CONTAINERTAG):$(VERSION) . && $(CLI) image tag $(CONTAINERTAG):$(VERSION) $(CONTAINERTAG):latest
	
# Push the binary
container-push: container
	$(CLI) push $(CONTAINER_NAME):$(COMMIT_HASH) && $(CLI) push $(CONTAINER_NAME):latest

# Ensure dependencies are available
dependencies:
	go mod tidy && go get -v -d ./...

# Clean the binary
clean:
	rm -f $(BINARY_NAME)

# Show help
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
