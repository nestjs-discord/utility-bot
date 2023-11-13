package cache

import (
	"github.com/nestjs-discord/utility-bot/internal/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler/automod"
)

var AutoMod *automod.AutoMod

func InitAutoMod() {
	AutoMod = automod.NewAutoMod(automod.Option{
		MessageTTL: config.GetConfig().AutoMod.MessageTTL,
		DenyTTL:    config.GetConfig().AutoMod.DenyTTL,
	})
}
