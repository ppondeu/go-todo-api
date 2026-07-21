package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ppondeu/go-todo-api/config"
	"github.com/ppondeu/go-todo-api/internal/handlers"
	"github.com/ppondeu/go-todo-api/internal/middlewares"
	"github.com/ppondeu/go-todo-api/internal/repositories"
	"github.com/ppondeu/go-todo-api/internal/usecases"
	"gorm.io/gorm"
)

type Server struct {
	Echo   *echo.Echo
	DB     *gorm.DB
	Config *config.Config
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	e := echo.New()
	return &Server{
		Echo:   e,
		Config: cfg,
		DB:     db,
	}
}

func (server *Server) Start(addr string) error {
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		<-sigs
		fmt.Println("\nShutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Echo.Shutdown(ctx); err != nil {
			server.Echo.Logger.Error("Error shutting down server:", err)
		}
	}()

	return server.Echo.Start(":" + addr)
}

func (server *Server) RegisterRoute(validator *validator.Validate) {
	userRepo := repositories.NewUserRepository(server.DB)
	userService := usecases.NewUserService(userRepo)

	todoRepo := repositories.NewTodoRepository(server.DB)
	todoService := usecases.NewTodoService(todoRepo)

	userHandler := handlers.NewUserHandler(userService, validator)

	todoHandler := handlers.NewTodoHandler(todoService, validator)

	jwtService := usecases.NewJwtService([]byte(server.Config.Auth.AccessSecret), []byte(server.Config.Auth.RefreshSecret))
	authService := usecases.NewAuthService(userService, todoService, jwtService)
	authHandler := handlers.NewAuthHandler(authService, validator)

	jwtAccessMiddleware := middlewares.JWTAccessMiddleware([]byte(server.Config.Auth.AccessSecret), userService)
	jwtRefreshMiddleware := middlewares.JWTRefreshMiddleware([]byte(server.Config.Auth.RefreshSecret), userService)

	routeGroup := server.Echo.Group("/api/v1")

	userGroup := routeGroup.Group("/users")
	userGroup.Use(jwtAccessMiddleware)
	userGroup.GET("", userHandler.GetUsers)
	userGroup.GET("/:id", userHandler.GetUser)
	userGroup.PATCH("/:id", userHandler.UpdateUser)
	userGroup.DELETE("/:id", userHandler.DeleteUser)
	userGroup.GET("/me", userHandler.GetMe)

	authGroup := routeGroup.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/refresh-token", jwtRefreshMiddleware(authHandler.Refresh))
	authGroup.POST("/logout", jwtRefreshMiddleware(authHandler.Logout))

	todoGroup := routeGroup.Group("/todos")
	todoGroup.Use(jwtAccessMiddleware)
	todoGroup.POST("", todoHandler.CreateTodo)
	todoGroup.GET("/:userId", todoHandler.GetTodosByUser)
	todoGroup.PATCH("/:id", todoHandler.UpdateTodo)
	todoGroup.PATCH("/state/:stateId", todoHandler.UpdateTodoState)
	todoGroup.DELETE("/:id", todoHandler.DeleteTodo)
}
