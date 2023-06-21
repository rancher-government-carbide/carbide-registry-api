.PHONY: dependencies test container container-push linux darwin windows run clean help 

BINARY_NAME=carbide-api
CONTAINERTAG=clanktron/$(BINARY_NAME)
SRC=./cmd
VERSION=0.1.0
GOENV=GOARCH=amd64 CGO_ENABLED=0
BUILD_FLAGS=-ldflags="-X 'main.Version=$(VERSION)'"
# change to docker if not using rancher desktop
CONTAINER_CLI=nerdctl

# Build the binary
$(BINARY_NAME):
	$(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(SRC)

# Clean the binary
clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-linux $(BINARY_NAME)-darwin $(BINARY_NAME)-windows

# Test the binary
test: $(BINARY_NAME)
	go test $(SRC)

# Build the binary for Linux
linux:
	GOOS=linux $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux $(SRC)
# Build the binary for MacOS
darwin:
	GOOS=darwin $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin $(SRC)
# Build the binary for Windows
windows:
	GOOS=windows $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows $(SRC)

# Build the container image
container:
	$(CONTAINER_CLI) build -t $(CONTAINERTAG):$(VERSION) . && $(CONTAINER_CLI) image tag $(CONTAINERTAG):$(VERSION) $(CONTAINERTAG):latest
	
# Push the binary
container-push: container
	$(CONTAINER_CLI) push $(CONTAINERTAG):$(VERSION) && $(CONTAINER_CLI) push $(CONTAINERTAG):latest 

# Ensure dependencies are available
dependencies:
	go mod tidy && go get -v -d ./...

# Show help
help:
	@printf "Available targets:\n"
	@printf "  $(BINARY_NAME) 		Build the binary (default)\n"
	@printf "  clean 		Clean build results\n"
	@printf "  test 			Build and test the binary\n"
	@printf "  linux 		Build the binary for Linux\n"
	@printf "  darwin 		Build the binary for MacOS\n"
	@printf "  windows 		Build the binary for Windows\n"
	@printf "  container 		Build the container\n"
	@printf "  container-push 	Build and push the container\n"
	@printf "  dependencies 		Ensure dependencies are available\n"
	@printf "  help 			Show help\n"
