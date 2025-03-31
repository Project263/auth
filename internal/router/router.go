package router

import (
	"auth/internal/config"
	"auth/internal/handlers"
	"auth/internal/repositories"
	"auth/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRouter(e *echo.Echo, db *pgxpool.Pool, cfg *config.Config) {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, cfg)

	e.Use(middleware.Logger())
	e.POST("/login", userHandler.Login)
}
