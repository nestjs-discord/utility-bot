# Utility Bot

[![build-and-test](https://github.com/nestjs-discord/utility-bot/actions/workflows/build-and-test.yaml/badge.svg)](https://github.com/nestjs-discord/utility-bot/actions/workflows/build-and-test.yaml)
[![CodeFactor](https://www.codefactor.io/repository/github/nestjs-discord/utility-bot/badge/main)](https://www.codefactor.io/repository/github/nestjs-discord/utility-bot/overview/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/nestjs-discord/utility-bot)](https://goreportcard.com/report/github.com/nestjs-discord/utility-bot)

A Discord bot designed to streamline the support process on the official NestJS Discord server.

As people who usually answer questions on [the official NestJS Discord server](https://discord.gg/nestjs), we've experienced that sometimes users ask questions that have already been answered many times before.

There are some common cases, like when they post a new support request, they don't provide a [minimal reproduction code](https://minimum-reproduction.wtf/), or sometimes they don't even share their code, and even if they do, they don't know how to put them in [code blocks](https://gist.github.com/matthewzring/9f7bbfd102003963f9be7dbcf7d40e51#code-blocks) properly.

So we devised the idea of having a [Discord bot](https://discord.com/developers/docs/intro#bots-and-apps) with predefined and well-written Markdown content as [slash commands](https://discord.com/developers/docs/interactions/application-commands) to reply to users instead of repeatedly writing and explaining.

## Discord commands overview

| Command                        | Description                                                                    | Who can execute it?    |
|--------------------------------|--------------------------------------------------------------------------------|------------------------|
| `/stats`                       | General statistics about the bot instance (resource usage)                     | Mods                   |
| `/solved`                      | To close and mark forum posts as solved by a tag called "solved"               | Mods, Forum post owner |
| `/archive`                     | To close and lock the forum posts                                              | Mods                   |
| `/mods cross-post`             | Warn users not to cross-post, ask in relevant channels, and be patient         | Mods                   |
| `/mods spam`                   | A friendly warning to prevent spamming in the chat                             | Mods                   |
| `/mods elaborate`              | Please provide more context so we can better assist you with your inquiry      | Everyone               |
| `/community awesome`           | A curated list of awesome things related to NestJS (`awesome-nestjs`)          | Everyone               |
| `/community testing`           | Show off to the community methods of testing NestJS (`jmcdo29/testing-nestjs`) | Everyone               |
| `/docs behind-proxy`           | How to enable the trust proxy option when NestJS app is behind a reverse proxy | Everyone               |
| `/docs circular-dependency`    | Circular dependencies in NestJS and how to avoid them                          | Everyone               |
| `/docs file-change-detected`   | How to fix when TypeScript 4.9+ on Windows causes infinite loop in watch mode  | Everyone               |
| `/docs nest-debug`             | Debug NestJS dependencies with environment variable                            | Everyone               |
| `/docs providers-vs-imports`   | Ensure accurate service/module listings                                        | Everyone               |
| `/docs request-lifecycle`      | Get a grip on NestJS' request lifecycle for smooth coding                      | Everyone               |
| `/docs resolve-dependencies`   | Nest not being able to resolve dependencies of something                       | Everyone               |
| `/rules dm`                    | Use server to get help on NestJS; don't DM other members                       | Everyone               |
| `/rules screenshot`            | When someone posts a screenshot instead of sharing their code                  | Everyone               |
| `/javascript floating-promise` | Introduction to floating Promises in JavaScript                                | Everyone               |

## Configuration

```shell
cp .env.example .env
```

Three `DISCORD_APP_ID`, `DISCORD_BOT_TOKEN`, and `DISCORD_GUILD_ID` environment variables are required, and the rest of
the configuration is located in `config.yml`.

Currently, the bot doesn't support hot-reloading. Instead, the application should restart to apply the changes.

## Adding the bot to the server

The `discord:invite` CLI command generates the server invite link,
as the permissions within the link will likely change based on the features we may add.

## Docker usage

```shell
docker compose --env-file ./.env up -d --build

# Generate an invite link to add the bot to the server
docker compose exec utility-bot ./utility-bot discord:invite

# Graceful shutdown
# docker compose down

# Removing the registered slash commands (in case of inconsistency between the configuration and the production)
# docker compose exec utility-bot ./utility-bot discord:clean
```

## Notes

- Slash commands
  - Discord will sort them alphabetically, regardless of their initial order in `config.yml`.
  - They update instantly for the end-user because this project uses [guild commands](https://discord.com/developers/docs/interactions/application-commands#registering-a-command) instead of global commands.
  - Discord has a global rate limit of [200 application command creations per day, per guild](https://discord.com/developers/docs/interactions/application-commands#registering-a-command).
  - Bot will automatically register slash commands on startup.
  - Registered slash commands can be removed by the `discord:clean` command.
  - Only one sub-command level is supported; for example, `/foo bar` is valid.
- Markdown content
  - Content within the slash commands can have a maximum of 2000 characters.
  - The bot will cache Markdown content to avoid spamming I/O.
- Moderators
  - They can be defined by their unique Discord ID in the `config.yml` file.
  - They bypass rate-limit policies.
  - They can execute `protected` commands in `config.yml`.
