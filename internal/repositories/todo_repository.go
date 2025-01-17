package repositories

import (
	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/pkg/logs"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Save(newTodo *domain.Todo) (*domain.Todo, error)
	Find(where interface{}) (*domain.Todo, error)
	FindAll() ([]domain.Todo, error)
	FindByUserID(userID uuid.UUID) ([]domain.Todo, error)
	Update(ID uuid.UUID, todo map[string]interface{}) (*domain.Todo, error)
	Delete(ID uuid.UUID) error
	UpdateTodoState(ID uuid.UUID, todoStateUpdateField map[string]interface{}) (*domain.TodoState, error)
	InitTodoState(userID uuid.UUID) ([]domain.TodoState, error)
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
	err := r.db.Create(newTodo).Error
	if err != nil {
		return nil, err
	}

	var todo domain.Todo
	err = r.db.Where("id = ?", newTodo.ID).Preload("State").First(&todo).Error
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *TodoRepositoryImpl) Find(where interface{}) (*domain.Todo, error) {
	todo := &domain.Todo{}
	err := r.db.Where(where).Preload("State").First(todo).Error
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (r *TodoRepositoryImpl) FindAll() ([]domain.Todo, error) {
	var todos []domain.Todo
	err := r.db.Preload("State").Find(&todos).Error
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepositoryImpl) FindByUserID(userID uuid.UUID) ([]domain.Todo, error) {
	var todos []domain.Todo
	err := r.db.Where("user_id = ?", userID).Preload("State").Find(&todos).Error
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
	err = r.db.Where("id = ?", ID).Preload("State").First(&updatedTodo).Error
	if err != nil {
		return nil, err
	}

	return &updatedTodo, nil
}

func (r *TodoRepositoryImpl) UpdateTodoState(ID uuid.UUID, todoStateUpdateField map[string]interface{}) (*domain.TodoState, error) {
	err := r.db.Model(&domain.TodoState{}).Where("id = ?", ID).Updates(todoStateUpdateField).Error
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	var todoState domain.TodoState
	err = r.db.Where("id = ?", ID).First(&todoState).Error
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	return &todoState, nil
}

func (r *TodoRepositoryImpl) InitTodoState(userID uuid.UUID) ([]domain.TodoState, error) {
	states := []string{domain.Backlock, domain.NotStarted, domain.InProgress, domain.Done}

	var todoStates []domain.TodoState

	for _, state := range states {
		todoStates = append(todoStates, domain.TodoState{
			Name:   state,
			UserID: userID,
		})
	}

	err := r.db.Create(todoStates).Error
	if err != nil {
		return nil, err
	}

	return todoStates, nil
}

func (r *TodoRepositoryImpl) Delete(ID uuid.UUID) error {
	err := r.db.Where("id = ?", ID).Delete(&domain.Todo{}).Error
	if err != nil {
		return err
	}

	return nil
}
