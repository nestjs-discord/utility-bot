package main

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/bot"
	"github.com/nestjs-discord/utility-bot/bot/commands"
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

	err = commands.CleanApplicationCommands(
		bot.ProvideSession(b),
		discordCfg.AppId,
		discordCfg.GuildId,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("cleaned application commands")
}
