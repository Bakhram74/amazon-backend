package service

import (
	"context"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/Bakhram74/amazon-backend.git/internal/repository"
	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUser(ctx context.Context, email string) (db.User, error)
	CreateSession(ctx context.Context, arg db.CreateSessionParams) error
	GetSession(ctx context.Context, id uuid.UUID) (db.Session, error)
}

type Profile interface {
	GetUserByID(ctx context.Context, id int64) (db.User, error)
}

type Service struct {
	Authorization
	Profile
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Profile:       NewProfileService(repos.Profile),
	}
}
