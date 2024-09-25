package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/validator"
)

type Category struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	Name        string     `json:"name" db:"name"`
	Slug        string     `json:"slug" db:"slug"`
	Description *string    `json:"description,omitempty" db:"description"`
	ImageUrl    string     `json:"image_url" db:"image_url"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	CreatedByID uuid.UUID  `json:"created_by_id" db:"created_by_id"`
	UpdatedByID uuid.UUID  `json:"updated_by_id" db:"updated_by_id"`
	Version     int        `json:"version" db:"version"`
}

func ValidateCategory(v *validator.Validator, category *Category) {
	v.Check(category.Name != "", "name", "must be provided")
	v.Check(category.Slug != "", "slug", "must be provided")
	v.Check(category.ImageUrl != "", "image_url", "must be provided")
	v.Check(category.CreatedByID != uuid.Nil, "created_by_id", "must be provided")
	v.Check(category.UpdatedByID != uuid.Nil, "updated_by_id", "must be provided")

	v.Check(len([]byte(category.Name)) <= 50, "name", "must not be more than 50 bytes long")
	v.Check(len([]byte(category.Slug)) <= 50, "slug", "must not be more than 50 bytes long")

	v.Check(category.UpdatedAt.After(time.Now()), "updated_at", "can not be in the future")
	v.Check(category.CreatedAt.After(category.UpdatedAt), "created_at", "can not be later than updated_at")
}

type CategoryModel struct {
	DBPOOL *pgxpool.Pool
}

func (c CategoryModel) Insert(category *Category) error {
	query := `
		INSERT INTO categories (
			parent_id, 
			name, 
			slug, 
			description, 
			image_url, 
			created_by_id, 
			updated_by_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, name, slug, created_at, version`

	args := []interface{}{
		category.ParentID,
		category.Name,
		category.Slug,
		category.Description,
		category.ImageUrl,
		category.CreatedByID,
		category.UpdatedByID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := c.DBPOOL.QueryRow(ctx, query, args...).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.CreatedAt,
		&category.Version,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return common.PgErrWithCode(pgErr)
		}
		return err
	}

	return nil
}

func (c CategoryModel) Get(id uuid.UUID) (*Category, error) {
	query := `
		SELECT * 
		FROM categories
		WHERE id = $1
	`
	var category Category

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := c.DBPOOL.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.ParentID,
		&category.Name,
		&category.Slug,
		&category.ImageUrl,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.CreatedByID,
		&category.UpdatedByID,
		&category.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (c CategoryModel) GetAll(
	names, slugs []string,
	parentIds []uuid.UUID,
	search *string,
	createdAtFrom, createdAtUpTo *time.Time,
	updatedAtFrom, updatedAtUpTo *time.Time,
	createdByIDs []uuid.UUID,
	updatedByIDs []uuid.UUID,
	sorts []string,
	sortSafeList []string,
	page, pageSize *int,
) ([]*Category, common.Metadata, error) {
	query := `
	SELECT (
		count(*) OVER(),
		id, 
		name, 
		parent_id, 
		slug, 
		description, 
		image_url, 
		created_at,
		updated_at,
		created_by_id, 
		updated_by_id
	)
	FROM categories
	WHERE 1=1`

	args := []interface{}{}
	argCounter := 1

	if len(names) > 0 {
		query += fmt.Sprintf(" AND LOWER(name) = ANY($%d)", argCounter)
		args = append(args, names)
		argCounter++
	}

	if len(slugs) > 0 {
		query += fmt.Sprintf(" AND LOWER(slug) = ANY($%d)", argCounter)
		args = append(args, slugs)
		argCounter++
	}

	if len(parentIds) > 0 {
		query += fmt.Sprintf(" AND parent_id = ANY($%d)", argCounter)
		args = append(args, parentIds)
		argCounter++
	}

	if search != nil {
		query += fmt.Sprintf(" AND to_tsvector('simple', name) @@ plainto_tsquery('simple', $%d)", argCounter)
		args = append(args, *search)
		argCounter++
	}

	if createdAtFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argCounter)
		args = append(args, *createdAtFrom)
		argCounter++
	}

	if createdAtUpTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argCounter)
		args = append(args, *createdAtUpTo)
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
		query += "ORDER BY"
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

	rows, err := c.DBPOOL.Query(ctx, query, args...)
	if err != nil {
		return nil, common.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	categories := []*Category{}

	for rows.Next() {
		var category Category
		err := rows.Scan(
			&totalRecords,
			&category.ID,
			&category.ParentID,
			&category.Name,
			&category.Slug,
			&category.ImageUrl,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.CreatedByID,
			&category.UpdatedByID,
			&category.Version,
		)
		if err != nil {
			return nil, common.Metadata{}, err
		}
		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, common.Metadata{}, err
	}

	metadata := common.CalculateMetadata(totalRecords, *page, *pageSize)

	return categories, metadata, nil
}

func (c CategoryModel) Update(category *Category) error {
	query := `
		UPDATE categories
		SET 
			parent_id = $1,
			name = $2,
			slug = $3,
			description = $4
			image_url = $5,
			updated_by_id = $6,
			version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING id, parent_id, name, slug, description, image_url, updated_by_id, version
	`

	args := []interface{}{
		category.ParentID,
		category.Name,
		category.Slug,
		category.Description,
		category.ImageUrl,
		category.UpdatedByID,
		category.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := c.DBPOOL.QueryRow(ctx, query, args...).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.ImageUrl,
		&category.UpdatedByID,
		&category.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (c CategoryModel) Delete(id uuid.UUID) error {
	query := `
		DELETE FROM categories
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := c.DBPOOL.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// CREATE TABLE IF NOT EXISTS products_categories (
//     id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
//     product_id uuid NOT NULL,
//     category_id uuid NOT NULL,

//     UNIQUE (product_id, category_id)
// );
