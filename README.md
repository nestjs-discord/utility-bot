<p align="center"><img src=".github/logo.png?raw=true" alt="Utility Bot" height="250px" /></p>

# Utility Bot

[![tag](https://img.shields.io/github/tag/nestjs-discord/utility-bot.svg)](https://github.com/nestjs-discord/utility-bot/releases)
[![build-and-test](https://github.com/nestjs-discord/utility-bot/actions/workflows/build-and-test.yaml/badge.svg)](https://github.com/nestjs-discord/utility-bot/actions/workflows/build-and-test.yaml)
[![golangci-lint](https://github.com/nestjs-discord/utility-bot/actions/workflows/golangci-lint.yml/badge.svg?branch=main)](https://github.com/nestjs-discord/utility-bot/actions/workflows/golangci-lint.yml)
[![CodeFactor](https://www.codefactor.io/repository/github/nestjs-discord/utility-bot/badge/main)](https://www.codefactor.io/repository/github/nestjs-discord/utility-bot/overview/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/nestjs-discord/utility-bot)](https://goreportcard.com/report/github.com/nestjs-discord/utility-bot)
[![Discord](https://img.shields.io/discord/520622812742811698?logo=nestjs&logoColor=%23e0234e&label=Discord&color=%235765F2)](<https://discord.gg/nestjs>)

Crafted with love for the official NestJS Discord server. :heart:

# Goals

The primary objectives of this project are:

1. Automate moderation actions.
2. Enable members to respond to common scenarios using [application commands](<https://discord.com/developers/docs/interactions/application-commands>).
3. Provide tools to encourage engagement and collaboration.

> [!NOTE]
> This bot is custom-made for our server, and all features are subject to change based on our community's needs.
> We reserve the right to introduce breaking changes at any time.

# Configuration

```sh
cp .env.example .env
```

# At a glance

```sh
# validates the configuration files (can be used by the contributors or in a CI/CD pipeline)
go run cmd/validate/main.go

# prints out the server invite link
go run cmd/invite/main.go

# launches the Discord bot
go run cmd/run/main.go

# cleans the registered application commands
go run cmd/clean/main.go

# generates the bot status texts (just for fun 🙂)
cd cmd/status/
go run .
```

# Donations

TBD.
