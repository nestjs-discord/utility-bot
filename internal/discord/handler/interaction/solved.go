package interaction

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/solved"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/rs/zerolog/log"
	"strings"
)

func validateInteractionForThreadPost(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.Channel, bool) {
	currentChannelInfo, err := s.Channel(i.ChannelID)
	if err != nil {
		util.InteractionRespondError(
			fmt.Errorf("failed to get current channel info: %s", err),
			s, i)

		return nil, false
	}

	if currentChannelInfo.Type != discordgo.ChannelTypeGuildPublicThread &&
		currentChannelInfo.Type != discordgo.ChannelTypeGuildPrivateThread {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: ":warning: You can only use this command in forum posts.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		return nil, false
	}

	if currentChannelInfo.ThreadMetadata.Locked {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: ":warning: Cannot perform this action on a locked thread post",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return nil, false
	}

	return currentChannelInfo, true
}

func SolvedHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	currentChannelInfo, isValid := validateInteractionForThreadPost(s, i)
	if !isValid {
		return
	}

	// Restrict further actions to the original post owner and moderators
	userId := i.Member.User.ID
	postOwnerId := currentChannelInfo.OwnerID
	if userId != postOwnerId && !util.IsUserModerator(userId) {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: ":warning: Only forum post owner and moderators can use this command.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	parentChannelInfo, err := s.Channel(currentChannelInfo.ParentID)
	if err != nil {
		util.InteractionRespondError(
			fmt.Errorf("failed to get parent channel info: %s", err),
			s, i)
		return
	}

	solvedTag, err := findSolvedTag(parentChannelInfo.AvailableTags)
	if err != nil {
		util.InteractionRespondError(err, s, i)
		return
	}

	hasSolvedTag := false

	for _, appliedTag := range currentChannelInfo.AppliedTags {
		if appliedTag == solvedTag.ID {
			hasSolvedTag = true
			break
		}
	}
	if !hasSolvedTag {
		currentChannelInfo.AppliedTags = append(currentChannelInfo.AppliedTags, solvedTag.ID)
	}

	// https://discord.com/developers/docs/resources/channel#modify-channel-json-params-thread
	if len(currentChannelInfo.AppliedTags) > 5 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: ":warning: The current post already has five tags applied to it. " +
					"To apply the \"Solved\" tag, please remove at least one tag, " +
					"as Discord allows a maximum of 5 tags per forum post.",
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	//
	// Assign solved tag
	//
	// Discord doesn't allow responding to an interaction when the thread post is archived or closed.
	// Hence, editing the channel twice is necessary: first to apply tags, and second to close the thread post.
	//
	_, err = s.ChannelEdit(currentChannelInfo.ID, &discordgo.ChannelEdit{
		AppliedTags: &currentChannelInfo.AppliedTags,
	})
	if err != nil {
		util.InteractionRespondError(
			fmt.Errorf("failed to edit the channel to apply the solved tag: %s", err),
			s, i)
		return
	}

	// Send the canned response
	content := "This post has been marked as resolved. :white_check_mark:\n" +
                   "Please read through the conversation and resolution, if you are having the same issue. " +
                   "If you were the original author of the post and the issue is still fresh (within a few days) " +
	           "and you are still have having trouble, continue to reply here. If you are not the original " + 
	           "author of the post or the post has aged, start a new thread linking this one as relevant to " +
                   "your problem, providing as much additional information as possible."

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
	if err != nil {
		util.InteractionRespondError(fmt.Errorf("failed to respond to interaction: %s", err), s, i)
		return
	}

	// Default values when "auto-close" option isn't specified
	archived := false              // aka close
	autoArchiveDuration := 60 * 24 // a day | unit is minutes

	for _, option := range i.ApplicationCommandData().Options { // Check whether the "auto-close" option is specified
		if option.Name != solved.AutoClose {
			continue
		}

		optionValue, err := convertToInteger(option.Value)
		if err != nil {
			log.Err(err).Interface("value", option.Value).Msg("float64 to int conversion failed on auto-close option's value")
			return
		}

		if optionValue == 1 { // close right after
			archived = true
			continue
		}

		autoArchiveDuration = optionValue // but set its auto archive duration value to auto close in the future
	}

	_, err = s.ChannelEdit(currentChannelInfo.ID, &discordgo.ChannelEdit{
		Archived:            &archived,
		AutoArchiveDuration: autoArchiveDuration,
	})
	if err != nil {
		log.Err(err).Str("channel-id", currentChannelInfo.ID).Msg("solved command failed to edit the channel")
	}

	log.Debug().
		Int("auto-archive-dur", autoArchiveDuration).
		Bool("archived", archived).
		Msg("solved command executed")
}

func findSolvedTag(tags []discordgo.ForumTag) (*discordgo.ForumTag, error) {
	for _, tag := range tags {
		if strings.ToLower(tag.Name) == "solved" {
			return &tag, nil
		}
	}

	return nil, errors.New("failed to find the solved tag")
}

func convertToInteger(value interface{}) (int, error) {
	if floatValue, ok := value.(float64); ok {
		optionValue := int(floatValue)
		return optionValue, nil
	}

	return 0, errors.New("value is not a float64")
}
