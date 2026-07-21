package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
)

type fakeUserService struct {
	user       *domain.User
	gotContext context.Context
	gotUserID  uuid.UUID
	gotUpdate  *dtos.UserUpdateDTO
	deleteErr  error
}

func (f *fakeUserService) Save(context.Context, *dtos.UserCreateDTO) (*domain.User, error) {
	return f.user, nil
}
func (f *fakeUserService) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	f.gotContext, f.gotUserID = ctx, userID
	return f.user, nil
}
func (f *fakeUserService) FindByEmail(context.Context, string) (*domain.User, error) {
	return f.user, nil
}
func (f *fakeUserService) Update(ctx context.Context, userID uuid.UUID, dto *dtos.UserUpdateDTO) (*domain.User, error) {
	f.gotContext, f.gotUserID, f.gotUpdate = ctx, userID, dto
	return f.user, nil
}
func (f *fakeUserService) UpdateRefreshToken(context.Context, uuid.UUID, *string) error {
	return nil
}
func (f *fakeUserService) Delete(ctx context.Context, userID uuid.UUID) error {
	f.gotContext, f.gotUserID = ctx, userID
	return f.deleteErr
}

func TestUserHandlerGetMeReturnsSafeResponse(t *testing.T) {
	userID := uuid.New()
	service := &fakeUserService{user: &domain.User{
		ID: userID, Email: "user@example.com", FirstName: "Test", Password: "secret-hash", RefreshToken: stringPtr("secret-token"),
	}}
	handler := NewUserHandler(service, validatorv10.New())
	e := echo.New()
	type contextKey string
	key := contextKey("request")
	request := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	request = request.WithContext(context.WithValue(request.Context(), key, "value"))
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	c.Set("user", &domain.User{ID: userID})

	if err := handler.GetMe(c); err != nil {
		t.Fatalf("GetMe() error = %v", err)
	}
	if recorder.Code != http.StatusOK || service.gotUserID != userID || service.gotContext.Value(key) != "value" {
		t.Fatal("GetMe() did not use the authenticated user and request context")
	}
	body := recorder.Body.String()
	for _, sensitiveField := range []string{"password", "refresh_token", "secret-hash", "secret-token", "todos", "todo_states"} {
		if strings.Contains(body, sensitiveField) {
			t.Fatalf("GetMe() response contains sensitive field or value %q", sensitiveField)
		}
	}
}

func TestUserHandlerUpdateMeUsesAuthenticatedUser(t *testing.T) {
	userID := uuid.New()
	service := &fakeUserService{user: &domain.User{ID: userID}}
	handler := NewUserHandler(service, validatorv10.New())
	e := echo.New()
	recorder := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodPatch, "/users/me", bytes.NewBufferString(`{"first_name":"Updated"}`)), recorder)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Set("user", &domain.User{ID: userID})

	if err := handler.UpdateMe(c); err != nil {
		t.Fatalf("UpdateMe() error = %v", err)
	}
	if service.gotUserID != userID || service.gotUpdate == nil || service.gotUpdate.FirstName == nil || *service.gotUpdate.FirstName != "Updated" {
		t.Fatal("UpdateMe() did not use the authenticated user and parsed input")
	}
}

func TestUserHandlerDeleteMeClearsAuthenticationCookies(t *testing.T) {
	userID := uuid.New()
	service := &fakeUserService{}
	handler := NewUserHandler(service, validatorv10.New())
	e := echo.New()
	recorder := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodDelete, "/users/me", nil), recorder)
	c.Set("user", &domain.User{ID: userID})

	if err := handler.DeleteMe(c); err != nil {
		t.Fatalf("DeleteMe() error = %v", err)
	}
	if service.gotUserID != userID {
		t.Fatalf("DeleteMe() user ID = %s, want %s", service.gotUserID, userID)
	}

	cookies := recorder.Result().Cookies()
	if len(cookies) != 2 {
		t.Fatalf("DeleteMe() cookies = %d, want 2", len(cookies))
	}
	for _, cookie := range cookies {
		if cookie.Value != "" || cookie.MaxAge >= 0 {
			t.Fatalf("DeleteMe() did not expire cookie %q", cookie.Name)
		}
	}
}

func TestUserHandlerRequiresAuthentication(t *testing.T) {
	handler := NewUserHandler(&fakeUserService{}, validatorv10.New())
	e := echo.New()
	recorder := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/users/me", nil), recorder)

	if err := handler.GetMe(c); err != nil {
		t.Fatalf("GetMe() error = %v", err)
	}
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("GetMe() status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestUserHandlerResponseShape(t *testing.T) {
	response := toUserResponse(&domain.User{ID: uuid.New(), Email: "user@example.com"})
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}
	if strings.Contains(string(data), "password") || strings.Contains(string(data), "refresh_token") {
		t.Fatal("user response exposes authentication fields")
	}
}

func stringPtr(value string) *string { return &value }
