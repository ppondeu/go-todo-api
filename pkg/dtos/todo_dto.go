package dtos

type CreateTodoDto struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateTodoDto struct {
	Title       *string `json:"title" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Priority    *string `json:"priority" validate:"omitempty,oneof=high medium low"`
	TodoStateID *string `json:"todo_state_id" validate:"omitempty,nullable_todo_state_id"`
	DueDate     *string `json:"due_date" validate:"omitempty,nullable_due_date"`
}

type UpdateTodoState struct {
	Name string `json:"name" validate:"required"`
}
