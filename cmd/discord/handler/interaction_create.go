package handler

import (
	"github.com/bwmarrin/discordgo"
	commands "github.com/erosdesire/discord-nestjs-utility-bot/cmd/discord/command"
	"github.com/erosdesire/discord-nestjs-utility-bot/config"
	"github.com/rs/zerolog/log"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Debug().
		Str("name", i.ApplicationCommandData().Name).
		Str("guild-id", i.GuildID).
		Str("channel-id", i.ChannelID).
		Str("user-id", i.Member.User.ID).
		Msg("event: interaction create")

	done := contentInteractionHandler(s, i)
	if done {
		return
	}

	switch i.ApplicationCommandData().Name {
	case commands.NpmSearch:
		npmSearchCommandHandler(s, i)
		return
	}

	defaultInteractionHandler(s, i)
}

func contentInteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	name := i.ApplicationCommandData().Name

	// Check cache if the tag exists
	cmd, cmdExist := config.GetConfig().Commands[name]
	if !cmdExist {
		return false
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: cmd.Content,
		},
	})
	if err == nil {
		return true
	}

	log.Error().
		Err(err).
		Str("name", name).
		Str("guild-id", i.GuildID).
		Str("channel-id", i.ChannelID).
		Str("user-id", i.Member.User.ID).
		Msg("failed to respond to interaction")

	_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "Something went wrong",
	})

	return false
}

func npmSearchCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := ""
	sortBy := ""
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case commands.NpmSearchNameOption:
			name = option.StringValue()
		case commands.NpmSearchSortOption:
			sortBy = option.StringValue()
		}
	}

	log.Debug().
		Str("name", name).
		Str("sort", sortBy).
		Send()

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Got npm: " + name,
		},
	})
}

func defaultInteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Content not found.",
		},
	})
	if err != nil {
		_, _ = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
	}

	// Delete the slash command when it doesn't have any registered handler
	// _ = s.ApplicationCommandDelete(s.State.User.ID, config.GetConfig().GuildID, i.ID)
}
