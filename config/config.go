package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var c Yaml

func GetYaml() *Yaml {
	return &c
}

func Bootstrap(configFilePath string) {
	setViperPath(configFilePath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("viper: config read failed")
	}

	log.Debug().Str("path", viper.ConfigFileUsed()).Msg("config: read success")

	if err := viper.UnmarshalExact(&c); err != nil {
		log.Fatal().Err(err).Msg("config: unmarshal failed")
	}

	validate()
}

func setViperPath(path string) {
	if path != "" {
		viper.SetConfigFile(path)
		return
	}

	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
}
