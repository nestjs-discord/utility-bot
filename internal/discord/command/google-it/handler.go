package google_it

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
)

func Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var suggestions []string

	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Type != discordgo.ApplicationCommandOptionString {
			continue
		}

		if !strings.HasSuffix(opt.Name, OptionSuffix) {
			continue
		}

		suggestions = append(suggestions, opt.Value.(string))
	}

	log.Debug().Strs("suggestions", suggestions).Msg("google-it handler")

	content := "Please make an effort to use services like Google or Stack Overflow to search for your question before submitting it here."
	content += " "
	content += "There is a very decent chance your problem has already been solved by someone else in some way."

	if len(suggestions) > 0 {
		content += "\n\n"
		content += "Here are some search suggestions that could yield some results:"
		content += "\n"

		for _, suggestion := range suggestions {
			content += "- <https://www.google.com/search?q=" + url.QueryEscape(suggestion) + ">" + "\n"
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		log.Err(err).Strs("suggestions", suggestions).Msg("interaction response failed in google-it handler")
	}
}
