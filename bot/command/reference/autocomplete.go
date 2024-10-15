package reference

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/pkg/algolia"
)

func AutocompleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rootOptions := i.ApplicationCommandData().Options
	var choices []*discordgo.ApplicationCommandOptionChoice

	for _, rootOption := range rootOptions {
		app, ok := algolia.Apps[rootOption.Name]
		if !ok {
			continue
		}

		query, err := getStringValueByName(QueryOption, rootOption.Options)
		if err != nil {
			break // reply empty choices at the bottom
		}

		hits, err := algolia.Search(app, query)
		if err != nil {
			break // reply empty choices at the bottom
		}

		// ().Interface("hits", hits).Msg("algolia search results")

		for _, hit := range hits {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  algolia.Truncate(algolia.GetFormattedHierarchy(hit), 95),
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
