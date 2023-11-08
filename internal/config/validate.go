package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
)

// MaxOneSpaceValidator validates that a string field contains at most one optional space character
var MaxOneSpaceValidator = func(fl validator.FieldLevel) bool {
	for _, key := range fl.Field().MapKeys() {
		if strings.Count(key.String(), " ") > 1 {
			return false
		}
	}

	return true
}

func ValidateConfig() error {
	validate := validator.New()

	// register all sql.Null* types to use the ValidateValuer CustomTypeFunc
	err := validate.RegisterValidation("max-one-space-allowed", MaxOneSpaceValidator)
	if err != nil {
		return errors.Wrap(err, "failed to register custom config validation function")
	}

	err = validate.Struct(c)
	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value, most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.Wrap(err, "invalid validation error")
		}

		for _, err := range err.(validator.ValidationErrors) {
			log.Error().
				Str("namespace", err.Namespace()).
				Str("field", err.Field()).
				Str("tag", err.Tag()).
				Interface("given-value", err.Value()).
				Msg("config: validation field")
		}

		return errors.New("config validation failed, please check the configuration file")
	}

	return nil
}
