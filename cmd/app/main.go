package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/config"
	// "github.com/ppondeu/go-todo-api/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	_ = cfg

	// db := database.NewPostgresDatabase(cfg)

	// _ = db

	router := echo.New()
	router.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Hello, World!"})
	})

	router.Logger.Fatal(router.Start(fmt.Sprintf(":%d", cfg.Server.Port)))

}
