# Utility Bot

Placeholder.

## Motivation

Placeholder.

## Deployment

Placeholder.

## Tests

Placeholder.

## Notes

- Configuration
  - Three `DISCORD_APP_ID`, `DISCORD_BOT_TOKEN`, and `DISCORD_GUILD_ID` environment variables are required.
  - Rest of the configuration will be loaded from the `config.yml` file.
  - Currently, the bot doesn't support hot-reloading. Instead, the application should restart to apply the changes.
- Slash commands
  - Once slash commands are registered, Discord will sort them alphabetically, regardless of their initial order in `config.yml`.
  - Once slash commands are registered or removed, they get updated instantly for the end-user because this project uses [guild commands](https://discord.com/developers/docs/interactions/application-commands#registering-a-command) instead of global commands.
  - Discord has a global rate limit of [200 application command creations per day, per guild](https://discord.com/developers/docs/interactions/application-commands#registering-a-command).
  - Bot will automatically register non-registered slash commands on bootstrap.
  - Registered slash commands can be removed by `discord:clean` command.
- Markdown content
  - Content within the slash commands can have a maximum of 3500 characters.
  - The bot will cache Markdown content on memory to avoid spamming I/O.
- Moderators
  - They can be defined by their unique Discord ID in `config.yml`.
  - They bypass rate-limit policies.
  - They can execute `protected` commands in `config.yml`.

## Dependencies overview

- [DiscordGo](https://github.com/bwmarrin/discordgo) - Provides low level bindings to the Discord chat client API
- [Cobra](https://github.com/spf13/cobra) - Commander for modern Go CLI interactions
- [Viper](https://github.com/spf13/viper) - Complete configuration solution for Go applications
- [Validator](https://github.com/go-playground/validator) - Implements value validations for structs based on tags
- [Zerolog](https://github.com/rs/zerolog) - Zero allocation JSON logger
- [Go-humanize](https://github.com/dustin/go-humanize) - Formatters for units to human friendly sizes
- [Testify](https://github.com/stretchr/testify) - A toolkit with common assertions and mocks

## Todo

- [ ] npm related
  - Before interacting with the npm registry API
    - [ ] validate npm package names
    - [ ] validate version numbers

- [ ] features
  - [ ] rate limit usage of some slash commands
  - [ ] mark slash commands as `protected`
  - [ ] "npm > inspect" slash command https://registry.npmjs.org/@nestjs/core/latest

- [ ] refactor
  - [ ] wrap errors

- [ ] deployment
  - [ ] Docker files
  - [ ] go releaser, maybe?

- [ ] `readme.md`
  - [ ] improve the `readme.md` file
  - [ ] badges
