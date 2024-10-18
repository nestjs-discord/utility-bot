# Utility Bot

[![tag](https://img.shields.io/github/tag/nestjs-discord/utility-bot.svg)](https://github.com/nestjs-discord/utility-bot/releases)
[![build-and-test](https://github.com/nestjs-discord/utility-bot/actions/workflows/build-and-test.yaml/badge.svg)](https://github.com/nestjs-discord/utility-bot/actions/workflows/build-and-test.yaml)
[![CodeFactor](https://www.codefactor.io/repository/github/nestjs-discord/utility-bot/badge/main)](https://www.codefactor.io/repository/github/nestjs-discord/utility-bot/overview/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/nestjs-discord/utility-bot)](https://goreportcard.com/report/github.com/nestjs-discord/utility-bot)
[![Discord](https://img.shields.io/discord/520622812742811698?logo=nestjs&logoColor=%23e0234e&label=Discord&color=%235765F2)](<https://discord.gg/nestjs>)

Carefully designed with love ❤️ for the official NestJS Discord server.

## Configuration

```sh
cp .env.example .env
```

## Running

```sh
# validates the configuration files (can be used by the contributors or in a CI/CD pipeline)
go run cmd/validate/main.go

# launches the Discord bot
go run cmd/run/main.go

# cleans the registered application commands
go run cmd/clean/main.go
```

## Docker

```sh
docker compose down --remove-orphans
docker compose --env-file ./.env up -d --build
docker compose ps
docker stats --no-stream
docker compose logs --follow
```
