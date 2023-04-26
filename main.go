package main

import (
	"github.com/joho/godotenv"
	"github.com/nestjs-discord/utility-bot/cmd"
	"github.com/rs/zerolog/log"
)

func main() {
	_ = godotenv.Load()

	if err := cmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("execution failed")
	}
}
