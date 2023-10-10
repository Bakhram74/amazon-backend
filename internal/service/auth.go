package service

import (
	"context"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/Bakhram74/amazon-backend.git/internal/repository"
	"github.com/Bakhram74/amazon-backend.git/pkg/utils"
	"github.com/google/uuid"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (service *AuthService) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	var err error
	arg.Password, err = utils.HashPassword(arg.Password)
	if err != nil {
		return db.User{}, nil
	}
	return service.repo.CreateUser(ctx, arg)
}

func (service *AuthService) GetUser(ctx context.Context, email string) (db.User, error) {

	return service.repo.GetUser(ctx, email)
}

func (service *AuthService) CreateSession(ctx context.Context, arg db.CreateSessionParams) error {
	return service.repo.CreateSession(ctx, arg)
}
func (service *AuthService) GetSession(ctx context.Context, id uuid.UUID) (db.Session, error) {
	return service.repo.GetSession(ctx, id)
}
