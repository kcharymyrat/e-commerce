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
	"github.com/kcharymyrat/e-commerce/internal/filters"
)

type TranslationRepository struct {
	DBPOOL *pgxpool.Pool
}

func (r TranslationRepository) Create(translation *data.Translation) error {
	query := `
		INSERT INTO translations (
			language_code,
			entity_id,
			table_name,
			field_name,
			translated_field_name,
			translated_value,
			created_by_id,
			updated_by_id,
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, language_code, table_name, field_name, translated_field_name, translated_value, version
	`

	args := []interface{}{
		&translation.LanguageCode,
		&translation.EntityID,
		&translation.TableName,
		&translation.FieldName,
		&translation.TranslatedFieldName,
		&translation.TranslatedValue,
		&translation.CreatedByID,
		&translation.UpdatedByID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.DBPOOL.QueryRow(ctx, query, args).Scan(
		&translation.ID,
		&translation.LanguageCode,
		&translation.TableName,
		&translation.FieldName,
		&translation.TranslatedFieldName,
		&translation.TranslatedValue,
		&translation.Version,
	)
}

func (r TranslationRepository) GetByID(id uuid.UUID) (*data.Translation, error) {
	query := `
		SELECT *
		FROM translations
		WHERE id = $1`

	var translation data.Translation

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.DBPOOL.QueryRow(ctx, query, id).Scan(
		&translation.ID,
		&translation.LanguageCode,
		&translation.EntityID,
		&translation.TableName,
		&translation.FieldName,
		&translation.TranslatedFieldName,
		&translation.TranslatedValue,
		&translation.CreatedAt,
		&translation.UpdatedAt,
		&translation.CreatedByID,
		&translation.UpdatedByID,
		&translation.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, common.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &translation, nil
}

func (r TranslationRepository) List(f *requests.TranslationsAdminFilters) ([]*data.Translation, common.Metadata, error) {

	query := `
		SELECT COUNT(*) OVER(), * 
		FROM translations
		WHERE 1=1
	`

	args := []interface{}{}
	argCounter := 1

	if len(f.LanguageCodes) > 0 {
		query += fmt.Sprintf(" AND language_code = ANY($%d)", argCounter)
		args = append(args, f.LanguageCodes)
		argCounter++
	}

	if len(f.TableNames) > 0 {
		query += fmt.Sprintf(" AND table_name = ANY($%d)", argCounter)
		args = append(args, f.TableNames)
		argCounter++
	}

	if len(f.FieldNames) > 0 {
		query += fmt.Sprintf(" AND field_name = ANY($%d)", argCounter)
		args = append(args, f.FieldNames)
		argCounter++
	}

	if len(f.EntityIDs) > 0 {
		query += fmt.Sprintf(" AND entity_id = ANY($%d)", argCounter)
		args = append(args, f.EntityIDs)
		argCounter++
	}

	if f.Search != nil {
		query += fmt.Sprintf(` 
		AND (
			to_tsvector('simple', id) @@ plainto_tsquery('simple', $%d) OR 
			to_tsvector('simple', table_name) @@ plainto_tsquery('simple', $%d) OR 
			to_tsvector('simple', field_name) @@ plainto_tsquery('simple', $%d) OR  
			to_tsvector('simple', translated_field_name) @@ plainto_tsquery('simple', $%d) OR 
			to_tsvector('simple', translated_value) @@ plainto_tsquery('simple', $%d)
		)`, argCounter, argCounter, argCounter, argCounter, argCounter)
		args = append(args, *f.Search)
		argCounter++
	}

	filters.AddCreatedUpdateAtFilterToSQL(&f.CreatedUpdatedAtFilter, &query, &argCounter, args)
	filters.AddCreatedUpdateByFilterToSQL(&f.CreatedUpdatedByFilter, &query, &argCounter, args)
	filters.AddSortListFilterToSQL(&f.SortListFilter, &query)
	filters.AddPaginationFilterToSQL(&f.PaginationFilter, &query, &argCounter, args)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.DBPOOL.Query(ctx, query, args...)
	if err != nil {
		return nil, common.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	trs := []*data.Translation{}

	for rows.Next() {
		var tr data.Translation
		err := rows.Scan(
			&totalRecords,
			&tr.ID,
			&tr.LanguageCode,
			&tr.EntityID,
			&tr.TableName,
			&tr.FieldName,
			&tr.TranslatedFieldName,
			&tr.TranslatedValue,
			&tr.CreatedAt,
			&tr.UpdatedAt,
			&tr.CreatedByID,
			&tr.UpdatedByID,
			&tr.Version,
		)
		if err != nil {
			return nil, common.Metadata{}, err
		}
		trs = append(trs, &tr)
	}

	if err = rows.Err(); err != nil {
		return nil, common.Metadata{}, err
	}

	metadata := common.CalculateMetadata(totalRecords, *f.Page, *f.PageSize)

	return trs, metadata, nil

}

func (r TranslationRepository) Update(tr *data.Translation) error {
	query := `
		UPDATE translations
		SET
			language_code = $1,
			entity_id = $2,
			table_name = $3,
			field_name = $4,
			translated_field_name = $5,
			translated_values = $6
			updated_by_id = $7,
			version = version + 1
		WHERE id = $8 AND version = $9
		RETURNING id, language_code, version
	`

	args := []interface{}{
		&tr.LanguageCode,
		&tr.EntityID,
		&tr.TableName,
		&tr.FieldName,
		&tr.TranslatedFieldName,
		&tr.TranslatedValue,
		&tr.UpdatedByID,
		&tr.ID,
		&tr.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.DBPOOL.QueryRow(ctx, query, args).Scan(
		&tr.LanguageCode,
		&tr.EntityID,
		&tr.TableName,
		&tr.FieldName,
		&tr.TranslatedFieldName,
		&tr.TranslatedValue,
		&tr.UpdatedByID,
		&tr.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return common.ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (r TranslationRepository) Delete(id uuid.UUID) error {
	query := `
		DELETE FROM translations
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func (r TranslationRepository) GetByEntityIDLangCodeFieldName(entityID uuid.UUID, languageCode, fieldName string) (*data.Translation, error) {
	query := `
		SELECT *
		FROM translations
		WHERE entity_id = $1 AND language_code = $2 AND field_name = $4`

	var translation data.Translation

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.DBPOOL.QueryRow(ctx, query, entityID, languageCode, fieldName).Scan(
		&translation.ID,
		&translation.LanguageCode,
		&translation.EntityID,
		&translation.TableName,
		&translation.FieldName,
		&translation.TranslatedFieldName,
		&translation.TranslatedValue,
		&translation.CreatedAt,
		&translation.UpdatedAt,
		&translation.CreatedByID,
		&translation.UpdatedByID,
		&translation.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, common.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &translation, nil
}
