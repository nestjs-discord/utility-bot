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

- [ ] investigate subcommands for npm related commands
- [ ] Cobra `discord:clean` command to remove registered application commands
  - [ ] Remove the registered commands array 
- [ ] validate npm packages name
- [ ] validate version numbers
- [ ] replace any follow-up message with interaction respond
- [ ] npm-inspect slash command https://registry.npmjs.org/@nestjs/core/latest
- [ ] rate limit usage of some slash commands
- [ ] improve the `readme.md` file
- [ ] refactor duplicate "Something went wrong" interaction responds
