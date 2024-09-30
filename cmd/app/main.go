package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ppondeu/go-todo-api/config"
	"github.com/ppondeu/go-todo-api/database"
	"github.com/ppondeu/go-todo-api/internal/domain"
	"github.com/ppondeu/go-todo-api/internal/repository"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	_ = cfg

	databaseInstance := database.NewPostgresDatabase(cfg)
	categoryID, _ := uuid.Parse("d0fe7c1b-4765-49e8-b5b2-3677215aef23")
	userID, err := uuid.Parse("a2723d86-6066-4aad-a04a-4758f040e855")
	todoID, _ := uuid.Parse("a21c7bac-1653-45fc-b952-3766821a2979")
	todoRepo := repository.NewTodoRepository(databaseInstance.GetDb())

	updateTodo := &domain.Todo{
		CategoryID: &categoryID,
	}

	todo, err := todoRepo.Update(todoID, updateTodo)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", todo.Title)
	if todo.Category != nil {
		fmt.Printf("%v\n", todo.Category.Name)
	} else {
		fmt.Println("No category")
	}

	todos, _ := todoRepo.FindByUserID(userID)

	for _, todo := range todos {
		fmt.Printf("%v\n", todo.Title)
		if todo.Category != nil {
			fmt.Printf("%v\n", todo.Category.Name)
		} else {
			fmt.Println("No category")
		}
	}

	// router := echo.New()
	// router.GET("/", func(c echo.Context) error {
	// 	return c.JSON(200, map[string]string{"message": "Hello, World!"})
	// })

	// router.Logger.Fatal(router.Start(fmt.Sprintf(":%d", cfg.Server.Port)))

}
