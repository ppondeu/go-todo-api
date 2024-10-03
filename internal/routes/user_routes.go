package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/internal/handler"
)

func RegisterUserRoute(router *echo.Echo, handler *handler.UserHandler) {
	router.POST("/users", handler.Register)
	router.GET("/users", handler.GetUsers)
	router.GET("/users/:id", handler.GetUser)
	router.PATCH("/users/:id", handler.UpdateUser)
	router.DELETE("/users/:id", handler.DeleteUser)
}
