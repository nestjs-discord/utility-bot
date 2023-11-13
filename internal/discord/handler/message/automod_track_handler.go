package message

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/internal/cache"
)

func AutoModTrackHandler(s *discordgo.Session, i *discordgo.MessageCreate) {
	content := "AutoMod is tracking the following text channels.\n"
	content += "> Forum channels and posts within them are ignored.\n"
	for _, channelId := range cache.AutoMod.GetTrackedChannelIds() {
		content += fmt.Sprintf("- <#%s>\n", channelId)
	}

	_, _ = s.ChannelMessageSend(i.ChannelID, content)
}
