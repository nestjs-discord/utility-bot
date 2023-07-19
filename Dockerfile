# Build stage for Golang application
FROM golang:1.20.6-alpine AS build-stage

RUN apk update \
    && apk add --no-cache ca-certificates \
    && rm -rf /var/cache/apk/*

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
RUN go build -trimpath -buildvcs=false -ldflags "-w" -o /usr/bin/ub ./main.go

# Deploy the application binary into a lean image
FROM scratch AS production-stage

# Set the working directory
WORKDIR /usr/app

# Copy CA certificates from the build-stage image
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy application binary from build-stage image
COPY --from=build-stage /usr/bin/ub /usr/bin/ub

# Set the entrypoint for the application
ENTRYPOINT ["ub", "discord:run"]
