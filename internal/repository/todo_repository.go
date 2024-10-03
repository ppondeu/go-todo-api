package repository

import (
	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Save(newTodo *domain.Todo) (*domain.Todo, error)
	Find(where interface{}) (*domain.Todo, error)
	FindAll() ([]domain.Todo, error)
	FindByUserID(userID uuid.UUID) ([]domain.Todo, error)
	Update(ID uuid.UUID, todo map[string]interface{}) (*domain.Todo, error)
	Delete(ID uuid.UUID) error
	SaveCategory(newCategory *domain.TodoCategory) (*domain.TodoCategory, error)
	FindCategories(userID uuid.UUID) ([]domain.TodoCategory, error)
	FindTodosByCategory(categoryID uuid.UUID) ([]domain.Todo, error)
	UpdateCategory(categoryID uuid.UUID, category *domain.TodoCategory) (*domain.TodoCategory, error)
	DeleteCategory(categoryID uuid.UUID) error
	FindCategory(categoryID uuid.UUID) (*domain.TodoCategory, error)
}

type TodoRepositoryImpl struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &TodoRepositoryImpl{
		db: db,
	}
}

func (r *TodoRepositoryImpl) Save(newTodo *domain.Todo) (*domain.Todo, error) {
	err := r.db.Omit("Category").Create(newTodo).Error
	if err != nil {
		return nil, err
	}

	var todo domain.Todo
	err = r.db.Where("id = ?", newTodo.ID).Preload("Category").First(&todo).Error
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepositoryImpl) Find(where interface{}) (*domain.Todo, error) {
	todo := &domain.Todo{}
	err := r.db.Where(where).Preload("Category").First(todo).Error
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (r *TodoRepositoryImpl) FindAll() ([]domain.Todo, error) {
	var todos []domain.Todo
	err := r.db.Preload("Category").Find(&todos).Error
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepositoryImpl) FindByUserID(userID uuid.UUID) ([]domain.Todo, error) {
	var todos []domain.Todo
	err := r.db.Where("user_id = ?", userID).Preload("Category").Find(&todos).Error
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepositoryImpl) Update(ID uuid.UUID, todo map[string]interface{}) (*domain.Todo, error) {
	err := r.db.Model(&domain.Todo{}).Where("id = ?", ID).Updates(todo).Error

	var updatedTodo domain.Todo
	err = r.db.Where("id = ?", ID).Preload("Category").First(&updatedTodo).Error
	if err != nil {
		return nil, err
	}

	return &updatedTodo, nil
}

func (r *TodoRepositoryImpl) Delete(ID uuid.UUID) error {
	err := r.db.Where("id = ?", ID).Delete(&domain.Todo{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *TodoRepositoryImpl) SaveCategory(newCategory *domain.TodoCategory) (*domain.TodoCategory, error) {
	err := r.db.Create(newCategory).Error
	if err != nil {
		return nil, err
	}

	return newCategory, nil
}

func (r *TodoRepositoryImpl) FindCategories(userID uuid.UUID) ([]domain.TodoCategory, error) {
	var categories []domain.TodoCategory
	err := r.db.Where("user_id = ?", userID).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *TodoRepositoryImpl) FindTodosByCategory(categoryID uuid.UUID) ([]domain.Todo, error) {
	var todos []domain.Todo
	err := r.db.Where("category_id = ?", categoryID).Preload("Category").Find(&todos).Error
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepositoryImpl) UpdateCategory(categoryID uuid.UUID, category *domain.TodoCategory) (*domain.TodoCategory, error) {
	err := r.db.Model(category).Where("id = ?", categoryID).Updates(category).Error
	if err != nil {
		return nil, err
	}

	var updatedCategory domain.TodoCategory
	err = r.db.Where("id = ?", categoryID).First(&updatedCategory).Error
	if err != nil {
		return nil, err
	}

	return &updatedCategory, nil
}

func (r *TodoRepositoryImpl) DeleteCategory(categoryID uuid.UUID) error {
	err := r.db.Where("id = ?", categoryID).Delete(&domain.TodoCategory{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *TodoRepositoryImpl) FindCategory(categoryID uuid.UUID) (*domain.TodoCategory, error) {
	category := &domain.TodoCategory{}
	err := r.db.Where("id = ?", categoryID).First(category).Error
	if err != nil {
		return nil, err
	}

	return category, nil
}
