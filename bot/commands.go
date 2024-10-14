package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) CleanApplicationCommands() error {
	emptyCmd := make([]*discordgo.ApplicationCommand, 0)
	_, err := b.session.ApplicationCommandBulkOverwrite(
		b.cfg.AppId,
		b.cfg.GuildId,
		emptyCmd,
	)
	if err != nil {
		return fmt.Errorf("unable to clean the app commands: %v", err)
	}

	b.logger.Info("cleaned application commands")

	return nil
}

func (b *Bot) RegisterApplicationCommands() error {
	// TODO: logic
	return nil
}
