package service

import (
	"context"
	db "github.com/Bakhram74/amazon.git/db/sqlc"
	"github.com/Bakhram74/amazon.git/internal/repository"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (service *AuthService) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return service.repo.CreateUser(ctx, arg)
}

//func (service *AuthService) CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error) {
//	return service.repo.CreateSession(ctx, arg)
//}
//func (service *AuthService) GetSession(ctx context.Context, id uuid.UUID) (db.Session, error) {
//	return service.repo.GetSession(ctx, id)
//}
