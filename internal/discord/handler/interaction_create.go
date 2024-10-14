package handler

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	forms2 "github.com/nestjs-discord/utility-bot/bot/forms"
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/cache"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/archive"
	dont_ping_mods "github.com/nestjs-discord/utility-bot/internal/discord/command/dont-ping-mods"
	google_it "github.com/nestjs-discord/utility-bot/internal/discord/command/google-it"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/reference"
	"github.com/nestjs-discord/utility-bot/internal/discord/command/solved"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler/interaction"
	"github.com/nestjs-discord/utility-bot/internal/discord/util"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"log/slog"
	"strings"
	"sync"
	"time"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		handleInteractionApplicationCommand(s, i)
		return
	case discordgo.InteractionMessageComponent: // interactive button (form)
		handleInteractionMessageComponent(s, i)
		return
	case discordgo.InteractionModalSubmit: // modal submit (form)
		handleInteractionModalSubmit(s, i)
		return
	case discordgo.InteractionApplicationCommandAutocomplete:
		handleInteractionApplicationCommandAutocomplete(s, i)
		return
	}
}

// interactionCommandHandlerMap maps command names against their handler
type interactionCommandHandlerMap map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)

func handleInteractionApplicationCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	userID := i.Member.User.ID

	slog.Debug("event: interaction app command",
		slog.String("userId", userID),
		slog.String("channelId", i.ChannelID),
		slog.String("name", data.Name),
		slog.Any("options", i.ApplicationCommandData().Options),
	)

	if checkRateLimit(userID) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: config.Yaml().RateLimit.Message,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			util.InteractionRespondError(err, s, i)
		}
		return
	}

	handlers := interactionCommandHandlerMap{
		solved.Name:         interaction.SolvedHandler,
		archive.Name:        interaction.ArchiveHandler,
		reference.Name:      reference.Handler,
		google_it.Name:      google_it.Handler,
		dont_ping_mods.Name: dont_ping_mods.Handler,
	}

	if handler, ok := handlers[data.Name]; ok {
		handler(s, i)
		return
	}

	if interaction.ContentHandler(s, i) {
		return
	}

	interaction.DefaultHandler(s, i)
}

// // TODO: refactor - to prevent race condition (when moderators click on the accept/reject/ban buttons)
var handleInteractionMessageComponentLock = sync.Mutex{}

func handleInteractionMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handleInteractionMessageComponentLock.Lock()
	defer handleInteractionMessageComponentLock.Unlock()

	// TODO: refactor this, it should only be called for the moderator actions, not the public interactive button!
	cacheKey := i.Message.ID
	cacheValue, cacheHit := forms2.ModActionsCache.Get(cacheKey)
	if cacheHit && cacheValue {
		msg := "Race condition detected! 😅\n"
		msg += "Another moderator has already handled this message."
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	cacheTtl := 30 * time.Second
	forms2.ModActionsCache.SetWithTTL(cacheKey, true, 1, cacheTtl)

	data := i.MessageComponentData()

	// Mod -> Accept button
	if strings.HasPrefix(data.CustomID, forms2.ModAcceptBtnIdPrefix) {
		formId := strings.TrimPrefix(data.CustomID, forms2.ModAcceptBtnIdPrefix)

		form, ok := config.Yaml().Forms[formId]
		if !ok {
			return
		}

		// receive the last message in the public channel
		messages, err := s.ChannelMessages(form.ChannelId, 1, "", "", "")
		if err != nil {
			util.InteractionRespondError(err, s, i)
			return
		}

		// Send the embed data into the public channel
		sentMessage, err := s.ChannelMessageSendComplex(form.ChannelId, &discordgo.MessageSend{
			Embeds: i.Message.Embeds,
		})
		if err != nil {
			msg := fmt.Errorf("failed to send the embeded message into the public channel: %s", err)
			util.InteractionRespondError(msg, s, i)
			return
		}

		// send the interactive form button again (since we deleted the last one)
		err = forms2.SendFormButton(s, form.ChannelId, formId, form.ButtonLabel)
		if err != nil {
			msg := fmt.Errorf("failed to send the interactive form button again (after deleting): %s", err)
			util.InteractionRespondError(msg, s, i)
			return
		}

		// at the point, since we know a new interactive form button is sent into the public channel
		// so it is safe to delete the old message that has the interactive form button
		if len(messages) == 1 && forms2.DoesHaveButtonComponentWithLabel(messages[0], form.ButtonLabel) {
			_ = s.ChannelMessageDelete(form.ChannelId, messages[0].ID)
		}

		userId := i.Member.User.ID
		content := fmt.Sprintf("Accepted by <@%s>, <t:%d:R>", userId, time.Now().UTC().Unix())

		_, err = s.ChannelMessageCrosspost(sentMessage.ChannelID, sentMessage.ID)
		if err != nil {
			content += "\nCross posting the message to the followers failed ❌: " + err.Error()
		} else {
			content += "\nCross posted the message to the followers ✅"
		}

		msgEdit := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)
		msgEdit.SetContent(content)

		// remove the message components
		emptyComponent := make([]discordgo.MessageComponent, 0)
		msgEdit.Components = &emptyComponent

		_, err = s.ChannelMessageEditComplex(msgEdit)
		if err != nil {
			util.InteractionRespondError(err, s, i)
		}
		return
	}

	// Mod -> Reject button
	if strings.HasPrefix(data.CustomID, forms2.ModRejectBtnIdPrefix) {
		//formId := strings.TrimPrefix(data.CustomID, forms.ModRejectBtnIdPrefix)

		msgEdit := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)
		userId := i.Member.User.ID
		content := fmt.Sprintf("Rejected by <@%s>, <t:%d:R>\n", userId, time.Now().UTC().Unix())
		msgEdit.SetContent(content)

		// remove the message components
		emptyComponent := make([]discordgo.MessageComponent, 0)
		msgEdit.Components = &emptyComponent

		_, err := s.ChannelMessageEditComplex(msgEdit)
		if err != nil {
			util.InteractionRespondError(err, s, i)
		}
		return
	}

	// Mod -> Ban button
	if strings.HasPrefix(data.CustomID, forms2.ModBanBtnIdPrefix) {
		userIdToBan := strings.TrimPrefix(data.CustomID, forms2.ModBanBtnIdPrefix)

		banReason := fmt.Sprintf("Banned by %s (%s)",
			i.Member.User.GlobalName,
			i.Member.User.Username,
		)

		err := s.GuildBanCreateWithReason(i.GuildID, userIdToBan, banReason, 7)
		if err != nil {
			util.InteractionRespondError(fmt.Errorf("failed to ban the given user: %s", err), s, i)
			return
		}

		msgEdit := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)
		userId := i.Member.User.ID
		content := fmt.Sprintf("Banned by <@%s>, <t:%d:R>\n", userId, time.Now().UTC().Unix())
		msgEdit.SetContent(content)

		// remove the message components
		emptyComponent := make([]discordgo.MessageComponent, 0)
		msgEdit.Components = &emptyComponent

		_, err = s.ChannelMessageEditComplex(msgEdit)
		if err != nil {
			util.InteractionRespondError(err, s, i)
		}
		return
	}

	// Form interactive button (to open modal)
	if !strings.HasPrefix(data.CustomID, forms2.FormButtonIdPrefix) {
		return // not a form interactive button, skip it
	}

	inputFormId := strings.TrimPrefix(data.CustomID, forms2.FormButtonIdPrefix)

	form, ok := config.Yaml().Forms[inputFormId]
	if !ok {
		return // form does not exist in the YAML config
	}

	var components []discordgo.MessageComponent
	for _, input := range form.Inputs {
		style := discordgo.TextInputShort
		if input.Multiline {
			style = discordgo.TextInputParagraph
		}

		label := lo.Capitalize(input.Id)

		comp := discordgo.TextInput{
			CustomID:    input.Id,
			Label:       label,
			Style:       style,
			Placeholder: input.Placeholder,
			Required:    input.Required,
			MaxLength:   input.Max,
			MinLength:   input.Min,
		}

		actionsRow := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				comp,
			},
		}

		components = append(components, actionsRow)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			//CustomID:   forms.FormModalIdPrefix + i.Interaction.Member.User.ID,
			CustomID:   forms2.FormModalIdPrefix + inputFormId,
			Title:      form.Title,
			Components: components,
		},
	})
	if err != nil {
		util.InteractionRespondError(err, s, i)
	}
}

func handleInteractionModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()

	modalCustomId := data.CustomID
	if !strings.HasPrefix(modalCustomId, forms2.FormModalIdPrefix) {
		return // skip it
	}

	formId := strings.TrimPrefix(modalCustomId, forms2.FormModalIdPrefix)
	form, ok := config.Yaml().Forms[formId]
	if !ok {
		return // skip invalid forms
	}

	// map of the 'input id' to the 'user given value'
	var userInput []forms2.UserInput
	for _, parentComp := range data.Components {
		switch row := parentComp.(type) {
		case *discordgo.ActionsRow:
			for _, childComp := range row.Components {
				switch child := childComp.(type) {
				case *discordgo.TextInput:
					val := strings.TrimSpace(child.Value)       // basic space trim
					val = strings.ReplaceAll(val, "\n\n", "\n") // remove double next lines
					val = strings.ReplaceAll(val, "\t", " ")    // replace the tab character
					val = strings.ReplaceAll(val, "  ", " ")    // remove double spaces
					val = forms2.ConvertLinksToHyperlinks(val)

					if val == "" {
						continue
					}

					userInput = append(userInput, forms2.UserInput{
						InputId: child.CustomID,
						Value:   val,
					})
				}
			}
		}
	}

	// safety check, in case Discord updated its response
	if len(userInput) == 0 {
		util.InteractionRespondError(errors.New("failed to extract the modal values! please report this issue"), s, i)
		return
	}

	// generate embed
	embed := &discordgo.MessageEmbed{
		Color: form.Color,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    form.Footer,
			IconURL: i.Message.Author.AvatarURL(""), // the interaction author is the bot itself!
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: i.Member.User.AvatarURL("100"),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Discord Name",
				Value:  i.Member.User.GlobalName,
				Inline: true,
			},
			{
				Name:   "Username",
				Value:  "`" + util.FormatUsername(i.Member.User) + "`",
				Inline: true,
			},
		},
	}

	if authorAccCreatedAt, err := discordgo.SnowflakeTimestamp(i.Member.User.ID); err == nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Account Created",
			//Value: humanize.Time(authorAccCreatedAt),
			Value:  fmt.Sprintf("<t:%d:R>", authorAccCreatedAt.UTC().Unix()),
			Inline: false,
		})
	}

	// append user inputs
	for _, inp := range userInput {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   lo.Capitalize(inp.InputId),
			Value:  inp.Value,
			Inline: false,
		})
	}

	channelId := form.ChannelId
	message := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	}
	respContent := "Thank you for taking your time to fill this form. ✅"

	if !form.ModSkipApproval {
		channelId = form.ModChannelId // overwrite the public channel with the private one.

		// the user who fills the modal should know their request is going to be in a pending state.
		respContent += "\n\nModerators will review your request shortly. 🔎"

		message.Components = append(message.Components, forms2.GenerateModComponents(formId, i.Member.User.ID))
	}

	sentMessage, err := s.ChannelMessageSendComplex(channelId, message)
	if err != nil {
		util.InteractionRespondError(fmt.Errorf("failed to send form embed: %s", err), s, i)
		return
	}

	// publish the message
	if form.ModSkipApproval {
		_, _ = s.ChannelMessageCrosspost(sentMessage.ChannelID, sentMessage.ID)
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: respContent,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func handleInteractionApplicationCommandAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name
	log.Debug().Str("name", name).
		Str("guild-id", i.GuildID).
		Str("channel-id", i.ChannelID).
		Interface("options", i.ApplicationCommandData().Options).
		Msg("event: interaction application command autocomplete")

	switch name {
	case reference.Name:
		reference.AutocompleteHandler(s, i)
		return
	case google_it.Name:
		google_it.AutocompleteHandler(s, i)
		return
	}
}

func checkRateLimit(userID string) bool {
	if util.IsUserModerator(userID) {
		return false
	}

	cache.RateLimit.IncrementUsage(userID)

	maxUsage := config.Yaml().RateLimit.Usage
	return cache.RateLimit.GetUsageCount(userID) > maxUsage
}
