package main

import (
	"github.com/joho/godotenv"
	"github.com/nestjs-discord/utility-bot/cmd"
	"github.com/nestjs-discord/utility-bot/internal/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	_ = godotenv.Load()
	logger.Register()

	if err := cmd.Execute(); err != nil {
		log.Panic().Err(err).Msg("execution failed")
	}
}
