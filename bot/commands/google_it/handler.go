package google_it

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/security"
	"log/slog"
	"net/url"
	"strings"
)

func (g *GoogleIt) Handler(s *dgo.Session, i *dgo.InteractionCreate) error {
	var suggestions []string

	options := i.ApplicationCommandData().Options
	for _, opt := range options {
		if opt.Type != dgo.ApplicationCommandOptionString {
			continue
		}

		if !strings.HasSuffix(opt.Name, OptionSuffix) {
			continue
		}

		suggestions = append(suggestions, opt.Value.(string))
	}

	g.logger.Debug("triggered",
		slog.Any("suggestions", suggestions),
	)

	content := "Please make an effort to use services like Google or Stack Overflow " +
		"to search for your question before submitting it here." +
		" " +
		"There is a very decent chance your problem has already been solved by someone else in some way."

	if len(suggestions) > 0 {
		content += "\n\n"
		content += "Here are some search suggestions that could yield some results:"
		content += "\n"

		for _, suggestion := range suggestions {
			suggestion = security.RemoveDangerousMentions(suggestion)
			content += "- <https://google.com/search?q=" + url.QueryEscape(suggestion) + ">" + "\n"
		}
	}

	err := s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		return fmt.Errorf("interaction respond failed: %s", err)
	}

	return nil
}
