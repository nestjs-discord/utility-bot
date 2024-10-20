package google_it

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/security"
	"log/slog"
)

func (g *GoogleIt) AutocompleteHandler(s *dgo.Session, i *dgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	focusedValue := extractFocusedValue(options)
	var choices []*dgo.ApplicationCommandOptionChoice

	// to avoid spamming google api
	if len(focusedValue) < 3 {

		choices = append(choices, &dgo.ApplicationCommandOptionChoice{
			Name:  focusedValue,
			Value: focusedValue,
		})

		_ = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
			Type: dgo.InteractionApplicationCommandAutocompleteResult,
			Data: &dgo.InteractionResponseData{
				Choices: choices,
			},
		})
		return
	}

	g.logger.Debug("autocomplete",
		slog.String("focusedValue", focusedValue),
	)

	focusedValue = security.RemoveDangerousMentions(focusedValue)

	res, err := g.client.Search(focusedValue)
	if err != nil {
		g.logger.Error("google client search failed",
			slog.String("query", focusedValue),
		)
		return
	}

	for _, item := range res {
		choices = append(choices, &dgo.ApplicationCommandOptionChoice{
			Name:  item,
			Value: item,
		})
	}

	err = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionApplicationCommandAutocompleteResult,
		Data: &dgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		g.logger.Error("auto complete interaction respond failed",
			slog.Any("err", err),
			slog.String("focusedValue", focusedValue),
		)
	}
}

func extractFocusedValue(options []*dgo.ApplicationCommandInteractionDataOption) string {
	for _, opt := range options {
		if opt.Type != dgo.ApplicationCommandOptionString {
			continue
		}

		if opt.Focused {
			return opt.Value.(string)
		}
	}

	return ""
}
