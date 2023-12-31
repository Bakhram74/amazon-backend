// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateReview(ctx context.Context, arg CreateReviewParams) (Review, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
}

var _ Querier = (*Queries)(nil)
