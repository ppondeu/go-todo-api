package main

import (
	"github.com/ppondeu/go-todo-api/config"
	"github.com/ppondeu/go-todo-api/db"
	"github.com/ppondeu/go-todo-api/internal/domain"
)

func main() {
	cfg := config.LoadConfig()

	dbInstance := db.NewDatabase(cfg)
	db := dbInstance.GetInstance()

	// err = db.Migrator().DropTable(&domain.User{}, &domain.UserSession{}, &domain.Todo{})
	// Migrate the schema
	if err := db.AutoMigrate(&domain.User{}, &domain.UserSession{}, &domain.Todo{}, &domain.TodoState{}); err != nil {
		panic(err)
	}
}
