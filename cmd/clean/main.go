package main

import (
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"log"
)

func main() {
	discordCfg, err := env.NewDiscordConfig()
	if err != nil {
		log.Fatal(err)
	}

	// handler := handler.NewHandler()
	b, err := bot.NewBot(discordCfg)
	if err != nil {
		log.Fatal(err)
	}

	err = b.CleanApplicationCommands()
	if err != nil {
		log.Fatal(err)
	}
}
