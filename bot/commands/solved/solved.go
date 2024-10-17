package solved

import (
	"errors"
	"fmt"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
	"github.com/nestjs-discord/utility-bot/bot/moderators"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"github.com/rs/zerolog/log"
	"log/slog"
	"strings"

	dgo "github.com/bwmarrin/discordgo"
)

const (
	Name      = "solved"
	AutoClose = "auto-close"
)

type Solved struct {
	logger     *slog.Logger
	moderators *moderators.Moderators
	// TODO: yaml config
}

func New(moderators *moderators.Moderators) *Solved {
	return &Solved{
		logger:     logger.NewWithSubsystem("bot", "commands", "solved"),
		moderators: moderators,
	}
}

func (c *Solved) Handler(s *dgo.Session, i *dgo.InteractionCreate) error {
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		c.logger.Error("unable to get channel info",
			slog.String("channelId", i.ChannelID),
		)
		return fmt.Errorf("unable to get channel info: %w", err)
	}

	if !c.validateChannelType(s, i, channel) ||
		!c.validateThreadLock(s, i, channel) ||
		!c.validateChannelOwner(s, i, channel) {
		return nil
	}

	// TODO: can this be loaded form the config file?
	parentChannelInfo, err := s.Channel(channel.ParentID)
	if err != nil {
		return errors.New("TODO")
	}

	solvedTag, err := c.findSolvedTag(parentChannelInfo.AvailableTags)
	if err != nil {
		return errors.New("TODO")
	}

	hasSolvedTag := false

	for _, appliedTag := range channel.AppliedTags {
		if appliedTag == solvedTag.ID {
			hasSolvedTag = true
			break
		}
	}
	if !hasSolvedTag {
		channel.AppliedTags = append(channel.AppliedTags, solvedTag.ID)
	}

	// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-thread
	if len(channel.AppliedTags) > 5 {
		msg := ":warning: The current post already has five tags applied to it. " +
			"To apply the \"Solved\" tag, please remove at least one tag, " +
			"as Discord allows a maximum of 5 tags per forum post."
		respond.InteractionWithEphemeralMessage(s, i, msg)
		return nil
	}

	//
	// Assign solved tag
	//
	// Discord doesn't allow responding to an interaction when the thread post is archived or closed.
	// Hence, editing the channel twice is necessary: first to apply tags, and second to close the thread post.
	//
	_, err = s.ChannelEdit(channel.ID, &dgo.ChannelEdit{
		AppliedTags: &channel.AppliedTags,
	})
	if err != nil {
		return fmt.Errorf("failed to edit the channel to apply the solved tag: %w", err)
	}

	// Send the canned response
	content := "This post has been marked as resolved. :white_check_mark:\n" +
		"Please read through the conversation and resolution, if you are having the same issue. " +
		"If you were the original author of the post and the issue is still fresh (within a few days) " +
		"and you are still have having trouble, continue to reply here. If you are not the original " +
		"author of the post or the post has aged, start a new thread linking this one as relevant to " +
		"your problem, providing as much additional information as possible."

	err = s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{Content: content},
	})
	if err != nil {
		return fmt.Errorf("failed to respond to interaction: %w", err)
	}

	// Default values when "auto-close" option isn't specified
	archived := false              // aka close
	autoArchiveDuration := 60 * 24 // a day | unit is minutes

	for _, option := range i.ApplicationCommandData().Options { // Check whether the "auto-close" option is specified
		if option.Name != AutoClose {
			continue
		}

		optionValue, err := c.convertToInteger(option.Value)
		if err != nil {
			return fmt.Errorf("float64 to int conversion failed on auto-close option value: %s", err)
		}

		if optionValue == 1 { // close right after
			archived = true
			continue
		}

		autoArchiveDuration = optionValue // but set its auto archive duration value to auto close in the future
	}

	_, err = s.ChannelEdit(channel.ID, &dgo.ChannelEdit{
		Archived:            &archived,
		AutoArchiveDuration: autoArchiveDuration,
	})
	if err != nil {
		log.Err(err).Str("channel-id", channel.ID).Msg("solved command failed to edit the channel")
	}

	log.Debug().
		Int("auto-archive-dur", autoArchiveDuration).
		Bool("archived", archived).
		Msg("solved command executed")
	return nil
}

func (c *Solved) findSolvedTag(tags []dgo.ForumTag) (*dgo.ForumTag, error) {
	for _, tag := range tags {
		if strings.ToLower(tag.Name) == "solved" {
			return &tag, nil
		}
	}

	return nil, errors.New("failed to find the solved tag")
}

func (c *Solved) convertToInteger(value interface{}) (int, error) {
	if floatValue, ok := value.(float64); ok {
		optionValue := int(floatValue)
		return optionValue, nil
	}

	return 0, errors.New("value is not a float64")
}
