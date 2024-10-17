package archive

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
	"github.com/nestjs-discord/utility-bot/bot/moderators"
	"github.com/nestjs-discord/utility-bot/bot/permissions"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log/slog"
)

const Name = "archive"

type Archive struct {
	logger     *slog.Logger
	moderators *moderators.Moderators
	// TODO: yaml config
}

func New(moderators *moderators.Moderators) *Archive {
	return &Archive{
		logger:     logger.NewWithSubsystem("bot", "commands", "archive"),
		moderators: moderators,
	}
}

func (a *Archive) Command() *dgo.ApplicationCommand {
	var perm = int64(permissions.ProtectedCommands)

	return &dgo.ApplicationCommand{
		Name:                     Name,
		Description:              "Close and lock a forum post.",
		DefaultMemberPermissions: &perm,
	}
}

func (a *Archive) Handler(s *dgo.Session, i *dgo.InteractionCreate) error {
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		return fmt.Errorf("unable to get channel info: %s", err)
	}

	if !a.moderators.IsUserModerator(i.Member.User.ID) {
		respond.InteractionWithEphemeralMessage(s, i, "⚠️ Only the moderators can perform this task.")
		return nil
	}

	if !a.validateChannelType(s, i, channel) ||
		!a.validateThreadLock(s, i, channel) {
		return nil // already sent the response
	}

	content := "This post has been marked as \"archived\".\n" +
		"Please use it as a reference, but do not re-open it. " +
		"If you have a similar issue and cannot resolve it after reading this thread, please open a new post."

	err = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{Content: content},
	})
	if err != nil {
		return fmt.Errorf("unable to respond to interaction: %s", err)
	}

	archived := true
	locked := true

	_, err = s.ChannelEdit(i.ChannelID, &dgo.ChannelEdit{
		Archived: &archived,
		Locked:   &locked,
	})
	if err != nil {
		return fmt.Errorf("failed to update channel: %s", err)
	}

	return nil
}
