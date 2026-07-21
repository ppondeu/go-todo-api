package middlewares

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/internal/usecases"
	"github.com/ppondeu/go-todo-api/pkg/errs"
	"github.com/ppondeu/go-todo-api/pkg/response"
)

func JWTRefreshMiddleware(secret []byte, userService usecases.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("refresh_token")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing refresh token in cookie")
			}

			tokenString := cookie.Value
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errs.NewBadRequestError("unexpected signing method")
				}
				return secret, nil
			})
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired refresh token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
			}
			userIDStr, ok := claims["sub"].(string)
			if !ok || userIDStr == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token subject")
			}
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return response.NewErrorResponse(c, errs.NewBadRequestError("Failed to convert string to uuid"))
			}

			user, err := userService.FindByUserID(c.Request().Context(), userID)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
			}
			if user.RefreshToken == nil || *user.RefreshToken != tokenString {
				return echo.NewHTTPError(http.StatusUnauthorized, "Refresh token has been revoked")
			}

			c.Set("user", user)

			return next(c)
		}
	}
}
