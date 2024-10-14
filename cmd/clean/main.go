package main

import (
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/config"
	"log"
)

func main() {
	botCfg, err := config.NewBotConfig()
	if err != nil {
		log.Fatal(err)
	}

	// handler := handler.NewHandler()
	b, err := bot.NewBot(botCfg, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = b.CleanApplicationCommands()
	if err != nil {
		log.Fatal(err)
	}
}
