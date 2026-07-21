package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/pkg/dtos"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/utils"
)

type fakeUserRepository struct {
	createErr      error
	findByIDErr    error
	updateErr      error
	deleteErr      error
	gotContext     context.Context
	gotUserID      uuid.UUID
	gotFields      map[string]any
	gotCreatedUser *domain.User
}

func (f *fakeUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	f.gotContext, f.gotCreatedUser = ctx, user
	if f.createErr != nil {
		return nil, f.createErr
	}
	user.ID = uuid.New()
	return user, nil
}
func (f *fakeUserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	f.gotContext, f.gotUserID = ctx, userID
	if f.findByIDErr != nil {
		return nil, f.findByIDErr
	}
	return &domain.User{ID: userID, Email: "user@example.com"}, nil
}
func (f *fakeUserRepository) FindByEmail(context.Context, string) (*domain.User, error) {
	return &domain.User{}, nil
}
func (f *fakeUserRepository) Update(ctx context.Context, userID uuid.UUID, fields map[string]any) (*domain.User, error) {
	f.gotContext, f.gotUserID, f.gotFields = ctx, userID, fields
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	return &domain.User{ID: userID}, nil
}
func (f *fakeUserRepository) UpdateRefreshToken(context.Context, uuid.UUID, *string) error {
	return nil
}
func (f *fakeUserRepository) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	f.gotContext, f.gotUserID = ctx, userID
	return f.deleteErr
}

func TestUserServiceUpdateScopesToAuthenticatedUserAndHashesPassword(t *testing.T) {
	repo := &fakeUserRepository{}
	service := NewUserService(repo)
	userID := uuid.New()
	password := "new-password"
	ctx := context.WithValue(context.Background(), struct{}{}, "request")

	_, err := service.Update(ctx, userID, &dtos.UserUpdateDTO{Password: &password})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if repo.gotContext != ctx || repo.gotUserID != userID {
		t.Fatal("Update() did not preserve request context and authenticated user ID")
	}
	hash, ok := repo.gotFields["password"].(string)
	if !ok || hash == password || utils.ComparePassword(hash, password) != nil {
		t.Fatal("Update() did not persist a valid password hash")
	}
	if password != "new-password" {
		t.Fatal("Update() mutated the request DTO")
	}
}

func TestUserServiceUpdateRejectsEmptyInput(t *testing.T) {
	service := NewUserService(&fakeUserRepository{})

	_, err := service.Update(context.Background(), uuid.New(), &dtos.UserUpdateDTO{})
	var appErr *errs.AppError
	if !errors.As(err, &appErr) || appErr.Code != 400 {
		t.Fatalf("Update() error = %v, want 400 app error", err)
	}
}

func TestUserServiceMapsRepositoryErrors(t *testing.T) {
	tests := []struct {
		name     string
		repoErr  error
		wantCode int
	}{
		{name: "not found", repoErr: domain.ErrUserNotFound, wantCode: 404},
		{name: "email conflict", repoErr: domain.ErrEmailConflict, wantCode: 409},
		{name: "database failure", repoErr: errors.New("sensitive database detail"), wantCode: 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewUserService(&fakeUserRepository{createErr: tt.repoErr})
			_, err := service.Save(context.Background(), &dtos.UserCreateDTO{Email: "user@example.com", Password: "password"})
			var appErr *errs.AppError
			if !errors.As(err, &appErr) || appErr.Code != tt.wantCode {
				t.Fatalf("Save() error = %v, want status %d", err, tt.wantCode)
			}
			if appErr.Message == "sensitive database detail" {
				t.Fatal("Save() leaked repository error details")
			}
		})
	}
}

func TestUserServiceDeleteUsesAuthenticatedUser(t *testing.T) {
	repo := &fakeUserRepository{}
	service := NewUserService(repo)
	userID := uuid.New()

	if err := service.Delete(context.Background(), userID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if repo.gotUserID != userID {
		t.Fatalf("Delete() user ID = %s, want %s", repo.gotUserID, userID)
	}
}
