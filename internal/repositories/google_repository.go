package repositories

import (
	"auth/internal/models"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type GoogleRepository struct {
	db *pgxpool.Pool
}

func NewGoogleRepository(db *pgxpool.Pool) *GoogleRepository {
	return &GoogleRepository{db: db}
}

func (r *GoogleRepository) CreateUserByGoogle(ctx context.Context, userData models.User) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	query, args, err := squirrel.Select("id").
		From("users").
		Where(squirrel.Eq{"email": userData.Email}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	row := tx.QueryRow(ctx, query, args...)
	var userID string
	if err := row.Scan(&userID); err == pgx.ErrNoRows {
		query, args, err = squirrel.Insert("users").
			Columns("email", "username", "avatar", "password", "role").
			Values(userData.Email, userData.Username, userData.Avatar, "", "user").
			Suffix("RETURNING id").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			logrus.Error(err)
			return "", err
		}

		var userID string
		if err := r.db.QueryRow(ctx, query, args...).Scan(&userID); err != nil {
			logrus.Error(err)
			return "", err
		}
		return userID, nil
	}

	return userID, nil
}
