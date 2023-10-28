package service

import (
	"context"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/Bakhram74/amazon-backend.git/internal/repository"
)

type ProfileService struct {
	repo repository.Profile
}

func NewProfileService(repo repository.Profile) *ProfileService {
	return &ProfileService{repo: repo}
}

func (p ProfileService) GetUserByID(ctx context.Context, id int64) (db.User, error) {
	return p.repo.GetUserByID(ctx, id)
}
