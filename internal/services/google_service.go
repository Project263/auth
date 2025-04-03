package services

import (
	"auth/internal/models"
	"auth/internal/repositories"
	"context"
)

type GoogleService struct {
	repo *repositories.GoogleRepository
}

func NewGoogleService(repo *repositories.GoogleRepository) *GoogleService {
	return &GoogleService{repo: repo}
}

func (s *GoogleService) CreateUserByGoogle(ctx context.Context, userData models.User) (string, error) {
	return s.repo.CreateUserByGoogle(ctx, userData)
}
