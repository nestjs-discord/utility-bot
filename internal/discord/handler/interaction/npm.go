package interaction

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/erosdesire/discord-nestjs-utility-bot/internal/discord/command/npm"
	npmAPI "github.com/erosdesire/discord-nestjs-utility-bot/internal/npm"
	"github.com/rs/zerolog/log"
)

func NpmSearchHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := ""
	sortBy := ""
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case npm.SearchOptionName:
			name = option.StringValue()
		case npm.SearchOptionSort:
			// option.IntValue()
			// TODO:  logic
		}
	}

	log.Debug().
		Str("name", name).
		Str("sort", sortBy).
		Send()

	options := &npmAPI.SearchOptions{
		Text: name,
	}

	data, err := npmAPI.Search(options)
	if err != nil {
		_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// TODO: generate proper message
			Content: "Got npm: " + name + " " + fmt.Sprint(data.Total),
		},
	})
}
