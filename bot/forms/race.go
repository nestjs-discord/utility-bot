package forms

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/handler/respond"
	"time"
)

func (f *Forms) raceConditionCheck(messageId string) bool {
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

func (f *Forms) raceConditionRespond(s *discordgo.Session, i *discordgo.InteractionCreate) {
	msg := "Race condition detected! 😅\n"
	msg += "Another moderator has already handled this message."
	respond.InteractionWithEphemeralMessage(s, i, msg)
}
