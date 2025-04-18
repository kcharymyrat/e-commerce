package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/types"
)

type LanguageRepository struct {
	DBPOOL *pgxpool.Pool
}

func (r LanguageRepository) Create(language *data.Language) error {
	query := `
		INSERT INTO languages (
			code,
			name,
			created_by_id,
			updated_by_id,
		) VALUES $1, $2, $3, $4
		RETURNING id, code, name, version
	`

	args := []interface{}{
		language.Code,
		language.Name,
		language.CreatedByID,
		language.CreatedByID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return r.DBPOOL.QueryRow(ctx, query, args).Scan(
		&language.ID,
		&language.Code,
		&language.Name,
		&language.Version,
	)
}

func (r LanguageRepository) GetByID(id uuid.UUID) (*data.Language, error) {
	query := `
		SELECT * 
		FROM languages
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	language := data.Language{}

	err := r.DBPOOL.QueryRow(ctx, query, id).Scan(
		&language.ID,
		&language.Code,
		&language.Name,
		&language.CreatedAt,
		&language.UpdatedAt,
		&language.CreatedByID,
		&language.UpdatedByID,
		&language.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, common.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &language, nil
}

func (r LanguageRepository) GetByCode(code string) (*data.Language, error) {
	query := `
		SELECT * 
		FROM languages
		WHERE code = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	language := data.Language{}

	err := r.DBPOOL.QueryRow(ctx, query, code).Scan(
		&language.ID,
		&language.Code,
		&language.Name,
		&language.CreatedAt,
		&language.UpdatedAt,
		&language.CreatedByID,
		&language.UpdatedByID,
		&language.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, common.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &language, nil
}

func (r LanguageRepository) List(f *requests.LanguagesAdminFilters) ([]*data.Language, types.PaginationMetadata, error) {
	query := `
		SELECT *
		FROM languages
	`

	argCounter := 1
	args := []interface{}{}

	fallbackPageSize := 20 // FIXME: make a constant number
	if f.PageSize != nil {
		query += fmt.Sprintf(" LIMIT $%d", argCounter)
		args = append(args, *f.PageSize)
		argCounter++
		fallbackPageSize = *f.PageSize
	} else {
		f.PageSize = &fallbackPageSize
		query += fmt.Sprintf(" LIMIT %d", fallbackPageSize)
	}

	defaultPage := 1 // FIXME: make a constant number
	if f.Page != nil {
		offset := fallbackPageSize * (*f.Page - 1)
		query += fmt.Sprintf(" OFFSET $%d", argCounter)
		args = append(args, offset)
		argCounter++
	} else {
		f.Page = &defaultPage
		query += fmt.Sprintf(" OFFSET %d", defaultPage)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := r.DBPOOL.Query(ctx, query, args...)
	if err != nil {
		return nil, types.PaginationMetadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	languages := []*data.Language{}

	for rows.Next() {
		var language data.Language
		err := rows.Scan(
			&language.ID,
			&language.Code,
			&language.Name,
			&language.CreatedAt,
			&language.UpdatedAt,
			&language.CreatedByID,
			&language.UpdatedByID,
			&language.Version,
		)
		if err != nil {
			return nil, types.PaginationMetadata{}, err
		}
		languages = append(languages, &language)
		totalRecords++
	}

	if err = rows.Err(); err != nil {
		return nil, types.PaginationMetadata{}, err
	}

	metadata := common.CalculateMetadata(totalRecords, *f.Page, *f.PageSize)

	return languages, metadata, nil
}

func (r LanguageRepository) Update(language *data.Language) error {
	query := `
		UPDATE languages
		SET
			code = $1,
			name = $2,
			updated_by_id = $3,
			version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING id, code, name, version	
	`

	args := []interface{}{
		language.Code,
		language.Name,
		language.UpdatedByID,
		language.ID,
		language.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := r.DBPOOL.QueryRow(ctx, query, args).Scan(
		&language.ID,
		&language.Code,
		&language.Name,
		&language.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return common.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r LanguageRepository) Delete(id uuid.UUID) error {
	query := `
		DELETE FROM languages
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
