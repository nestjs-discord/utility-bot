package cache

import (
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler/automod"
)

var AutoMod *automod.AutoMod

func InitAutoMod() {
	AutoMod = automod.NewAutoMod(automod.Option{
		MessageTTL: config.Yaml().AutoMod.MessageTTL,
		DenyTTL:    config.Yaml().AutoMod.DenyTTL,
	})
}
