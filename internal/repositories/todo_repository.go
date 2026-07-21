package repositories

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/application/ports"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) ports.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Save(ctx context.Context, userID uuid.UUID, todo *domain.Todo) (*domain.Todo, error) {
	todo.UserID = userID
	if err := r.db.WithContext(ctx).Create(todo).Error; err != nil {
		return nil, err
	}

	return r.FindByID(ctx, userID, todo.ID)
}

func (r *todoRepository) FindByID(ctx context.Context, userID, todoID uuid.UUID) (*domain.Todo, error) {
	var todo domain.Todo
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", todoID, userID).
		Preload("State").
		First(&todo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrTodoNotFound
	}
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *todoRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Todo, error) {
	todos := make([]domain.Todo, 0)
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_deleted = ?", userID, false).
		Preload("State").
		Order("created_at DESC").
		Find(&todos).Error
	return todos, err
}

func (r *todoRepository) List(ctx context.Context, userID uuid.UUID, filter ports.TodoListFilter) (*ports.TodoPage, error) {
	var total int64
	if err := r.filteredQuery(ctx, userID, filter).Count(&total).Error; err != nil {
		return nil, err
	}

	todos := make([]domain.Todo, 0)
	offset := (filter.Page - 1) * filter.Limit
	orderBy := allowedTodoSorts[filter.Sort] + " " + strings.ToUpper(filter.Order)
	if err := r.filteredQuery(ctx, userID, filter).
		Preload("State").
		Order(orderBy).
		Offset(offset).
		Limit(filter.Limit).
		Find(&todos).Error; err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
	}

	return &ports.TodoPage{
		Items:      todos,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

var allowedTodoSorts = map[string]string{
	"created_at": "created_at",
	"updated_at": "updated_at",
	"due_date":   "due_date",
	"priority":   "priority",
	"title":      "title",
}

func (r *todoRepository) filteredQuery(ctx context.Context, userID uuid.UUID, filter ports.TodoListFilter) *gorm.DB {
	query := r.db.WithContext(ctx).Model(&domain.Todo{}).
		Where("user_id = ? AND is_deleted = ?", userID, false)

	if filter.Search != "" {
		pattern := "%" + escapeLike(filter.Search) + "%"
		query = query.Where("(title ILIKE ? ESCAPE E'\\\\' OR description ILIKE ? ESCAPE E'\\\\')", pattern, pattern)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}
	if filter.StateID != nil {
		query = query.Where("state_id = ?", *filter.StateID)
	}
	if filter.DueFrom != nil {
		query = query.Where("due_date >= ?", *filter.DueFrom)
	}
	if filter.DueTo != nil {
		query = query.Where("due_date <= ?", *filter.DueTo)
	}

	return query
}

func escapeLike(value string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(value)
}

func (r *todoRepository) Update(ctx context.Context, userID, todoID uuid.UUID, fields map[string]any) (*domain.Todo, error) {
	result := r.db.WithContext(ctx).Model(&domain.Todo{}).
		Where("id = ? AND user_id = ?", todoID, userID).
		Updates(fields)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrTodoNotFound
	}

	return r.FindByID(ctx, userID, todoID)
}

func (r *todoRepository) Delete(ctx context.Context, userID, todoID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", todoID, userID).Delete(&domain.Todo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrTodoNotFound
	}

	return nil
}

func (r *todoRepository) FindStateByID(ctx context.Context, userID, stateID uuid.UUID) (*domain.TodoState, error) {
	var state domain.TodoState
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", stateID, userID).First(&state).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrTodoStateNotFound
	}
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func (r *todoRepository) UpdateTodoState(ctx context.Context, userID, stateID uuid.UUID, fields map[string]any) (*domain.TodoState, error) {
	result := r.db.WithContext(ctx).Model(&domain.TodoState{}).
		Where("id = ? AND user_id = ?", stateID, userID).
		Updates(fields)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrTodoStateNotFound
	}

	return r.FindStateByID(ctx, userID, stateID)
}

func (r *todoRepository) InitTodoState(userID uuid.UUID) ([]domain.TodoState, error) {
	states := []string{domain.Backlog, domain.NotStarted, domain.InProgress, domain.Done}
	todoStates := make([]domain.TodoState, 0, len(states))
	for _, state := range states {
		todoStates = append(todoStates, domain.TodoState{Name: state, UserID: userID})
	}

	if err := r.db.Create(todoStates).Error; err != nil {
		return nil, err
	}

	return todoStates, nil
}
