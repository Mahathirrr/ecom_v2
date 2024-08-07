package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func getCustomValidationError(err validator.FieldError) ValidationError {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return ValidationError{Field: field, Error: "This field is required"}
	case "email":
		return ValidationError{Field: field, Error: "Invalid email format"}
	case "min":
		if field == "Password" {
			return ValidationError{Field: field, Error: fmt.Sprintf("Password must be at least %s characters long", param)}
		}
		return ValidationError{Field: field, Error: fmt.Sprintf("Minimum value is %s", param)}
	case "max":
		return ValidationError{Field: field, Error: fmt.Sprintf("Maximum value is %s", param)}
	case "url":
		return ValidationError{Field: field, Error: "Invalid URL format"}
	case "gt":
		return ValidationError{Field: field, Error: fmt.Sprintf("Value must be greater than %s", param)}
	case "oneof":
		return ValidationError{Field: field, Error: fmt.Sprintf("Value must be one of: %s", param)}
	default:
		return ValidationError{Field: field, Error: fmt.Sprintf("Failed on %s validation", tag)}
	}
}

func ValidateStruct(v *validator.Validate, s interface{}) []ValidationError {
	var errors []ValidationError

	err := v.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, getCustomValidationError(err))
		}
	}

	return errors
}
