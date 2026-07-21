package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/application/ports"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/utils"
)

type TodoService interface {
	Create(context.Context, uuid.UUID, *dtos.CreateTodoDto) (*domain.Todo, error)
	Update(context.Context, uuid.UUID, uuid.UUID, *dtos.UpdateTodoDto) (*domain.Todo, error)
	UpdateTodoState(context.Context, uuid.UUID, uuid.UUID, *dtos.UpdateTodoState) (*domain.TodoState, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	FindByTodoID(context.Context, uuid.UUID, uuid.UUID) (*domain.Todo, error)
	FindByUserID(context.Context, uuid.UUID) ([]domain.Todo, error)
	List(context.Context, uuid.UUID, ports.TodoListFilter) (*ports.TodoPage, error)
	InitTodoState(context.Context, uuid.UUID) ([]domain.TodoState, error)
}

type todoServiceImpl struct {
	todoRepo ports.TodoRepository
}

func NewTodoService(todoRepo ports.TodoRepository) TodoService {
	return &todoServiceImpl{todoRepo: todoRepo}
}

func (s *todoServiceImpl) Create(ctx context.Context, userID uuid.UUID, dto *dtos.CreateTodoDto) (*domain.Todo, error) {
	if _, err := s.todoRepo.FindStateByID(ctx, userID, dto.TodoStateID); err != nil {
		return nil, mapTodoError(err)
	}

	todo, err := s.todoRepo.Save(ctx, userID, &domain.Todo{
		Title:       dto.Title,
		Description: dto.Description,
		StateID:     dto.TodoStateID,
	})
	if err != nil {
		return nil, errs.NewInternalError("could not create todo")
	}

	return todo, nil
}

func (s *todoServiceImpl) Update(ctx context.Context, userID, todoID uuid.UUID, dto *dtos.UpdateTodoDto) (*domain.Todo, error) {
	fields := make(map[string]any)
	if dto.Title != nil {
		fields["title"] = *dto.Title
	}
	if dto.Description != nil {
		fields["description"] = *dto.Description
	}
	if dto.Priority != nil {
		fields["priority"] = *dto.Priority
	}
	if dto.TodoStateID != nil {
		stateID, err := uuid.Parse(*dto.TodoStateID)
		if err != nil {
			return nil, errs.NewBadRequestError("invalid todo state id")
		}
		if _, err := s.todoRepo.FindStateByID(ctx, userID, stateID); err != nil {
			return nil, mapTodoError(err)
		}
		fields["state_id"] = stateID
	}
	if dto.DueDate != nil {
		if *dto.DueDate == "" {
			fields["due_date"] = nil
		} else {
			dueDate, err := utils.ParseTime(*dto.DueDate)
			if err != nil {
				return nil, errs.NewBadRequestError("invalid due date")
			}
			fields["due_date"] = dueDate
		}
	}
	if len(fields) == 0 {
		return nil, errs.NewBadRequestError("at least one field is required")
	}

	todo, err := s.todoRepo.Update(ctx, userID, todoID, fields)
	if err != nil {
		return nil, mapTodoError(err)
	}

	return todo, nil
}

func (s *todoServiceImpl) UpdateTodoState(ctx context.Context, userID, stateID uuid.UUID, dto *dtos.UpdateTodoState) (*domain.TodoState, error) {
	state, err := s.todoRepo.UpdateTodoState(ctx, userID, stateID, map[string]any{"name": dto.Name})
	if err != nil {
		return nil, mapTodoError(err)
	}

	return state, nil
}

func (s *todoServiceImpl) Delete(ctx context.Context, userID, todoID uuid.UUID) error {
	return mapTodoError(s.todoRepo.Delete(ctx, userID, todoID))
}

func (s *todoServiceImpl) FindByTodoID(ctx context.Context, userID, todoID uuid.UUID) (*domain.Todo, error) {
	todo, err := s.todoRepo.FindByID(ctx, userID, todoID)
	if err != nil {
		return nil, mapTodoError(err)
	}
	return todo, nil
}

func (s *todoServiceImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Todo, error) {
	todos, err := s.todoRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errs.NewInternalError("could not list todos")
	}
	return todos, nil
}

func (s *todoServiceImpl) List(ctx context.Context, userID uuid.UUID, filter ports.TodoListFilter) (*ports.TodoPage, error) {
	if filter.Page < 1 || filter.Limit < 1 || filter.Limit > 100 {
		return nil, errs.NewBadRequestError("invalid pagination")
	}
	if !validTodoSort(filter.Sort) || (filter.Order != "asc" && filter.Order != "desc") {
		return nil, errs.NewBadRequestError("invalid sorting")
	}
	if filter.DueFrom != nil && filter.DueTo != nil && filter.DueFrom.After(*filter.DueTo) {
		return nil, errs.NewBadRequestError("invalid due date range")
	}

	page, err := s.todoRepo.List(ctx, userID, filter)
	if err != nil {
		return nil, errs.NewInternalError("could not list todos")
	}
	return page, nil
}

func validTodoSort(sort string) bool {
	switch sort {
	case "created_at", "updated_at", "due_date", "priority", "title":
		return true
	default:
		return false
	}
}

func (s *todoServiceImpl) InitTodoState(ctx context.Context, userID uuid.UUID) ([]domain.TodoState, error) {
	states, err := s.todoRepo.InitTodoState(ctx, userID)
	if err != nil {
		return nil, errs.NewInternalError("could not initialize todo states")
	}
	return states, nil
}

func mapTodoError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, domain.ErrTodoNotFound) {
		return errs.NewNotFoundError("todo not found")
	}
	if errors.Is(err, domain.ErrTodoStateNotFound) {
		return errs.NewNotFoundError("todo state not found")
	}

	return errs.NewInternalError("todo operation failed")
}
