-- name: CreateReview :one
INSERT INTO "review" (
                      "user_id",
                      "product_id",
                     "rating",
                       "text"

)
VALUES ($1, $2,$3,$4) RETURNING *;

