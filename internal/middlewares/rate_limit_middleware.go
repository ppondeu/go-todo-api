package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ppondeu/go-todo-api/config"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/response"
	"golang.org/x/time/rate"
)

func NewSessionRateLimiter(cfg config.RateLimitConfig) echo.MiddlewareFunc {
	cfg = cfg.WithDefaults()
	store := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      rate.Limit(cfg.RequestsPerSecond),
		Burst:     cfg.Burst,
		ExpiresIn: time.Duration(cfg.ExpiresInMinutes) * time.Minute,
	})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*domain.User)
			if !ok {
				return response.NewErrorResponse(c, errs.NewUnauthorizedError("Unauthorized"))
			}

			allowed, err := store.Allow(user.ID.String())
			if err != nil {
				return response.NewErrorResponse(c, errs.NewInternalError("rate limiter unavailable"))
			}
			if !allowed {
				return response.NewAPIResponse(c, http.StatusTooManyRequests, "rate limit exceeded", nil)
			}

			return next(c)
		}
	}
}
