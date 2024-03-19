package internal

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func GenerateRequestValidatorError(err error) error {
	validationErr, ok := err.(validator.FieldError)
	if !ok {
		return err
	}

	fieldName := validationErr.Field()
	fieldName = strings.ToLower(fieldName)

	tag := validationErr.Tag()

	if tag == "required" {
		return fmt.Errorf("%s is required in the request", fieldName)
	}

	if tag == "email" {
		return fmt.Errorf("%s is not a valid email address", fieldName)
	}

	return err
}
