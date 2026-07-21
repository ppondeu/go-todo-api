package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func nullableTodoStateID(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}

	_, err := uuid.Parse(value)

	return err == nil
}

func nullableDueDate(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}

	_, err := time.Parse(time.RFC3339, value)
	return err == nil
}
