package config

import "github.com/nestjs-discord/utility-bot/config/yaml"

var c yaml.Config // TODO: remove

func Yaml() *yaml.Config { // TODO: remove
	return &c
}
