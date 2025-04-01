package repositories

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type GoogleRepository struct {
	db *pgxpool.Pool
}

func NewGoogleRepository(db *pgxpool.Pool) *GoogleRepository {
	return &GoogleRepository{db: db}
}

func (r *GoogleRepository) CreateUser(ctx context.Context, email, username, password string) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	query, args, err := squirrel.Select("id").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	row := tx.QueryRow(ctx, query, args...)
	var userID string
	if err := row.Scan(&userID); err == sql.ErrNoRows {
		query, args, err = squirrel.Insert("users").
			Columns("email", "username", "password", "role").
			Values(email, username, "", "user").
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
