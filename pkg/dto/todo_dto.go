package dto

type CreateTodoDto struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateTodoDto struct {
	Title       *string `json:"title" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty"`
	DueDate     *string `json:"due_date" validate:"omitempty"`
	Priority    *string `json:"priority" validate:"omitempty,oneof=high medium low"`
	State       *string `json:"state" validate:"omitempty,oneof=not_started in_progress done"`
	CategoryID  *string `json:"category_id" validate:"omitempty,uuid"`
}
