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

func (r *UserRepository) CreateUser(ctx context.Context, email, name, password string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		logrus.Error(err)
		return err
	}
	query, args, err := squirrel.Select("id").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		logrus.Error(err)
		return err
	}

	row := tx.QueryRow(ctx, query, args...)
	var id string
	if err := row.Scan(&id); err != nil {
		logrus.Error(err)
		return err
	}

	logrus.Info(id)

	return nil
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

func (r *UserRepository) GetUserById(ctx context.Context, id string) (models.User, error) {
	query, args, err := squirrel.Select("id", "username", "password", "email", "role").
		From("users").
		Where(squirrel.Eq{"id": id}).
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
