package validation

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	eng "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate = validator.New(validator.WithRequiredStructEnabled())
	Trans    ut.Translator
)

func init() {
	// Register a custom validation rule
	err := Validate.RegisterValidation("max-one-space-allowed", MaxOneSpaceValidator)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to register custom config validation function: %s", err))
		return
	}

	uni := ut.New(en.New())
	trans, _ := uni.GetTranslator("en")
	if err := eng.RegisterDefaultTranslations(Validate, trans); err != nil {
		log.Fatal(fmt.Errorf("failed to register default validation translations: %s", err))
	}

	Trans = trans
}

// MaxOneSpaceValidator validates that a string field contains at most one optional space character
var MaxOneSpaceValidator = func(fl validator.FieldLevel) bool {
	for _, key := range fl.Field().MapKeys() {
		if strings.Count(key.String(), " ") > 1 {
			return false
		}
	}

	return true
}
