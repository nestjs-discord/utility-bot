# Contributing to Utility Bot

Hello there! Thank you for considering contributing to the Utility Bot project, part of the NestJS-Discord organization.
We're thrilled to have you here!

## How to Contribute

1. Fork the repository
2. Clone the repository to your local machine
3. Run `go mod tidy` to install the dependencies
4. Make your changes
5. Run tests and ensure they pass with `go test -v ./...`
6. Ensure Markdown content in the configuration have proper length with `go run cmd/validate/main.go`
7. Push your changes to your fork
8. Create a pull request

We welcome all contributions, including bug reports, feature requests, and code improvements.
So let's build something great together!

Happy coding! 🚀

## Running the bot in development mode

```shell
go run cmd/run/main.go --stage dev
```

## Build

Tou must [install Golang](https://go.dev/doc/install) in your system and execute the following command:

```shell
go build -o bin/ cmd/run/main.go
```

## Notes

- Slash commands
    - Discord will sort commands alphabetically, regardless of their initial order in the `config.yml` file.
    - Users will instantly see commands once registered because this project uses [guild commands](https://discord.com/developers/docs/interactions/application-commands#registering-a-command) instead of global ones.
    - Discord has a global rate limit of [200 application command creations per day, per guild](https://discord.com/developers/docs/interactions/application-commands#registering-a-command).
    - The bot will automatically register slash commands on startup.
    - Registered commands can be removed by running `go run cmd/clean/main.go`.
    - Only one sub-command level is supported; for example, `/foo bar` is valid.
- Markdown content
    - The content of each file can be up to 2000 characters.
    - The bot caches the content after execution (restarting is required to apply the changes).
- Moderators
    - They can be defined by their unique Discord ID in the `config.yml` file.
    - They bypass rate-limit policies.
    - They can execute `protected` commands in the `config.yml` file.
