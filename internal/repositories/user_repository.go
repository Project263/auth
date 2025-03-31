package repositories

import (
	"auth/internal/models"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	query, args, err := squirrel.Select("id", "username", "password", "email", "role").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		logrus.Error(err)
		return models.User{}, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	var user models.User
	if err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Role); err != nil {
		logrus.Error(err)
		return models.User{}, err
	}

	return user, nil
}
