# Discord NestJS Utility Bot

Placeholder.

## Motivation

Placeholder.

## Deployment

Placeholder.

## Dependencies

- [DiscordGo](https://github.com/bwmarrin/discordgo) - Provides low level bindings to the Discord chat client API
- [Cobra](https://github.com/spf13/cobra) - Commander for modern Go CLI interactions
- [Viper](https://github.com/spf13/viper) - Complete configuration solution for Go applications
- [Validator](https://github.com/go-playground/validator) - Implements value validations for structs based on tags
- [Zerolog](https://github.com/rs/zerolog) - Zero allocation JSON logger
- [Go-humanize](https://github.com/dustin/go-humanize) - Formatters for units to human friendly sizes
- [Testify](https://github.com/stretchr/testify) - A toolkit with common assertions and mocks

## Notes

It does not matter in what order commands are defined in `config.yml`; Discord will sort them alphabetically once they're registered.

## Todo

- [ ] refactor duplicate "Something went wrong" interaction responds
- [ ] replace any follow-up message with interaction respond

- [ ] npm related commands
  - [ ] validate npm packages name
  - [ ] validate version numbers
  - [ ] investigate subcommands
  - [ ] npm-inspect slash command https://registry.npmjs.org/@nestjs/core/latest

- [ ] rate limit usage of some slash commands
- [ ] move repository and go mod to the `nestjs-discord` organization
- [ ] research about sharding
- [ ] deployment
  - [ ] Docker files
  - [ ] CI build
  - [ ] CI tests
  - [ ] go releaser

- [ ] `readme.md`
  - [ ] improve the `readme.md` file
  - [ ] badges
