package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/internal/application/ports"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
)

type fakeTodoService struct {
	listPage   *ports.TodoPage
	gotUserID  uuid.UUID
	gotFilter  ports.TodoListFilter
	gotContext context.Context
}

func (f *fakeTodoService) Create(context.Context, uuid.UUID, *dtos.CreateTodoDto) (*domain.Todo, error) {
	return &domain.Todo{}, nil
}
func (f *fakeTodoService) Update(context.Context, uuid.UUID, uuid.UUID, *dtos.UpdateTodoDto) (*domain.Todo, error) {
	return &domain.Todo{}, nil
}
func (f *fakeTodoService) UpdateTodoState(context.Context, uuid.UUID, uuid.UUID, *dtos.UpdateTodoState) (*domain.TodoState, error) {
	return &domain.TodoState{}, nil
}
func (f *fakeTodoService) Delete(context.Context, uuid.UUID, uuid.UUID) error { return nil }
func (f *fakeTodoService) FindByTodoID(context.Context, uuid.UUID, uuid.UUID) (*domain.Todo, error) {
	return &domain.Todo{}, nil
}
func (f *fakeTodoService) FindByUserID(context.Context, uuid.UUID) ([]domain.Todo, error) {
	return []domain.Todo{}, nil
}
func (f *fakeTodoService) List(ctx context.Context, userID uuid.UUID, filter ports.TodoListFilter) (*ports.TodoPage, error) {
	f.gotContext, f.gotUserID, f.gotFilter = ctx, userID, filter
	return f.listPage, nil
}
func (f *fakeTodoService) InitTodoState(uuid.UUID) ([]domain.TodoState, error) {
	return []domain.TodoState{}, nil
}

func TestTodoHandlerListTodos(t *testing.T) {
	userID := uuid.New()
	service := &fakeTodoService{listPage: &ports.TodoPage{
		Items: []domain.Todo{}, Page: 2, Limit: 10, Total: 12, TotalPages: 2,
	}}
	handler := NewTodoHandler(service, validatorv10.New())
	e := echo.New()
	type contextKey string
	key := contextKey("test")
	request := httptest.NewRequest(http.MethodGet, "/todos?page=2&limit=10&priority=high&sort=title&order=asc", nil)
	request = request.WithContext(context.WithValue(request.Context(), key, "request-context"))
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	c.Set("user", &domain.User{ID: userID})

	if err := handler.ListTodos(c); err != nil {
		t.Fatalf("ListTodos() error = %v", err)
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("ListTodos() status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if service.gotUserID != userID || service.gotFilter.Page != 2 || service.gotFilter.Limit != 10 {
		t.Fatalf("ListTodos() did not pass authenticated user and pagination")
	}
	if service.gotContext.Value(key) != "request-context" {
		t.Fatal("ListTodos() did not pass request context")
	}

	var body struct {
		Data struct {
			Items      []domain.Todo `json:"items"`
			Pagination struct {
				TotalPages int `json:"totalPages"`
			} `json:"pagination"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Items == nil || body.Data.Pagination.TotalPages != 2 {
		t.Fatalf("ListTodos() returned invalid pagination response")
	}
}

func TestTodoHandlerListTodosRejectsInvalidQuery(t *testing.T) {
	tests := []string{
		"?page=0",
		"?limit=101",
		"?priority=urgent",
		"?state_id=invalid",
		"?due_from=invalid",
		"?due_from=2026-08-01T00:00:00Z&due_to=2026-07-01T00:00:00Z",
		"?sort=random",
		"?order=random",
	}

	for _, query := range tests {
		t.Run(query, func(t *testing.T) {
			service := &fakeTodoService{}
			handler := NewTodoHandler(service, validatorv10.New())
			e := echo.New()
			recorder := httptest.NewRecorder()
			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/todos"+query, nil), recorder)
			c.Set("user", &domain.User{ID: uuid.New()})

			if err := handler.ListTodos(c); err != nil {
				t.Fatalf("ListTodos() error = %v", err)
			}
			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("ListTodos() status = %d, want %d", recorder.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestTodoHandlerListTodosRequiresAuthentication(t *testing.T) {
	service := &fakeTodoService{}
	handler := NewTodoHandler(service, validatorv10.New())
	e := echo.New()
	recorder := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/todos", nil), recorder)

	if err := handler.ListTodos(c); err != nil {
		t.Fatalf("ListTodos() error = %v", err)
	}
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("ListTodos() status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestTodoHandlerLegacyListRejectsAnotherUser(t *testing.T) {
	service := &fakeTodoService{}
	handler := NewTodoHandler(service, validatorv10.New())
	e := echo.New()
	recorder := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/todos/other", nil), recorder)
	c.SetPath("/todos/:userId")
	c.SetParamNames("userId")
	c.SetParamValues(uuid.NewString())
	c.Set("user", &domain.User{ID: uuid.New()})

	if err := handler.GetTodosByUser(c); err != nil {
		t.Fatalf("GetTodosByUser() error = %v", err)
	}
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("GetTodosByUser() status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}
