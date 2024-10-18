package markdown

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
	"time"
)

const (
	AcknowledgeButtonAction = "content-ack"
)

func (m *Markdown) generateAckButton(user *dgo.User) (dgo.MessageComponent, error) {
	ackCustomId, err := components.EncodeCustomId(&components.CustomID{
		Action: AcknowledgeButtonAction,
		UserId: user.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding custom id: %s", err)
	}

	ackButton := dgo.Button{
		Label:    "After reading things above ☝️, please click here to acknowledge. ✅",
		Style:    dgo.SecondaryButton,
		Disabled: false,
		CustomID: ackCustomId,
		// Emoji:    &dgo.ComponentEmoji{Name: "✅", Animated: false},
	}

	return ackButton, nil
}

func (m *Markdown) AcknowledgeButtonClicked(s *dgo.Session, i *dgo.InteractionCreate, customId *components.CustomID) error {
	whoIsAcknowledging := i.Member.User.ID
	whoIsMentionedToAck := customId.UserId

	if whoIsAcknowledging != whoIsMentionedToAck {
		warnMessage := "Only the mentioned user on top can acknowledge this message. ⚠️\n\n"
		warnMessage += "-# Abuse of this functionality will lead to a ban. 🔴"
		respond.InteractionWithEphemeralMessage(s, i, warnMessage)
		return nil
	}

	content := fmt.Sprintf(
		"This message has been acknowledged by <@%s>, <t:%d:R>",
		whoIsAcknowledging,
		time.Now().UTC().Unix(),
	)

	msgEdit := dgo.NewMessageEdit(i.ChannelID, i.Message.ID)
	msgEdit.SetContent(content)

	// remove the message components
	emptyComponent := make([]dgo.MessageComponent, 0)
	msgEdit.Components = &emptyComponent

	_, err := s.ChannelMessageEditComplex(msgEdit)
	if err != nil {
		return fmt.Errorf("channel message edit complex failed: %s", err)
	}

	return nil
}
