package usecases

import (
	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/internal/repositories"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/logs"
)

type TodoService interface {
	Create(userID uuid.UUID, newTodo *dtos.CreateTodoDto) (*domain.Todo, error)
	Update(ID uuid.UUID, updatedTodo *dtos.UpdateTodoDto) (*domain.Todo, error)
	Delete(ID uuid.UUID) error
	FindByTodoID(ID uuid.UUID) (*domain.Todo, error)
	FindByUserID(userID uuid.UUID) ([]domain.Todo, error)
	FindAll() ([]domain.Todo, error)
}

type todoServiceImpl struct {
	todoRepo repositories.TodoRepository
}

func NewTodoService(todoRepo *repositories.TodoRepository) TodoService {
	return &todoServiceImpl{
		todoRepo: *todoRepo,
	}
}

func (s *todoServiceImpl) Create(userID uuid.UUID, todoDto *dtos.CreateTodoDto) (*domain.Todo, error) {
	newTodo := &domain.Todo{
		Title:       todoDto.Title,
		Description: todoDto.Description,
		UserID:      userID,
	}

	todo, err := s.todoRepo.Save(newTodo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoServiceImpl) Update(ID uuid.UUID, todoUpdateDTO *dtos.UpdateTodoDto) (*domain.Todo, error) {
	todoUpdate := map[string]interface{}{}

	if todoUpdateDTO.Title != nil {
		todoUpdate["title"] = *todoUpdateDTO.Title
	}

	if todoUpdateDTO.Description != nil {
		todoUpdate["description"] = *todoUpdateDTO.Description
	}

	if todoUpdateDTO.Priority != nil {
		todoUpdate["priority"] = *todoUpdateDTO.Priority
	}

	if todoUpdateDTO.TodoStateID != nil {
		todoUpdate["todo_state_id"] = *todoUpdateDTO.TodoStateID
	}

	if todoUpdateDTO.DueDate != nil {
		if *todoUpdateDTO.DueDate == "" {
			todoUpdate["due_date"] = nil
		} else {
			todoUpdate["due_date"] = *todoUpdateDTO.DueDate
		}
	}

	logs.Info(todoUpdate)

	todo, err := s.todoRepo.Update(ID, todoUpdate)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoServiceImpl) Delete(ID uuid.UUID) error {
	err := s.todoRepo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *todoServiceImpl) FindByTodoID(ID uuid.UUID) (*domain.Todo, error) {
	todo, err := s.todoRepo.Find(ID)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoServiceImpl) FindByUserID(userID uuid.UUID) ([]domain.Todo, error) {
	todos, err := s.todoRepo.FindByUserID(userID)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewBadRequestError("user not found")
	}

	return todos, nil
}

func (s *todoServiceImpl) FindAll() ([]domain.Todo, error) {
	todos, err := s.todoRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return todos, nil
}
