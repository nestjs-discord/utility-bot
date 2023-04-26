package interaction

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/erosdesire/discord-nestjs-utility-bot/internal/discord/util"
)

func NpmInspectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: logic
	util.InteractionRespondError(errors.New("not implemented yet"), s, i)
}
