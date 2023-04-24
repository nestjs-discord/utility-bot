package main

import (
	"github.com/erosdesire/discord-nestjs-utility-bot/cmd"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	_ = godotenv.Load()

	if err := cmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("execution failed")
	}
}
