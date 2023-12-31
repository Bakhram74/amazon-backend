// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: category_query.sql

package db

import (
	"context"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO "category" (
                     "id",
                       "name"

)
VALUES ($1, $2) RETURNING id, updated_at, created_at, name
`

type CreateCategoryParams struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	row := q.db.QueryRow(ctx, createCategory, arg.ID, arg.Name)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.Name,
	)
	return i, err
}
