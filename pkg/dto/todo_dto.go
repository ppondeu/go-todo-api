package dto

type CreateTodoDto struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateTodoDto struct {
    Title       *string     `json:"title" validate:"omitempty,min=3,max=100"`
    Description *string     `json:"description" validate:"omitempty,max=500"`
    Priority    *string     `json:"priority" validate:"omitempty,oneof=high medium low"`
    State       *string     `json:"state" validate:"omitempty,oneof=not_started in_progress done"`
    CategoryID  *string     `json:"category_id" validate:"omitempty,nullable_category"`
    DueDate     *string    `json:"due_date" validate:"omitempty,nullable_due_date"`
}

type CreateCategoryDto struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCategoryDto struct {
	Name string `json:"name" validate:"required"`
}
