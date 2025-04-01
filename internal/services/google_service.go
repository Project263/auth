package services

import (
	"auth/internal/repositories"
	"context"
)

type GoogleService struct {
	repo *repositories.GoogleRepository
}

func NewGoogleService(repo *repositories.GoogleRepository) *GoogleService {
	return &GoogleService{repo: repo}
}

func (s *GoogleService) CreateUser(ctx context.Context, email, username, password string) (string, error) {
	return s.repo.CreateUser(ctx, email, username, password)
}
