package yaml

import (
	"fmt"
	"log/slog"
	"reflect"
)

type withValidateFunc interface {
	validate() error
}

func (c Config) fieldName(f reflect.StructField) string {
	// return f.Name
	return f.Tag.Get("yaml")
}

func (c Config) validate() error {
	v := reflect.ValueOf(c)
	fieldsCount := v.NumField()

	for i := 0; i < fieldsCount; i++ {
		structWithValidateFunc, ok := v.Field(i).Interface().(withValidateFunc)
		if !ok {
			continue
		}
		fieldName := c.fieldName(v.Type().Field(i))
		err := structWithValidateFunc.validate()
		if err != nil {
			return fmt.Errorf("the '%s' key has error: %v", fieldName, err)
		}

		slog.Debug("validated", // TODO: replace with logger instance
			slog.String("field", fieldName),
		)
	}

	return nil
}
