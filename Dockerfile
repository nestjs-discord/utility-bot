# Build stage for Golang application
FROM golang:1.20.3-alpine as build-stage

# Set the working directory for COPY
WORKDIR /app

# Set environment variables
ENV CGO_ENABLED=0 \
    GOPROXY=https://proxy.golang.org \
    GOOS=linux

# Copy Go module files
COPY go.mod go.sum ./

# Download Go modules and verify
RUN go mod download && \
    go mod verify

# Copy application source code
COPY . .

# Build application
RUN go build -trimpath -buildvcs=false -ldflags "-w" -o /utility-bot ./main.go

# Deploy the application binary into a lean image
FROM debian:bullseye-20230320 AS build-release-stage

# Set the working directory
WORKDIR /

# Update the package lists and install required packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl

# Update the trusted CA certificates
RUN update-ca-certificates

# Copy application binary from build-stage
COPY --from=build-stage /utility-bot /utility-bot

# Set the entrypoint for the application
ENTRYPOINT ["/utility-bot", "discord:run"]
