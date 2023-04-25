package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func validateConfig() {
	validate := validator.New()

	err := validate.Struct(c)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Fatal().Err(err).Msg("config: validation invalid error")
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			log.Error().
				Str("namespace", err.Namespace()).
				Str("field", err.Field()).
				Str("tag", err.Tag()).
				Interface("given-value", err.Value()).
				Msg("config: validation field")
		}

		log.Fatal().Msg("config: validation error, please check the configuration!")
	}
}
