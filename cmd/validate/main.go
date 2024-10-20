package main

import (
	"github.com/nestjs-discord/utility-bot/infra/config/env"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log"
	"log/slog"
)

func main() {
	_, err := logger.NewLogger(env.StageDev)
	if err != nil {
		log.Fatal(err)
	}

	_, err = yaml.NewConfig("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Good job! everything looks fine :)")
}
