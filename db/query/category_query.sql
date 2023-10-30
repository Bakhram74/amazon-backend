-- name: CreateCategory :one
INSERT INTO "category" (
                     "id",
                       "name"

)
VALUES ($1, $2) RETURNING *;

