package repository

import (
	db "github.com/Bakhram74/amazon-backend.git/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*db.Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {
	return Store{
		connPool: connPool,
		Queries:  db.New(connPool),
	}
}
