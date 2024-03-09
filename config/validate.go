package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/nestjs-discord/utility-bot/pkg/colors"
	"github.com/nestjs-discord/utility-bot/pkg/validation"
)

// validate runs both environment variable and YAML validations.
func validate() {
	validateEnvVars()
	validateYaml()
}

// validateEnvVars checks if required environment variables are set.
func validateEnvVars() {
	for _, v := range toValidate {
		if os.Getenv(v) != "" {
			continue
		}

		log.Fatalf("config: '%s' environment variable is required", v)
	}
}

// validateYaml validates the configuration struct using the validator library.
func validateYaml() {
	err := validation.Validate.Struct(c)
	if err == nil {
		return
	}

	// this check is only needed when your code could produce
	// an invalid value for validation such as interface with nil
	// value, most including myself do not usually have code like this.
	var invalidValidationError *validator.InvalidValidationError
	if errors.As(err, &invalidValidationError) {
		log.Fatal(fmt.Errorf("invalid validation error: %s", err))
		return
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		log.Fatal(fmt.Errorf("invalid validation error: %s", err))
		return
	}

	// Print detailed validation errors.
	divider := strings.Repeat("-", 50)

	for _, err := range err.(validator.ValidationErrors) {
		fmt.Println(divider)

		msg := colors.Red.Render(err.Translate(validation.Trans))
		msg += fmt.Sprintf("\n struct namespace:\t%s", colors.Yellow.Render(err.StructNamespace()))
		msg += fmt.Sprintf("\n validation tag:\t%s", colors.Blue.Render(err.Tag()))
		msg += fmt.Sprintf("\n current value:\t\t'%s'", colors.Purple.Render(fmt.Sprint(err.Value())))
		// msg += fmt.Sprintf("\n\traw error message: %s", err.Error())

		fmt.Println(msg)
	}

	fmt.Println(divider)

	log.Fatal(colors.Red.Render("config validation failed"))
}
