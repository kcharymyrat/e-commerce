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

func (r LanguageRepository) Get(id uuid.UUID) (*data.Language, error) {
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

func (r LanguageRepository) List(page, pageSize *int) ([]*data.Language, common.Metadata, error) {
	query := `
		SELECT *
		FROM languages
	`

	argCounter := 1
	args := []interface{}{}

	fallbackPageSize := 20 // FIXME: make a constant number
	if pageSize != nil {
		query += fmt.Sprintf(" LIMIT $%d", argCounter)
		args = append(args, *pageSize)
		argCounter++
		fallbackPageSize = *pageSize
	} else {
		query += fmt.Sprintf(" LIMIT %d", fallbackPageSize)
	}

	defaultPage := 1 // FIXME: make a constant number
	if page != nil {
		offset := fallbackPageSize * (*page - 1)
		query += fmt.Sprintf(" OFFSET $%d", argCounter)
		args = append(args, offset)
		argCounter++
	} else {
		query += fmt.Sprintf(" LIMIT %d", defaultPage)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := r.DBPOOL.Query(ctx, query, args...)
	if err != nil {
		return nil, common.Metadata{}, err
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
			return nil, common.Metadata{}, err
		}
		languages = append(languages, &language)
		totalRecords++
	}

	if err = rows.Err(); err != nil {
		return nil, common.Metadata{}, err
	}

	metadata := common.CalculateMetadata(totalRecords, *page, *pageSize)

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
		WHERE version = $4
		RETURNING id, code, name, version	
	`

	args := []interface{}{
		language.Code,
		language.Name,
		language.UpdatedByID,
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
