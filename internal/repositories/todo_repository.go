package repositories

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
	if err != nil {
		return nil, err
	}
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
