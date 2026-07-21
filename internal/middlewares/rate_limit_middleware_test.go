package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/config"
	"github.com/ppondeu/go-todo-api/internal/domain"
)

func TestSessionRateLimiter(t *testing.T) {
	cfg := config.RateLimitConfig{
		RequestsPerSecond: 0.000001,
		Burst:             2,
		ExpiresInMinutes:  10,
	}

	tests := []struct {
		name       string
		user       *domain.User
		requests   int
		wantStatus []int
	}{
		{
			name:     "rejects request without an authenticated session",
			requests: 1,
			wantStatus: []int{
				http.StatusUnauthorized,
			},
		},
		{
			name:     "limits requests from the same session",
			user:     &domain.User{ID: uuid.New()},
			requests: 3,
			wantStatus: []int{
				http.StatusOK,
				http.StatusOK,
				http.StatusTooManyRequests,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			handler := NewSessionRateLimiter(cfg)(func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			})

			for i := 0; i < tt.requests; i++ {
				recorder := httptest.NewRecorder()
				ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), recorder)
				if tt.user != nil {
					ctx.Set("user", tt.user)
				}

				if err := handler(ctx); err != nil {
					t.Fatalf("request %d returned an error: %v", i+1, err)
				}
				if recorder.Code != tt.wantStatus[i] {
					t.Errorf("request %d status = %d, want %d", i+1, recorder.Code, tt.wantStatus[i])
				}
				if recorder.Code == http.StatusTooManyRequests {
					assertRateLimitResponse(t, recorder)
				}
			}
		})
	}
}

func assertRateLimitResponse(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	var body struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode rate limit response: %v", err)
	}
	if body.StatusCode != http.StatusTooManyRequests {
		t.Errorf("response statusCode = %d, want %d", body.StatusCode, http.StatusTooManyRequests)
	}
	if body.Message != "rate limit exceeded" {
		t.Errorf("response message = %q, want %q", body.Message, "rate limit exceeded")
	}
}

func TestSessionRateLimiterSeparatesSessions(t *testing.T) {
	cfg := config.RateLimitConfig{
		RequestsPerSecond: 0.000001,
		Burst:             1,
		ExpiresInMinutes:  10,
	}
	e := echo.New()
	handler := NewSessionRateLimiter(cfg)(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	for _, userID := range []uuid.UUID{uuid.New(), uuid.New()} {
		recorder := httptest.NewRecorder()
		ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), recorder)
		ctx.Set("user", &domain.User{ID: userID})

		if err := handler(ctx); err != nil {
			t.Fatalf("session %s returned an error: %v", userID, err)
		}
		if recorder.Code != http.StatusOK {
			t.Errorf("session %s status = %d, want %d", userID, recorder.Code, http.StatusOK)
		}
	}
}
