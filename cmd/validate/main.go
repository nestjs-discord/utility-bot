package main

import (
	"flag"
	"github.com/nestjs-discord/utility-bot/config/yaml"
	"github.com/nestjs-discord/utility-bot/internal/cache"
	"github.com/nestjs-discord/utility-bot/logger"
	"log"
	"log/slog"
)

var yamlConfigPath = flag.String("yaml-config-path", "./config.yml", "")

func main() {
	err := logger.Initialize("dev")
	if err != nil {
		log.Fatal(err)
	}

	ymlCfg, err := yaml.NewConfig(*yamlConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	err = cache.MarkdownContent(ymlCfg.Commands) // TODO: avoid global instance
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Good job! everything looks fine :)")
}
