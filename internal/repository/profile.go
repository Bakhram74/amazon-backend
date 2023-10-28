package repository

import (
	"context"
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
)

type ProfileRepository struct {
	store Store
}

func NewProfileRepository(store Store) *ProfileRepository {
	return &ProfileRepository{
		store: store,
	}
}
func (p ProfileRepository) GetUserByID(ctx context.Context, id int64) (db.User, error) {
	return p.store.GetUserByID(ctx, id)
}
