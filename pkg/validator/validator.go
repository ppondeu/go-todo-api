package validator

import v "github.com/go-playground/validator/v10"

func NewValidator() *v.Validate {
	validator := v.New()
	if err := validator.RegisterValidation("nullable_todo_state_id", nullableTodoStateID); err != nil {
		panic(err)
	}
	if err := validator.RegisterValidation("nullable_due_date", nullableDueDate); err != nil {
		panic(err)
	}
	return validator
}
