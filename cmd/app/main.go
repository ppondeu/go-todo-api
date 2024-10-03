package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/config"
	"github.com/ppondeu/go-todo-api/database"
	"github.com/ppondeu/go-todo-api/internal/handler"
	"github.com/ppondeu/go-todo-api/internal/repository"
	"github.com/ppondeu/go-todo-api/internal/routes"
	"github.com/ppondeu/go-todo-api/internal/usecase"
	v "github.com/ppondeu/go-todo-api/pkg/validator"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	validator := v.NewValidator()
	databaseInstance := database.NewPostgresDatabase(cfg)

	userRepo := repository.NewUserRepository(databaseInstance.GetDb())
	userService := usecase.NewUserService(&userRepo)
	userHandler := handler.NewUserHandler(userService, validator)

	todoRepo := repository.NewTodoRepository(databaseInstance.GetDb())
	todoService := usecase.NewTodoService(&todoRepo)
	todoHandler := handler.NewTodoHandler(&todoService, validator)

	router := echo.New()
	defer router.Close()
	router.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Hello, World!"})
	})

	routes.RegisterUserRoute(router, userHandler)
	routes.RegisterTodoRoute(router, todoHandler)

	router.Logger.Fatal(router.Start(fmt.Sprintf(":%d", cfg.Server.Port)))

	
}
