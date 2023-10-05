package service

import (
	"context"
	db "github.com/Bakhram74/amazon.git/db/sqlc"
	"github.com/Bakhram74/amazon.git/internal/repository"

	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUser(ctx context.Context, phoneNumber string) (db.User, error)
	CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (db.Session, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
