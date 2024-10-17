package google_it

import (
	"github.com/bwmarrin/discordgo"
	"log/slog"
)

func (g *GoogleIt) AutocompleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	focusedValue := extractFocusedValue(options)
	var choices []*discordgo.ApplicationCommandOptionChoice

	// to avoid spamming google api
	if len(focusedValue) < 3 {

		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  focusedValue,
			Value: focusedValue,
		})

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
		return
	}

	g.logger.Debug("autocomplete",
		slog.String("focusedValue", focusedValue),
	)

	res, err := g.client.Search(focusedValue)
	if err != nil {
		g.logger.Error("google client search failed",
			slog.String("query", focusedValue),
		)
		return
	}

	for _, item := range res {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  item,
			Value: item,
		})
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
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

func extractFocusedValue(options []*discordgo.ApplicationCommandInteractionDataOption) string {
	for _, opt := range options {
		if opt.Type != discordgo.ApplicationCommandOptionString {
			continue
		}

		if opt.Focused {
			return opt.Value.(string)
		}
	}

	return ""
}
