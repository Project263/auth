package database

import (
	"auth/internal/config"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func ConnectDB(cfg *config.Config, ctx context.Context) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, cfg.POSTGRES_DSN)

	if err != nil {
		logrus.Errorf("ошибка подключения к базе данных: %s", err.Error())
		return nil
	}

	err = pool.Ping(ctx)

	if err != nil {
		logrus.Errorf("ошибка подключения к базе данных: %s", err.Error())
		return nil
	}

	return pool
}
