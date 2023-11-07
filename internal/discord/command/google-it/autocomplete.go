package google_it

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func AutocompleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	focusedValue := extractFocusedValue(options)
	var choices []*discordgo.ApplicationCommandOptionChoice

	// to avoid spamming search api
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

	log.Debug().Str("focused-value", focusedValue).Msg("google-it: autocomplete")

	res, err := searchClient.Search(focusedValue)
	if err != nil {
		log.Err(err).Str("query", focusedValue).Msg("search client failed to query a value")
		return
	}

	for _, item := range res {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  item,
			Value: item,
		})
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	}); err != nil {
		log.Err(err).Str("focused-value", focusedValue).Msg("interaction respond failed on google-it autocomplete handler")
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
