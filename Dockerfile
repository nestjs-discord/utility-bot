# Build stage for Golang application
FROM golang:1.20.6-alpine AS build-stage

RUN apk update \
    && apk add --no-cache ca-certificates make \
    && rm -rf /var/cache/apk/*

# Set the working directory
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum Makefile ./

RUN make install

# Copy application source code
COPY . .

# Build application using the Makefile
RUN make build

# Deploy the application binary into a lean image
FROM scratch AS production-stage

# Set the working directory
WORKDIR /usr/app

# Copy CA certificates from the build-stage image
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy application binary from build-stage image
COPY --from=build-stage /app/bin/ub /usr/bin/ub

# Set the entrypoint for the application
ENTRYPOINT ["ub", "discord:run"]
