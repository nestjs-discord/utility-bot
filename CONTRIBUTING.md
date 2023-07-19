# Contributing to Utility Bot

Hello there! Thank you for considering contributing to the Utility Bot project, part of the NestJS-Discord organization.
We're thrilled to have you here!

## How to Contribute

1. Fork the repository
2. Clone the repository to your local machine
3. Make your changes
4. Run tests and ensure they pass with `go test -v ./...`
5. Ensure Markdown content in the configuration have proper length with `go run . validate:content`
6. Push your changes to your fork
7. Create a pull request

We welcome all contributions, including bug reports, feature requests, and code improvements.
So let's build something great together!

Happy coding! 🚀

## Running the Discord bot

```shell
go run . discord:run --debug
```

## Build

To build this project, you must [install Golang](https://go.dev/doc/install) in your system
and execute the following command.

```shell
go build -trimpath -buildvcs=false -ldflags "-w" -o ./bin/utility-bot ./main.go
```

## Dependencies overview

- [DiscordGo](https://github.com/bwmarrin/discordgo) - Provides low-level bindings to the Discord chat client API
- [Cobra](https://github.com/spf13/cobra) - Commander for modern Go CLI interactions
- [Viper](https://github.com/spf13/viper) - Complete configuration solution for Go applications
- [Validator](https://github.com/go-playground/validator) - Implements value validations for structs based on tags
- [Zerolog](https://github.com/rs/zerolog) - Zero allocation JSON logger
- [Go-humanize](https://github.com/dustin/go-humanize) - Formatters for units to human friendly sizes
- [Testify](https://github.com/stretchr/testify) - A toolkit with common assertions and mocks
- [Prometheus Golang](https://github.com/prometheus/client_golang) - Prometheus instrumentation library for Go apps
