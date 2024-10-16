package forms

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/components"
	"github.com/nestjs-discord/utility-bot/bot/user"
	"time"
)

func (f *Forms) ModBanButtonClicked(s *discordgo.Session, i *discordgo.InteractionCreate, customId *components.CustomID) error {
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
		user.Mention(i.Member.User),
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

	return nil
}
