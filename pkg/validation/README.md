# Validation

The `validation` package provides functionality for validating data using
the Go [`go-playground/validator`](<github.com/go-playground/validator>) package.

- It initializes the validator instance with required struct enabled.
- Registers a custom validation rule, `max-one-space-allowed`,
which ensures that a string field contains at most one optional space character.
- Handles translations for validation error messages.
