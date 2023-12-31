package repository

import (
	"context"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/google/uuid"
)

type AuthRepository struct {
	store Store
}

func NewAuthRepository(store Store) *AuthRepository {
	return &AuthRepository{
		store: store,
	}
}

func (a AuthRepository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	user, err := a.store.CreateUser(ctx, arg)
	if err != nil {
		return db.User{}, err
	}
	return user, err
}

func (a AuthRepository) GetUser(ctx context.Context, email string) (db.User, error) {
	user, err := a.store.GetUser(ctx, email)
	if err != nil {
		return db.User{}, err
	}
	return user, err
}

func (a AuthRepository) CreateSession(ctx context.Context, arg db.CreateSessionParams) error {
	err := a.store.CreateSession(ctx, arg)
	if err != nil {
		return err
	}
	return nil
}
func (a AuthRepository) GetSession(ctx context.Context, id uuid.UUID) (db.Session, error) {
	getSession, err := a.store.GetSession(ctx, id)
	if err != nil {
		return db.Session{}, err
	}
	return getSession, err
}
