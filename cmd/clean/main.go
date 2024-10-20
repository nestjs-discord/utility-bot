package main

import (
	"fmt"
	dgo "github.com/bwmarrin/discordgo"
	"github.com/nestjs-discord/utility-bot/bot/commands"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"log"
)

func main() {
	discordCfg, err := env.NewDiscordConfig()
	if err != nil {
		log.Fatal(err)
	}

	session, err := dgo.New("Bot " + discordCfg.Token)
	if err != nil {
		log.Fatalf("unable to create the session: %v", err)
	}

	err = commands.CleanApplicationCommands(
		session,
		discordCfg.AppId,
		discordCfg.GuildId.String(),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("cleaned application commands")
}
