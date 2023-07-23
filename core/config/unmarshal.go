package config

import (
	"github.com/spf13/viper"
)

func Unmarshal() error {
	err := viper.UnmarshalExact(&c)
	if err != nil {
		return err
	}

	return nil
}
