// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: sessions_query.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSession = `-- name: CreateSession :exec
INSERT INTO sessions (
    id,
    user_id,
    refresh_token,
    user_agent,
    client_ip,
    is_blocked,
    expires_at
) VALUES (
             $1, $2, $3, $4, $5, $6,$7
         )
`

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) error {
	_, err := q.db.Exec(ctx, createSession,
		arg.ID,
		arg.UserID,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ExpiresAt,
	)
	return err
}

const getSession = `-- name: GetSession :one
SELECT id, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at, user_id FROM sessions
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, id uuid.UUID) (Session, error) {
	row := q.db.QueryRow(ctx, getSession, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UserID,
	)
	return i, err
}
