package repository

import (
	"context"
	"errors"
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
			translated_value,
			created_by_id,
			updated_by_id,
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, language_code, table_name, field_name, translated_value, version`

	args := []interface{}{
		&translation.LanguageCode,
		&translation.EntityID,
		&translation.TableName,
		&translation.FieldName,
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
		&translation.TranslatedValue,
		&translation.Version,
	)
}

// TODO: put the index on DB for entityID uuid.UUID, languageCode, tableName, fieldName string
func (r TranslationRepository) Get(entityID uuid.UUID, languageCode, tableName, fieldName string) (*data.Translation, error) {
	query := `
		SELECT *
		FROM translations
		WHERE entity_id = $1 AND language_code = $2 AND table_name = $3 AND field_name = $4`

	var translation data.Translation

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.DBPOOL.QueryRow(ctx, query, entityID, languageCode, tableName, fieldName).Scan(
		&translation.ID,
		&translation.LanguageCode,
		&translation.EntityID,
		&translation.TableName,
		&translation.FieldName,
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
) ([]*data.Category, common.Metadata, error) {
	query := `
		SELECT COUNT(*) OVER(), * 
		FROM translations
		WHERE 1=1
	`

	args := []interface{}{}
	argCounter := 1

	if len(langCodes) > 0 {
		query := ` AND `
	}
}
