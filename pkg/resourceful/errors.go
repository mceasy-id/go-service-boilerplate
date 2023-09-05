package resourceful

import (
	"errors"
	"fmt"
	"strings"
)

var ErrPagination = errors.New("invalid pagination parameter")

type ErrInit struct {
	Message string
}

func (e *ErrInit) Error() string {
	return fmt.Sprintf("goresourceful: %s", e.Message)
}

func createError(message string) error {
	return &ErrInit{Message: message}
}

// Validation Error
type ValidationErrors []FieldError

func (v ValidationErrors) Error() string {
	var errMessages []string
	for _, field := range v {
		errMessage := fmt.Sprintf("%s: ", field.FieldName)

		for _, err := range field.Errors {
			errMessage += err
		}

		errMessages = append(errMessages, errMessage)
	}

	return strings.Join(errMessages, ", ")
}

type FieldError struct {
	FieldName string
	Errors    []string
}

func (e *ValidationErrors) appendFieldError(fieldName string, message string) {
	*e = append(*e, FieldError{
		FieldName: fieldName,
		Errors:    []string{message},
	})
}
