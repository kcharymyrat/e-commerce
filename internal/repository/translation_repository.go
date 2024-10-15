package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
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

func (r TranslationRepository) List(
	langCodes, tableNames, fieldNames []string,
	entityIDs []uuid.UUID,
	search *string,
	createdAtFrom, createdAtUpTo *time.Time,
	updatedAtFrom, updatedAtUpTo *time.Time,
	createdByIDs []uuid.UUID,
	updatedByIDs []uuid.UUID,
	sorts []string,
	sortSafeList []string,
	page, pageSize *int,
) ([]*data.Translation, common.Metadata, error) {

	query := `
		SELECT COUNT(*) OVER(), * 
		FROM translations
		WHERE 1=1
	`

	args := []interface{}{}
	argCounter := 1

	if len(langCodes) > 0 {
		query += fmt.Sprintf(" AND language_code = ANY($%d)", argCounter)
		args = append(args, langCodes)
		argCounter++
	}

	if len(tableNames) > 0 {
		query += fmt.Sprintf(" AND table_names = ANY($%d)", argCounter)
		args = append(args, tableNames)
		argCounter++
	}

	if len(fieldNames) > 0 {
		query += fmt.Sprintf(" AND field_name = ANY($%d)", argCounter)
		args = append(args, fieldNames)
		argCounter++
	}

	if len(entityIDs) > 0 {
		query += fmt.Sprintf(" AND entity_id = ANY($%d)", argCounter)
		args = append(args, entityIDs)
		argCounter++
	}

	if createdAtFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argCounter)
		args = append(args, createdAtFrom)
		argCounter++
	}

	if createdAtUpTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argCounter)
		args = append(args, createdAtFrom)
		argCounter++
	}

	if updatedAtFrom != nil {
		query += fmt.Sprintf(" AND updated_at >= $%d", argCounter)
		args = append(args, *updatedAtFrom)
		argCounter++
	}

	if updatedAtUpTo != nil {
		query += fmt.Sprintf(" AND updated_at <= $%d", argCounter)
		args = append(args, *updatedAtUpTo)
		argCounter++
	}

	if len(createdByIDs) > 0 {
		query += fmt.Sprintf(" AND created_by = ANY($%d)", argCounter)
		args = append(args, createdByIDs)
		argCounter++
	}

	if len(updatedByIDs) > 0 {
		query += fmt.Sprintf(" AND updated_by_id = ANY($%d)", argCounter)
		args = append(args, updatedByIDs)
		argCounter++
	}

	if len(sorts) > 0 {
		query += " ORDER BY"
		for _, sort := range sorts {
			direction := "ASC"
			sortField := strings.TrimSpace(strings.ToLower(sort))
			if strings.HasPrefix(sort, "-") {
				direction = "DESC"
				sortField = strings.TrimPrefix(sort, "-")
			}
			query += fmt.Sprintf(" %s %s,", sortField, direction)
		}
		query += " id ASC"
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.DBPOOL.Query(ctx, query, args)
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

	metadata := common.CalculateMetadata(totalRecords, *page, *pageSize)

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
