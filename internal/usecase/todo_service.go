package usecase

import (
	// "fmt"
	// "time"

	"time"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/internal/repository"
	"github.com/ppondeu/go-todo-api/pkg/dto"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/logs"
	"github.com/ppondeu/go-todo-api/pkg/utils"
	// "github.com/ppondeu/go-todo-api/pkg/utils"
)

type TodoService interface {
	Create(userID uuid.UUID, newTodo *dto.CreateTodoDto) (*domain.Todo, error)
	Update(ID uuid.UUID, updatedTodo *dto.UpdateTodoDto) (*domain.Todo, error)
	Delete(ID uuid.UUID) error
	FindByTodoID(ID uuid.UUID) (*domain.Todo, error)
	FindByUserID(userID uuid.UUID) ([]domain.Todo, error)
	FindAll() ([]domain.Todo, error)
	CreateCategory(userID uuid.UUID, name string) (*domain.TodoCategory, error)
	FindCategories(userID uuid.UUID) ([]domain.TodoCategory, error)
	UpdateCategory(categoryID uuid.UUID, name string) (*domain.TodoCategory, error)
	DeleteCategory(categoryID uuid.UUID) error
	FindTodosByCategory(categoryID uuid.UUID) ([]domain.Todo, error)
}

type todoServiceImpl struct {
	todoRepo repository.TodoRepository
}

func NewTodoService(todoRepo *repository.TodoRepository) TodoService {
	return &todoServiceImpl{
		todoRepo: *todoRepo,
	}
}

func (s *todoServiceImpl) Create(userID uuid.UUID, todoDto *dto.CreateTodoDto) (*domain.Todo, error) {
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

func (s *todoServiceImpl) Update(ID uuid.UUID, todoDto *dto.UpdateTodoDto) (*domain.Todo, error) {
	updateTodo := map[string]interface{}{}

	if todoDto.Title != nil {
		updateTodo["title"] = *todoDto.Title
	}

	if todoDto.Description != nil {
		updateTodo["description"] = *todoDto.Description
	}

	if todoDto.Priority != nil {
		updateTodo["priority"] = *todoDto.Priority
	}

	if todoDto.State != nil {
		updateTodo["state"] = *todoDto.State
		if *todoDto.State == "done" {
			updateTodo["is_completed"] = true
		} else {
			updateTodo["is_completed"] = false
		}
	}

	if todoDto.CategoryID != nil {
		if *todoDto.CategoryID == "" {
			updateTodo["category_id"] = nil
		} else {
			updateTodo["category_id"] = *todoDto.CategoryID
		}
	}

	if todoDto.DueDate != nil {
		if *todoDto.DueDate == "" {
			updateTodo["due_date"] = nil
			updateTodo["is_overdue"] = false
		} else {
			updateTodo["due_date"] = *todoDto.DueDate
			dueDate, err := utils.ParseTime(*todoDto.DueDate)
			if err != nil {
				return nil, err
			}

			if dueDate.Before(time.Now()) {
				updateTodo["is_overdue"] = true
			} else {
				updateTodo["is_overdue"] = false
			}
		}
	}

	logs.Info(updateTodo)

	todo, err := s.todoRepo.Update(ID, updateTodo)
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

func (s *todoServiceImpl) CreateCategory(userID uuid.UUID, name string) (*domain.TodoCategory, error) {
	newCategory := &domain.TodoCategory{
		Name:   name,
		UserID: userID,
	}

	category, err := s.todoRepo.SaveCategory(newCategory)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *todoServiceImpl) FindCategories(userID uuid.UUID) ([]domain.TodoCategory, error) {
	categories, err := s.todoRepo.FindCategories(userID)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *todoServiceImpl) UpdateCategory(categoryID uuid.UUID, name string) (*domain.TodoCategory, error) {
	updateCategory := &domain.TodoCategory{
		Name: name,
	}

	category, err := s.todoRepo.UpdateCategory(categoryID, updateCategory)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *todoServiceImpl) DeleteCategory(categoryID uuid.UUID) error {
	err := s.todoRepo.DeleteCategory(categoryID)
	if err != nil {
		return err
	}

	return nil
}

func (s *todoServiceImpl) FindTodosByCategory(categoryID uuid.UUID) ([]domain.Todo, error) {
	todos, err := s.todoRepo.FindTodosByCategory(categoryID)
	if err != nil {
		return nil, err
	}

	return todos, nil
}
