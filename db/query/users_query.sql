-- name: CreateUser :one
INSERT INTO users ("username",
                   "phone_number",
                   "hashed_password")
VALUES ($1, $2, $3) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE phone_number = $1 LIMIT 1;


-- name: PartialUpdateUser :one
UPDATE users
SET username = CASE WHEN @update_username::boolean THEN @username::TEXT ELSE username END,
    phone_number  = CASE WHEN @update_phone_number::boolean THEN @phone_number::TEXT ELSE phone_number END,
    hashed_password  = CASE WHEN @update_password::boolean THEN @password::TEXT ELSE hashed_password END
WHERE id = @id
RETURNING *;

