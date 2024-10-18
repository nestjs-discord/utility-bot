FROM golang:1.23.2-alpine AS build
WORKDIR /usr/app
RUN apk update && apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./
RUN CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="-s -w -buildid=" -o ./bin/bot cmd/run/main.go

FROM scratch AS prod
WORKDIR /usr/app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/app/bin/bot /usr/bin/bot
ENTRYPOINT ["bot"]
