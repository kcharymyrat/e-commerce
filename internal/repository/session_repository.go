package repository

import "github.com/jackc/pgx/v5/pgxpool"

type SessionRepository struct {
	DBPOOL *pgxpool.Pool
}
