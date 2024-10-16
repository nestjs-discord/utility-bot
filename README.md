# Utility Bot

[![build-and-test](https://github.com/nestjs-discord/utility-bot/actions/workflows/build-and-test.yaml/badge.svg)](https://github.com/nestjs-discord/utility-bot/actions/workflows/build-and-test.yaml)
[![CodeFactor](https://www.codefactor.io/repository/github/nestjs-discord/utility-bot/badge/main)](https://www.codefactor.io/repository/github/nestjs-discord/utility-bot/overview/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/nestjs-discord/utility-bot)](https://goreportcard.com/report/github.com/nestjs-discord/utility-bot)

Carefully designed to streamline the support process for [the official NestJS Discord server](https://discord.gg/nestjs).

## Configuration

```sh
cp .env.example .env
```

## Running

```sh
# validates the YAML config (to be used by the contributors or in a CI/CD pipeline)
go run cmd/validate/main.go

# launches the discord bot
go run cmd/run/main.go

# cleans the registered application commands on the server
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
