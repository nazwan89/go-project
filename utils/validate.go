package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// pathParam wraps a path parameter value for struct-tag-based validation.
// Rules: required (non-empty), min 1 char, max 100 chars, alphanumeric only.
type pathParam struct {
	Value string `validate:"required,min=1,max=100,alphanum"`
}

// queryParam wraps a query string value for struct-tag-based validation.
// Rules: optional (omitempty), but if present: min 1 char, max 200 chars.
type queryParam struct {
	Value string `validate:"omitempty,min=1,max=200"`
}

// ValidatePathParam validates a path parameter value.
// Returns a non-nil error if the value is empty, exceeds 100 characters,
// or contains non-alphanumeric characters.
// Call this before using any c.Params() value in business logic.
func ValidatePathParam(value string) error {
	return validate.Struct(pathParam{Value: value})
}

// ValidateQueryParam validates a query string parameter value.
// Empty values are allowed (optional query params). Non-empty values must
// not exceed 200 characters.
// Call this before using any c.Query() value in business logic.
func ValidateQueryParam(value string) error {
	if value == "" {
		return nil
	}
	return validate.Struct(queryParam{Value: value})
}
