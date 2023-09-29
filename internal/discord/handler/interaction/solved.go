package interaction

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
)

func validateInteractionForThreadPost(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.Channel, bool) {
	currentChannelInfo, err := s.Channel(i.ChannelID)
	if err != nil {
		util.InteractionRespondError(
			errors.Wrap(err, "failed to get current channel info"),
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
			errors.Wrap(err, "failed to get parent channel info"),
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

			//_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			//	Type: discordgo.InteractionResponseChannelMessageWithSource,
			//	Data: &discordgo.InteractionResponseData{
			//		Content: ":warning: The solved tag is already applied to this channel.",
			//		Flags:   discordgo.MessageFlagsEphemeral,
			//	},
			//})
			//return
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

	// Assign solved tag
	//
	// Discord doesn't allow us to send interaction responses when the thread post is archived or locked
	// that's why I decided to edit the channel twice, the first time for applying tags
	// and the second time for archiving or closing (locking) the thread post.

	_, err = s.ChannelEdit(currentChannelInfo.ID, &discordgo.ChannelEdit{
		AppliedTags: &currentChannelInfo.AppliedTags,
	})
	if err != nil {
		util.InteractionRespondError(
			errors.Wrap(err, "failed to edit the channel to apply the solved tag"),
			s, i)
		return
	}

	// Send the canned response
	content := "This post has been marked as resolved. :white_check_mark:\n" +
		"Please read through the conversation and resolution if you are having the same issue, " +
		"and then re-open the post if you are still having trouble, " +
		"providing as much extra information as possible."

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: content},
	})
	if err != nil {
		util.InteractionRespondError(errors.Wrap(err, "failed to respond to interaction"), s, i)
		return
	}

	// Mark thread post as closed

	archived := true
	_, err = s.ChannelEdit(currentChannelInfo.ID, &discordgo.ChannelEdit{
		Archived:            &archived,
		AutoArchiveDuration: 60, // 60 minutes
	})
	if err != nil {
		log.Err(err).Str("channel-id", currentChannelInfo.ID).Msg("solved command failed to edit the channel")
	}
}

func findSolvedTag(tags []discordgo.ForumTag) (*discordgo.ForumTag, error) {
	for _, tag := range tags {
		if strings.ToLower(tag.Name) == "solved" {
			return &tag, nil
		}
	}

	return nil, errors.New("failed to find the solved tag")
}
