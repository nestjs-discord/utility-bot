package google_it

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"github.com/nestjs-discord/utility-bot/infra/services/google"
	"log/slog"
)

const (
	Name         = "google-it"
	OptionSuffix = "-suggestion"
)

type GoogleIt struct {
	logger *slog.Logger
	client *google.Google
}

func NewGoogleIt() *GoogleIt {
	return &GoogleIt{
		logger: logger.NewWithSubsystem("bot", "commands", "googleIt"),
		client: google.NewGoogle(),
	}
}

func (g *GoogleIt) Command() *dgo.ApplicationCommand {
	cmd := &dgo.ApplicationCommand{
		Name:        Name,
		Description: "Tell someone to search on Google or StackOverflow!",
		Type:        dgo.ChatApplicationCommand,
	}

	elements := []string{"first", "second", "third", "fourth"}
	minLength := 3

	for i, opt := range elements {
		cmd.Options = append(cmd.Options, &dgo.ApplicationCommandOption{
			Name:         opt + OptionSuffix,
			Description:  opt + " suggestion",
			Type:         dgo.ApplicationCommandOptionString,
			Required:     i == 0, // only the first item is required
			Autocomplete: true,
			MinLength:    &minLength,
		})
	}

	return cmd
}
