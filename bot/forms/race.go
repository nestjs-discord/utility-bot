package forms

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

func (f *Forms) RaceConditionCheck(messageId string) bool {
	f.modActionLock.Lock()
	defer f.modActionLock.Unlock()

	cacheKey := messageId
	cacheValue, cacheHit := f.modActionsCache.Get(cacheKey)
	if cacheHit && cacheValue {
		return true
	}

	f.raceConditionStoreMessageId(messageId)

	return false
}

func (f *Forms) raceConditionStoreMessageId(messageId string) {
	cacheKey := messageId

	cacheTtl := 30 * time.Second
	f.modActionsCache.SetWithTTL(cacheKey, true, 1, cacheTtl)
}

func (f *Forms) RaceConditionRespond(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg := "Race condition detected! 😅\n"
	msg += "Another moderator has already handled this message."
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
