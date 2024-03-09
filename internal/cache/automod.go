package cache

import (
	"github.com/nestjs-discord/utility-bot/config"
	"github.com/nestjs-discord/utility-bot/internal/discord/handler/automod"
)

var AutoMod *automod.AutoMod

func InitAutoMod() {
	AutoMod = automod.NewAutoMod(automod.Option{
		MessageTTL: config.GetYaml().AutoMod.MessageTTL,
		DenyTTL:    config.GetYaml().AutoMod.DenyTTL,
	})
}
