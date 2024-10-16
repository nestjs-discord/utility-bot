package reference

import (
	"github.com/bwmarrin/discordgo"
	algolia2 "github.com/nestjs-discord/utility-bot/infra/services/algolia"
)

func AutocompleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rootOptions := i.ApplicationCommandData().Options
	var choices []*discordgo.ApplicationCommandOptionChoice

	for _, rootOption := range rootOptions {
		app, ok := algolia2.Apps[rootOption.Name]
		if !ok {
			continue
		}

		query, err := getStringValueByName(QueryOption, rootOption.Options)
		if err != nil {
			break // reply empty choices at the bottom
		}

		hits, err := algolia2.Search(app, query)
		if err != nil {
			break // reply empty choices at the bottom
		}

		// ().Interface("hits", hits).Msg("algolia search results")

		for _, hit := range hits {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  algolia2.Truncate(algolia2.GetFormattedHierarchy(hit), 95),
				Value: hit.ObjectID,
			})
		}

		// Break the loop after processing the first valid rootOption
		break
	}

	response := &discordgo.InteractionResponseData{
		Choices: choices,
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: response,
	})
}
