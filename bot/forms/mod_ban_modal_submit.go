package forms

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
	"strings"
	"time"
)

func (f *Forms) ModBanModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate, customId *components.CustomID) error {
	data := i.ModalSubmitData()
	val := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	if strings.ToLower(val) != "yes" {
		respond.InteractionWithEphemeralMessage(s, i, "Submit 'yes' to confirm the action.")
		return nil
	}

	userIdToBan := customId.UserId

	banReason := fmt.Sprintf("Banned by %s (%s)",
		i.Member.User.GlobalName,
		i.Member.User.Username,
	)

	err := s.GuildBanCreateWithReason(i.GuildID, userIdToBan, banReason, 7)
	if err != nil {
		return fmt.Errorf("failed to ban the given user: %s", err)
	}

	msgEdit := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)
	content := fmt.Sprintf("Banned by %s, <t:%d:R>\n",
		i.Member.User.Mention(),
		time.Now().UTC().Unix(),
	)
	msgEdit.SetContent(content)

	// remove the message components
	emptyComponent := make([]discordgo.MessageComponent, 0)
	msgEdit.Components = &emptyComponent

	_, err = s.ChannelMessageEditComplex(msgEdit)
	if err != nil {
		return fmt.Errorf("failed to edit the message: %s", err)
	}

	respond.InteractionWithEphemeralMessage(s, i, "Successfully banned the user.")

	return nil
}
