package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

type SessionRepository struct {
	DBPOOL *pgxpool.Pool
}

func (r SessionRepository) Create(session *data.Session) error {
	query := `INSERT INTO sessions (id, user_phone, refresh_token, expires_at)
	VALUES ($1, $2, $3, $4)
	`

	args := []interface{}{
		&session.ID,
		&session.UserPhone,
		&session.RefreshToken,
		&session.ExpiresAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.DBPOOL.QueryRow(ctx, query, args...).Scan(
		&session.ID,
		&session.UserPhone,
		&session.RefreshToken,
		&session.IsRevoked,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
}

func (r SessionRepository) GetByID(id uuid.UUID) (*data.Session, error) {
	query := `SELECT * FROM sessions
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session data.Session

	err := r.DBPOOL.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.UserPhone,
		&session.RefreshToken,
		&session.IsRevoked,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, common.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &session, nil
}

func (r SessionRepository) GetByRefreshToken(refreshToken string) (*data.Session, error) {
	query := `SELECT * FROM sessions
	WHERE refresh_token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session data.Session

	err := r.DBPOOL.QueryRow(ctx, query, refreshToken).Scan(
		&session.ID,
		&session.UserPhone,
		&session.RefreshToken,
		&session.IsRevoked,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, common.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &session, nil
}

func (r SessionRepository) RevokeSessionByID(id uuid.UUID) error {
	query := `UPDATE sessions
	SET is_revoked = true
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session data.Session

	err := r.DBPOOL.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.UserPhone,
		&session.RefreshToken,
		&session.IsRevoked,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return common.ErrEditConflict
		default:
			return err
		}
	}

	fmt.Println("session =", session)

	return nil
}

func (r SessionRepository) DeleteByID(id uuid.UUID) error {
	query := `
	DELETE FROM categories
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := r.DBPOOL.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected < 1 {
		return common.ErrRecordNotFound
	}

	return nil
}
