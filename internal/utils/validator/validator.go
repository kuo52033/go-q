package validator

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func GetValidator() *validator.Validate {
	return validate
}

func FormateValidationError(err error) map[string]interface{} {
	validationErrors := err.(validator.ValidationErrors)

	errors := make(map[string]interface{})

	for _, err := range validationErrors {
		errors[err.Field()] = err.Tag()
		errors["full_error"] = err.Error()
	}

	return errors
}
