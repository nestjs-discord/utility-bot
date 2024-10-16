package main

import (
	"flag"
	"github.com/nestjs-discord/utility-bot/infra/config/yaml"
	"github.com/nestjs-discord/utility-bot/infra/logger"
	"log"
	"log/slog"
)

var yamlConfigPath = flag.String("yaml-config-path", "./config.yml", "")

func main() {
	err := logger.Initialize("dev")
	if err != nil {
		log.Fatal(err)
	}

	_, err = yaml.NewConfig(yaml.Path(*yamlConfigPath))
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Good job! everything looks fine :)")
}
