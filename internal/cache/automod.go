package cache

import (
	"github.com/nestjs-discord/utility-bot/bot/automod"
	"github.com/nestjs-discord/utility-bot/config"
)

var AutoMod *automod.AutoMod // TODO: avoid global instance

func InitAutoMod() {
	AutoMod = automod.NewAutoMod(automod.Option{
		MessageTTL: config.Yaml().AutoMod.MessageTTL,
		DenyTTL:    config.Yaml().AutoMod.DenyTTL,
	})
}
