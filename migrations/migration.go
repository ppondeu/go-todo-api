package main

import (
	"github.com/ppondeu/go-todo-api/config"
	"github.com/ppondeu/go-todo-api/db"
	"github.com/ppondeu/go-todo-api/internal/domain"
)

var err error

func main() {
	cfg := config.LoadConfig()

	dbInstance := db.NewDatabase(cfg)
	db := dbInstance.GetInstance()

	// err = db.Migrator().DropTable(&domain.User{}, &domain.UserSession{}, &domain.Todo{})
	// Migrate the schema
	err = db.AutoMigrate(&domain.User{}, &domain.UserSession{}, &domain.Todo{})
	if err != nil {
		panic(err)
	}
}
