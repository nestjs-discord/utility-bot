package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func Unmarshal() error {
	err := viper.UnmarshalExact(&c)
	if err != nil {
		return err
	}

	log.Debug().Interface("instance", c).Msg("config")

	validateConfig()

	return nil
}
