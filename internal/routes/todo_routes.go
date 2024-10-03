package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/internal/handler"
)

func RegisterTodoRoute(router *echo.Echo, handler *handler.TodoHandler) {
	router.POST("/todos/:userId", handler.CreateTodo)
	router.GET("/todos/:userId", handler.GetTodosByUser)
	router.DELETE("/todos/:id", handler.DeleteTodo)
	router.PATCH("/todos/:id", handler.UpdateTodo)
	router.POST("/todos/:userId/category", handler.CreateCategory)
	router.GET("/todos/:userId/category", handler.GetCategoriesByUser)
	router.PUT("/todos/category/:categoryId", handler.UpdateCategory)
	router.DELETE("/todos/category/:categoryId", handler.DeleteCategory)
	router.GET("/todos/category", handler.GetTodosByCategory)
}
