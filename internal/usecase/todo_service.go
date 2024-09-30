package usecase

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/internal/repository"
	"github.com/ppondeu/go-todo-api/pkg/dto"
	"github.com/ppondeu/go-todo-api/pkg/utils"
)

type TodoService interface {
	Create(userID uuid.UUID, newTodo *dto.CreateTodoDto) (*domain.Todo, error)
	Update(ID uuid.UUID, updatedTodo *dto.UpdateTodoDto) (*domain.Todo, error)
	Delete(ID uuid.UUID) error
	FindByTodoID(ID uuid.UUID) (*domain.Todo, error)
	FindAll() ([]domain.Todo, error)
	CreateCategory(userID uuid.UUID, name string) (*domain.TodoCategory, error)
	FindCategories(userID uuid.UUID) ([]domain.TodoCategory, error)
	UpdateCategory(categoryID uuid.UUID, name string) (*domain.TodoCategory, error)
	DeleteCategory(categoryID uuid.UUID) error
	FindTodosByCategory(categoryID uuid.UUID) ([]domain.Todo, error)
	UpdateTodoCategory(todoID uuid.UUID, categoryID uuid.UUID) (*domain.Todo, error)
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

	updateTodo := &domain.Todo{
		Title:       *todoDto.Title,
		Description: *todoDto.Description,
		DueDate:     utils.ParseTime(*todoDto.DueDate),
		Priority:    domain.Priority(*todoDto.Priority),
		State:       domain.TodoState(*todoDto.State),
		CategoryID:  utils.ParseUUID(*todoDto.CategoryID),
	}

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

func (s *todoServiceImpl) UpdateTodoCategory(todoID uuid.UUID, categoryID uuid.UUID) (*domain.Todo, error) {
	category, err := s.todoRepo.FindCategory(categoryID)
	if err != nil {
		return nil, err
	}

	todo, err := s.todoRepo.Find(todoID)
	if err != nil {
		return nil, err
	}

	if category.UserID != todo.UserID {
		return nil, errors.New("category not found")
	}

	updateTodo := &domain.Todo{
		CategoryID: &categoryID,
	}

	updatedTodo, err := s.todoRepo.Update(todoID, updateTodo)
	if err != nil {
		return nil, err
	}

	return updatedTodo, nil
}
