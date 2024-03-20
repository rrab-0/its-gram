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

type ValidationErrors []error

func (errs ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range errs {
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}

func GenerateRequestValidatorError(err error) error {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	validationErr := validationErrs[0]
	fieldName := validationErr.Field()
	val := validationErr.Value()
	fieldName = strings.ToLower(fieldName)

	tag := validationErr.Tag()

	if tag == "required" {
		return fmt.Errorf("%s is required in the request", fieldName)
	}

	if tag == "email" {
		return fmt.Errorf("%s is not a valid email address", fieldName)
	}

	if tag == "uuid" {
		return fmt.Errorf("%v is not a valid uuid", val)
	}

	return err
}
