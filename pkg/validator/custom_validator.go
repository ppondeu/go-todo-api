package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/pkg/logs"
)

func nullableCategory(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	logs.Info(value)
	if value == "" {
		return true 
	}

	_, err := uuid.Parse(value)

	return err == nil

}

func NullableDueDate(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	logs.Info(value)
	if value == "" {
		return true 
	}

	_, err := time.Parse(time.RFC3339, value)
	return err == nil

}

