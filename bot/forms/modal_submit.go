package forms

import (
	"errors"
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/forPelevin/gomoji"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
	"github.com/nestjs-discord/utility-bot/bot/security"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"strings"
)

func (f *Forms) getInputLabelByInputId(inputs []yaml.FormInput, id string) string {
	for _, formInput := range inputs {
		if formInput.Id == id {
			return formInput.Label
		}
	}

	return "Label not found!"
}

func (f *Forms) ModalSubmitted(s *dgo.Session, i *dgo.InteractionCreate, customId *components.CustomID) error {
	form, err := f.getFormById(customId.FormId)
	if err != nil {
		return err
	}

	data := i.ModalSubmitData()

	// map of the 'input id' to the 'user given value'
	var userInputs []userInputType
	for _, parentComp := range data.Components {
		row, isRow := parentComp.(*dgo.ActionsRow)
		if !isRow {
			break
		}

		for _, childComp := range row.Components {
			child, isText := childComp.(*dgo.TextInput)
			if !isText {
				break
			}

			userInput := strings.TrimSpace(child.Value)                // basic space trim
			userInput = strings.ReplaceAll(userInput, "\n\n", "\n")    // remove double next lines
			userInput = strings.ReplaceAll(userInput, "\t", " ")       // replace the tab character
			userInput = strings.ReplaceAll(userInput, "  ", " ")       // remove double spaces
			userInput = gomoji.RemoveEmojis(userInput)                 // we don't want emojis
			userInput = f.markdown.ConvertLinksToHyperlinks(userInput) // improves embed visualization
			userInput = security.RemoveDangerousMentions(userInput)

			if userInput == "" {
				continue
			}

			// Discord only returns us the input CustomID
			// We have to manually loop through the config file to fetch the input label
			inputId := child.CustomID
			inputLabel := f.getInputLabelByInputId(form.Inputs, inputId)

			userInputs = append(userInputs, userInputType{
				InputId:      inputId,
				InputLabel:   inputLabel,
				InputValue:   userInput,
				AutoModCheck: f.autoMod.ValidateUserInputAgainstServerRules(userInput),
			})
		}
	}

	// safety check, in case Discord updated its response
	if len(userInputs) == 0 {
		return errors.New("failed to extract the modal values! please report this issue")
	}

	// generate embed
	discordProfileEmbed := &dgo.MessageEmbed{
		// Title: "Discord profile",
		Title: i.Member.User.GlobalName,
		// Color: i.Member.User.AccentColor,
		Color: form.Color,
		Thumbnail: &dgo.MessageEmbedThumbnail{
			URL: i.Member.User.AvatarURL("4096"),
		},
		Fields: []*dgo.MessageEmbedField{
			{
				Name:   "Username",
				Value:  "`" + i.Member.User.String() + "`",
				Inline: true,
			},
			{
				Name:   "Profile",
				Value:  i.Member.User.Mention(),
				Inline: false,
			},
		},
	}

	if authorAccCreatedAt, err := dgo.SnowflakeTimestamp(i.Member.User.ID); err == nil {
		discordProfileEmbed.Fields = append(discordProfileEmbed.Fields, &dgo.MessageEmbedField{
			Name:   "Account created",
			Value:  fmt.Sprintf("<t:%d:R>", authorAccCreatedAt.UTC().Unix()),
			Inline: true,
		})
	}

	if !i.Member.JoinedAt.IsZero() {
		discordProfileEmbed.Fields = append(discordProfileEmbed.Fields, &dgo.MessageEmbedField{
			Name:   "Joined the server",
			Value:  fmt.Sprintf("<t:%d:R>", i.Member.JoinedAt.UTC().Unix()),
			Inline: true,
		})
	}

	formDataEmbed := &dgo.MessageEmbed{
		Color: form.Color,
		Footer: &dgo.MessageEmbedFooter{
			Text:    form.Footer,
			IconURL: i.Message.Author.AvatarURL(""), // the interaction author is the bot itself!
		},
	}

	message := &dgo.MessageSend{}

	// append user inputs
	for _, inp := range userInputs {
		formDataEmbed.Fields = append(formDataEmbed.Fields, &dgo.MessageEmbedField{
			Name:   inp.InputLabel,
			Value:  inp.InputValue,
			Inline: false,
		})

		if inp.AutoModCheck != nil {
			if message.Content == "" {
				message.Content = "### 🚨 Auto Mod rules have been triggered: 🚨\n\n"
			}
			message.Content += fmt.Sprintf(
				"_%s_: ||%s||\n",
				inp.InputLabel,
				inp.AutoModCheck,
			)
		}
	}

	channelId := form.ChannelId

	message.Embeds = []*dgo.MessageEmbed{
		discordProfileEmbed,
		formDataEmbed,
	}

	respContent := "Thank you for taking your time to fill this form. ✅"

	if !form.ModSkipApproval {
		channelId = form.ModChannelId // overwrite the public channel with the private one.

		// the user who fills the modal must know their request is going to be in a pending state.
		respContent += "\n\nModerators will review your request shortly. 🔎"

		modComps, err := f.generateModComponents(customId.FormId, i.Member.User.ID)
		if err != nil {
			return fmt.Errorf("failed to generate mod components: %s", err)
		}

		message.Components = append(message.Components, modComps)
	}

	sentMessage, err := s.ChannelMessageSendComplex(channelId, message)
	if err != nil {
		return fmt.Errorf("failed to send form embed: %s", err)
	}

	// publish the message
	if form.ModSkipApproval {
		_, _ = s.ChannelMessageCrosspost(sentMessage.ChannelID, sentMessage.ID)
	}

	respond.InteractionWithEphemeralMessage(s, i, respContent)

	return nil
}
