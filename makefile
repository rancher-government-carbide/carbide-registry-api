BINARY_NAME=carbide-registry-api
ORG=rancher-government-carbide
CONTAINERNAME=$(ORG)/$(BINARY_NAME)
CONTAINERTAG=dev
CONTAINERFILE=./Containerfile
COMPILATION_SRC=./cmd
SRC=./cmd
VERSION=0.1.2
COMMIT_HASH=$(shell git rev-parse HEAD)
GOENV=CGO_ENABLED=0
BUILD_FLAGS=-ldflags="-X 'main.Version=$(VERSION)'"
TEST_FLAGS=-v -cover -count 1
ARTIFACT_DIR=dist
CLI=sudo nerdctl

$(BINARY_NAME):
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(SRC) ## Build binary (default)

.PHONY: check
check: test lint ## Test and lint

.PHONY: test
test: ## Run go tests
	go test $(TEST_FLAGS) ./...

.PHONY: lint
lint: ## Run go vet and staticcheck against codebase
	go vet ./...
	staticcheck ./...

.PHONY: build-container
build-container: clean ## Build the container
	$(CLI) build -t $(CONTAINERNAME):$(CONTAINERTAG) .
	
.PHONY: push-container
push-container: ## Push the container
	$(CLI) push $(CONTAINERNAME):$(CONTAINERTAG)

.PHONY: container-push
build-and-push-container: build-container push-container ## Build and push the container

.PHONY: dependencies
dependencies: ## Run go mod and go get to ensure dependencies
	go mod tidy && go get -v -d ./...

.PHONY: release
release: build-linux build-darwin build-windows package-chart checksums ## Build helm chart and binaries for all platforms

.PHONY: release-windows
build-windows: ## Build all arches for windows
	make GOOS=windows  GOARCH=amd64 BINARY_NAME=$(ARTIFACT_DIR)/$(BINARY_NAME)-windows-amd64-$(VERSION)
	make GOOS=windows GOARCH=arm64 BINARY_NAME=$(ARTIFACT_DIR)/$(BINARY_NAME)-windows-arm64-$(VERSION)

.PHONY: release-darwin
build-darwin: ## Build all arches for darwin
	make GOOS=darwin GOARCH=amd64 BINARY_NAME=$(ARTIFACT_DIR)/$(BINARY_NAME)-darwin-amd64-$(VERSION)
	make GOOS=darwin GOARCH=arm64 BINARY_NAME=$(ARTIFACT_DIR)/$(BINARY_NAME)-darwin-arm64-$(VERSION)

.PHONY: release-linux
build-linux: ## Build all arches for linux
	make GOOS=linux GOARCH=amd64 BINARY_NAME=$(ARTIFACT_DIR)/$(BINARY_NAME)-linux-amd64-$(VERSION)
	make GOOS=linux  GOARCH=arm64 BINARY_NAME=$(ARTIFACT_DIR)/$(BINARY_NAME)-linux-arm64-$(VERSION)
	make GOOS=linux  GOARCH=riscv64 BINARY_NAME=$(ARTIFACT_DIR)/$(BINARY_NAME)-linux-riscv64-$(VERSION)

.PHONY: release-chart
package-chart: ## Package helm chart
	helm package -u ./chart -d $(ARTIFACT_DIR)

.PHONY: checksums
checksums: ## Generate release asset checksums
	shasum -a 256 $(ARTIFACT_DIR)/* | tee $(ARTIFACT_DIR)/checksums.txt

.PHONY: clean
clean: ## Clean workspace
	rm -rf $(BINARY_NAME) $(ARTIFACT_DIR)/*

.PHONY: help
help:
	@echo "Available targets:"
	@if [ -t 1 ]; then \
		awk -F ':|##' '/^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST) | grep -v '^help:'; \
	else \
		awk -F ':|##' '/^[a-zA-Z0-9_-]+:.*?##/ { printf "  %-20s %s\n", $$1, $$NF }' $(MAKEFILE_LIST) | grep -v '^help:'; \
	fi
