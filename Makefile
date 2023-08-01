# Set the Go application binary name
BINARY_PATH := ./bin/ub

# Set the Go package path
PACKAGE_PATH := ./

# Set Go linker flags
LDFLAGS := -s -w -buildid=

# Set Go build options
GO_BUILD_OPTIONS := -trimpath -buildvcs=false -ldflags "$(LDFLAGS)"

# Set the Go proxy
GOPROXY := https://proxy.golang.org

# Set the GoOS to linux for the Docker build
GOOS := linux

# Set the CGO_ENABLED to 0 to disable CGO
CGO_ENABLED := 0

.PHONY: install
## Install Go dependencies.
install:
	# Download and verify Go module dependencies.
	go mod download
	go mod verify

.PHONY: run
## Run the application in development mode.
run:
	go run . discord:run --debug

.PHONY: build
## Build the Go application.
build:
	CGO_ENABLED=$(CGO_ENABLED) GOPROXY=$(GOPROXY) GOOS=$(GOOS) go build $(GO_BUILD_OPTIONS) -o $(BINARY_PATH) $(PACKAGE_PATH)

.PHONY: clean
## Clean the build artifacts.
clean:
	# Remove the binary file generated from the build process.
	rm -f $(BINARY_PATH)

.PHONY: docker-redeploy
docker-redeploy:
	docker compose down --remove-orphans
	docker compose --env-file ./.env up -d --build
	sleep 5
	docker compose ps
	docker stats --no-stream
	docker compose logs
