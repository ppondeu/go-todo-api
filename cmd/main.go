package main

import (
	"log"

	"github.com/ppondeu/go-todo-api/config"
	"github.com/ppondeu/go-todo-api/db"
	"github.com/ppondeu/go-todo-api/internal/server"
	v "github.com/ppondeu/go-todo-api/pkg/validator"
)

func main() {
	cfg := config.LoadConfig()

	validator := v.NewValidator()
	db := db.NewDatabase(cfg)

	server := server.NewServer(cfg, db.GetInstance())
	server.RegisterRoute(validator)

	if err := server.Start(cfg.Http.Port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
