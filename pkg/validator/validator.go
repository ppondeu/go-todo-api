package validator

import v "github.com/go-playground/validator/v10"

func NewValidator() *v.Validate {
	valivator := v.New()
	valivator.RegisterValidation("nullable_category", nullableCategory)
	valivator.RegisterValidation("nullable_due_date", NullableDueDate)
	return valivator
}
