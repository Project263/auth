package services

import (
	"auth/internal/models"
	"auth/internal/repositories"
	"context"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUserById(ctx context.Context, id string) (models.User, error) {
	return s.repo.GetUserById(ctx, id)
}
