package main

import (
	"fmt"
	"github.com/nestjs-discord/utility-bot/bot/invite"
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"log"
)

func main() {
	discordCfg, err := env.NewDiscordConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(invite.Link(discordCfg.AppId))
}
