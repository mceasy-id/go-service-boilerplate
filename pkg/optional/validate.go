package optional

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type FieldError struct {
	Field  string
	Errors []string
}

type ValidationErrors []FieldError

func (v ValidationErrors) Error() string {
	var errMessage string
	for _, field := range v {
		errMessage += fmt.Sprintf("%s: ", field.Field)

		for _, err := range field.Errors {
			errMessage += err
		}
	}

	return errMessage
}
func ValidateStruct(structure any) error {
	parent := reflect.ValueOf(structure)
	parentType := parent.Type()
	if parent.Kind() != reflect.Struct {
		return errors.New("value is not a struct")
	}

	var validationErr ValidationErrors
	for i := 0; i < parent.NumField(); i++ {

		valuer, ok := parent.Field(i).Interface().(driver.Valuer)
		if !ok {
			continue
		}
		
		set := parent.Field(i).FieldByName("isSet").Bool()
		if !set {
			continue
		}

		fieldTags := strings.Split(parentType.Field(i).Tag.Get("nullable"), ",")
		jsonTag := strings.Split(parentType.Field(i).Tag.Get("json"), ",")[0]
		for _, fieldTag := range fieldTags {
			switch fieldTag {
			case "required":
				err := validateRequired(valuer)
				if err != nil {
					var fieldErr FieldError
					fieldErr.Field = jsonTag
					fieldErr.Errors = append(fieldErr.Errors, err.Error())
					validationErr = append(validationErr, fieldErr)
				}
			case "gte0":
				err := validateGteZero(valuer)
				if err != nil {
					var fieldErr FieldError
					fieldErr.Field = jsonTag
					fieldErr.Errors = append(fieldErr.Errors, err.Error())
					validationErr = append(validationErr, fieldErr)
				}
			case "gt0":
				err := validateGtZero(valuer)
				if err != nil {
					var fieldErr FieldError
					fieldErr.Field = jsonTag
					fieldErr.Errors = append(fieldErr.Errors, err.Error())
					validationErr = append(validationErr, fieldErr)
				}
			case "notnull":
				err := validateNotNull(valuer)
				if err != nil {
					var fieldErr FieldError
					fieldErr.Field = jsonTag
					fieldErr.Errors = append(fieldErr.Errors, err.Error())
					validationErr = append(validationErr, fieldErr)
				}

			}
		}
	}
	if len(validationErr) != 0 {
		return validationErr
	}
	return nil
}

func validateNotNull(field driver.Valuer) error {
	value, err := field.Value()
	if err != nil {
		return err
	}

	if b, ok := value.(bool); ok && !b {
		return nil
	}
	err = notNull(value)
	if err != nil {
		return err
	}
	return nil
}

func validateGteZero(field driver.Valuer) error {
	value, err := field.Value()
	if err != nil {
		return err
	}

	if b, ok := value.(bool); ok && !b {
		return nil
	}

	err = notNegative(value)
	if err != nil {
		return err
	}
	return nil
}
func validateGtZero(field driver.Valuer) error {
	value, err := field.Value()
	if err != nil {
		return err
	}

	if b, ok := value.(bool); ok && !b {
		return nil
	}
	err = gtZero(value)
	if err != nil {
		return err
	}
	return nil
}
func validateRequired(field driver.Valuer) error {
	value, err := field.Value()
	if err != nil {
		return err
	}

	if b, ok := value.(bool); ok && !b {
		return nil
	}

	err = required(value)
	if err != nil {
		return err
	}
	return nil
}
func notNegative(val driver.Value) error {
	switch v := val.(type) {
	case nil:
		return nil
	case int32:
		if v < 0 {
			return errors.New("value must be greater than or equal zero")
		}
	case int64:
		if v < 0 {
			return errors.New("value must be greater than or equal zero")
		}
	case float32:
		if v < 0 {
			return errors.New("value must be greater than or equal zero")
		}
	case float64:
		if v < 0 {
			return errors.New("value must be greater than or equal zero")
		}
	}

	return nil
}
func gtZero(val driver.Value) error {
	switch v := val.(type) {
	case nil:
		return nil

	case int32:
		if v <= 0 {
			return errors.New("value must be greater than zero")
		}
	case int64:
		if v <= 0 {
			return errors.New("value must be greater than zero")
		}
	case float32:
		if v <= 0 {
			return errors.New("value must be greater than zero")
		}
	case float64:
		if v <= 0 {
			return errors.New("value must be greater than zero")
		}
	}
	return nil
}
func required(val driver.Value) error {
	switch v := val.(type) {
	case nil:
		return errors.New("value cannot be null")

	case string:
		if v == "" {
			return errors.New("value cannot be empty string")
		}
	}
	return nil
}
func notNull(field driver.Value) error {
	switch field.(type) {
	case nil:
		return errors.New("value cannot be null")
	}
	return nil
}
