package dont_ping_mods

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/commands/common"
	"github.com/nestjs-discord/utility-bot/bot/moderators"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/rs/zerolog/log"
	"time"
)

type DontPingMods struct {
	moderators *moderators.Moderators
}

func NewDontPingMods(
	moderators *moderators.Moderators,
) *DontPingMods {
	return &DontPingMods{
		moderators: moderators,
	}
}

func (d *DontPingMods) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	content := "Please **do not** tag the moderators unless someone is breaking server rules. " +
		"The mods are here to help enforce the rules of the server, " +
		"and while most of them are knowledgeable about Nest, " +
		"they are not the only ones able to solve your question."

	for _, opt := range i.ApplicationCommandData().Options {
		// Mention the "target" user
		if opt.Name == common.OptionTarget && opt.Value != "" {
			content = fmt.Sprintf("*Suggestion for* <@%v>:\n", opt.Value) + content
		}
	}

	button := discordgo.Button{
		Emoji: &discordgo.ComponentEmoji{Name: "🔗"},
		Label: "Server Rules",
		Style: discordgo.LinkButton,
		URL:   "https://discord.com/channels/520622812742811698/527853342152458287/769643761797431336",
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						button,
					},
				},
			},
		},
	})
	if err != nil {
		msg := fmt.Errorf("failed to respond to interaction: %s", err)
		util.InteractionRespondError(msg, s, i)
		return
	}

	currentChannelInfo, err := s.Channel(i.ChannelID)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch the channel information on dont-ping-mods command")
		return
	}

	// Skip further steps when the current channel is not a forum post (thread)
	if currentChannelInfo.Type != discordgo.ChannelTypeGuildPublicThread &&
		currentChannelInfo.Type != discordgo.ChannelTypeGuildPrivateThread {
		return
	}

	// Loop over moderators defined in the configuration file
	for _, modId := range d.moderators.UserIds() {
		// Skip removing the person who have executed the command
		if modId == i.Member.User.ID {
			continue
		}

		// Remove the moderator from the forum post
		if err := s.ThreadMemberRemove(i.ChannelID, modId); err != nil {
			log.Error().Err(err).
				Str("mod-user-id", modId).
				Str("thread-id", i.ChannelID).
				Msg("failed to remove the mod from the thread")
		}

		// Sleep for a bit to avoid flooding Discord API
		time.Sleep(100 * time.Millisecond)
	}

	// Silence is golden :)
}
