package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

type UserRepository struct {
	DBPOOL *pgxpool.Pool
}

func (r UserRepository) Create(user *data.User) error {
	args := []interface{}{
		user.Phone,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Patronomic,
		user.Email,
		user.IsActive,
	}

	query := `
		INSERT INTO users (
			phone, 
			password_hash, 
			first_name,
			last_name,
			patronmic,
			email,
			is_active 
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, phone, first_name, last_name, patronomic, email, is_active, created_by_id
	`

	if user.CreatedByID != nil {
		query = `
		INSERT INTO users (
			phone, 
			password_hash, 
			first_name,
			last_name,
			patronmic,
			email,
			is_active, 
			created_by_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, phone, first_name, last_name, patronomic, email, is_active, created_by_id
		`
		args = append(args, *user.CreatedByID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.DBPOOL.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Phone,
		&user.FirstName,
		&user.LastName,
		&user.Patronomic,
		&user.Email,
		&user.IsActive,
		&user.CreatedByID,
	)

	// TODO: add error handling - pgErrs
	return err
}
