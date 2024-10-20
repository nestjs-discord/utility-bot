package reference

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/infra/services/algolia"
)

func (r *Reference) AutocompleteHandler(s *dgo.Session, i *dgo.InteractionCreate) {
	rootOptions := i.ApplicationCommandData().Options
	var choices []*dgo.ApplicationCommandOptionChoice

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

		for _, hit := range hits {
			choices = append(choices, &dgo.ApplicationCommandOptionChoice{
				Name:  algolia.Truncate(algolia.GetFormattedHierarchy(hit), 95),
				Value: hit.ObjectID,
			})
		}

		// Break the loop after processing the first valid rootOption
		break
	}

	response := &dgo.InteractionResponseData{
		Choices: choices,
	}

	_ = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionApplicationCommandAutocompleteResult,
		Data: response,
	})
}
