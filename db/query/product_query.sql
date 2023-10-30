-- name: CreateProduct :one
INSERT INTO "product" (
                       "category_id",
                       "name",
                    "slug",
                    "description",
                       "price",
                       "images"
)
VALUES ($1, $2, $3,$4,$5,$6) RETURNING *;

