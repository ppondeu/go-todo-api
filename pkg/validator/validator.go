package validator

import v "github.com/go-playground/validator/v10"

func NewValidator() *v.Validate {
	valivator := v.New()
	valivator.RegisterValidation("nullable_todo_state_id", nullableTodoStateID)
	valivator.RegisterValidation("nullable_due_date", NullableDueDate)
	return valivator
}
