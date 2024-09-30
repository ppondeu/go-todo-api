package main

import (
	"github.com/ppondeu/go-todo-api/config"
	"github.com/ppondeu/go-todo-api/database"
	"github.com/ppondeu/go-todo-api/internal/domain"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbInstance := database.NewPostgresDatabase(cfg)
	db := dbInstance.GetDb()

	// err = db.Migrator().DropTable(&domain.User{}, &domain.UserSession{}, &domain.TodoCategory{}, &domain.Todo{})
	// Migrate the schema
	err = db.AutoMigrate(&domain.User{}, &domain.UserSession{}, &domain.TodoCategory{}, &domain.Todo{})
	if err != nil {
		panic(err)
	}
}
