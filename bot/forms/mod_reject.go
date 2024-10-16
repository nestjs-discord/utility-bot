package forms

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"time"
)

func (f *Forms) ModRejectButtonClicked(s *discordgo.Session, i *discordgo.InteractionCreate, _ *components.CustomID) error {
	// formId := customId.FormId
	if f.RaceConditionCheck(i.Message.ID) {
		f.RaceConditionRespond(s, i)
		return nil
	}

	msgEdit := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)
	userId := i.Member.User.ID
	content := fmt.Sprintf("Rejected by <@%s>, <t:%d:R>\n", userId, time.Now().UTC().Unix())
	msgEdit.SetContent(content)

	// remove the message components
	emptyComponent := make([]discordgo.MessageComponent, 0)
	msgEdit.Components = &emptyComponent

	_, err := s.ChannelMessageEditComplex(msgEdit)
	if err != nil {
		return fmt.Errorf("failed to edit the message: %s", err)
	}

	return nil
}
