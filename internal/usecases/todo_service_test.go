package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/application/ports"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
)

type fakeTodoRepository struct {
	findStateErr error
	updateErr    error
	deleteErr    error
	listResult   *ports.TodoPage
	listErr      error
	gotUserID    uuid.UUID
	gotTodoID    uuid.UUID
	gotFilter    ports.TodoListFilter
}

func (f *fakeTodoRepository) Save(_ context.Context, userID uuid.UUID, todo *domain.Todo) (*domain.Todo, error) {
	f.gotUserID = userID
	todo.UserID = userID
	return todo, nil
}
func (f *fakeTodoRepository) FindByID(context.Context, uuid.UUID, uuid.UUID) (*domain.Todo, error) {
	return nil, domain.ErrTodoNotFound
}
func (f *fakeTodoRepository) FindByUserID(context.Context, uuid.UUID) ([]domain.Todo, error) {
	return []domain.Todo{}, nil
}
func (f *fakeTodoRepository) List(_ context.Context, userID uuid.UUID, filter ports.TodoListFilter) (*ports.TodoPage, error) {
	f.gotUserID, f.gotFilter = userID, filter
	return f.listResult, f.listErr
}
func (f *fakeTodoRepository) Update(_ context.Context, userID, todoID uuid.UUID, _ map[string]any) (*domain.Todo, error) {
	f.gotUserID, f.gotTodoID = userID, todoID
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	return &domain.Todo{ID: todoID, UserID: userID}, nil
}
func (f *fakeTodoRepository) Delete(_ context.Context, userID, todoID uuid.UUID) error {
	f.gotUserID, f.gotTodoID = userID, todoID
	return f.deleteErr
}
func (f *fakeTodoRepository) FindStateByID(context.Context, uuid.UUID, uuid.UUID) (*domain.TodoState, error) {
	if f.findStateErr != nil {
		return nil, f.findStateErr
	}
	return &domain.TodoState{}, nil
}
func (f *fakeTodoRepository) UpdateTodoState(context.Context, uuid.UUID, uuid.UUID, map[string]any) (*domain.TodoState, error) {
	return &domain.TodoState{}, nil
}
func (f *fakeTodoRepository) InitTodoState(uuid.UUID) ([]domain.TodoState, error) {
	return []domain.TodoState{}, nil
}

func TestTodoServiceCreateRejectsStateOutsideUserScope(t *testing.T) {
	repo := &fakeTodoRepository{findStateErr: domain.ErrTodoStateNotFound}
	service := NewTodoService(repo)

	_, err := service.Create(context.Background(), uuid.New(), &dtos.CreateTodoDto{TodoStateID: uuid.New()})
	var appErr *errs.AppError
	if !errors.As(err, &appErr) || appErr.Code != 404 {
		t.Fatalf("Create() error = %v, want 404 app error", err)
	}
}

func TestTodoServiceUpdateScopesOperationToUser(t *testing.T) {
	repo := &fakeTodoRepository{}
	service := NewTodoService(repo)
	userID, todoID := uuid.New(), uuid.New()
	title := "updated title"

	_, err := service.Update(context.Background(), userID, todoID, &dtos.UpdateTodoDto{Title: &title})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if repo.gotUserID != userID || repo.gotTodoID != todoID {
		t.Fatalf("Update() scope = (%s, %s), want (%s, %s)", repo.gotUserID, repo.gotTodoID, userID, todoID)
	}
}

func TestTodoServiceDeleteHidesCrossUserTodo(t *testing.T) {
	repo := &fakeTodoRepository{deleteErr: domain.ErrTodoNotFound}
	service := NewTodoService(repo)

	err := service.Delete(context.Background(), uuid.New(), uuid.New())
	var appErr *errs.AppError
	if !errors.As(err, &appErr) || appErr.Code != 404 {
		t.Fatalf("Delete() error = %v, want 404 app error", err)
	}
}

func TestTodoServiceListPassesAuthenticatedUserAndFilter(t *testing.T) {
	wantPage := &ports.TodoPage{Items: []domain.Todo{}, Page: 2, Limit: 10, Total: 12, TotalPages: 2}
	repo := &fakeTodoRepository{listResult: wantPage}
	service := NewTodoService(repo)
	userID := uuid.New()
	filter := ports.TodoListFilter{Page: 2, Limit: 10, Sort: "title", Order: "asc"}

	page, err := service.List(context.Background(), userID, filter)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if page != wantPage || repo.gotUserID != userID || repo.gotFilter != filter {
		t.Fatalf("List() did not preserve page, user, and filter")
	}
}
