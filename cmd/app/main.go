package main

import (
	"auth/internal/config"
	"auth/internal/database"
	"auth/internal/logger"
	"auth/internal/router"
	"context"

	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()
	// init config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	// init logger
	logger.InitLogger(cfg.LOG_LEVEL, cfg.MODE)
	// init database
	pool := database.ConnectDB(cfg, ctx)
	// init echo
	e := echo.New()
	router.InitRouter(e, pool, cfg)

	e.Logger.Fatal(e.Start(":8080"))

	// graceful shotdown
}
