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
	"github.com/kcharymyrat/e-commerce/internal/types"
)

type CategoryRepository struct {
	DBPOOL *pgxpool.Pool
}

func (r CategoryRepository) Create(category *data.Category) error {
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

	return r.DBPOOL.QueryRow(ctx, query, args...).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.CreatedAt,
		&category.Version,
	)
}

func (r CategoryRepository) GetByID(id uuid.UUID) (*data.Category, error) {
	query := `
	SELECT * 
	FROM categories
	WHERE id = $1
`
	var category data.Category

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := r.DBPOOL.QueryRow(ctx, query, id).Scan(
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
		case errors.Is(err, pgx.ErrNoRows):
			return nil, common.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (r CategoryRepository) GetBySlug(slug string) (*data.Category, error) {
	query := `
	SELECT * 
	FROM categories
	WHERE slug = $1
`
	var category data.Category

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := r.DBPOOL.QueryRow(ctx, query, slug).Scan(
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
		case errors.Is(err, pgx.ErrNoRows):
			return nil, common.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (r CategoryRepository) List(f *requests.CategoriesAdminFilters) ([]*data.Category, types.PaginationMetadata, error) {
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
			updated_by_id,
			version
		)
		FROM categories
		WHERE 1=1
	`

	args := []interface{}{}
	argCounter := 1

	if len(f.Names) > 0 {
		query += fmt.Sprintf(" AND name = ANY($%d)", argCounter)
		args = append(args, f.Names)
		argCounter++
	}

	if len(f.Slugs) > 0 {
		query += fmt.Sprintf(" AND LOWER(slug) = ANY($%d)", argCounter)
		args = append(args, f.Slugs)
		argCounter++
	}

	if len(f.ParentIDs) > 0 {
		query += fmt.Sprintf(" AND parent_id = ANY($%d)", argCounter)
		args = append(args, f.ParentIDs)
		argCounter++
	}

	if f.Search != nil {
		query += fmt.Sprintf(" AND to_tsvector('simple', name) @@ plainto_tsquery('simple', $%d)", argCounter)
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
		return nil, types.PaginationMetadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	categories := []*data.Category{}

	for rows.Next() {
		var category data.Category
		err := rows.Scan(
			&totalRecords,
			&category.ID,
			&category.Name,
			&category.ParentID,
			&category.Slug,
			&category.Description,
			&category.ImageUrl,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.CreatedByID,
			&category.UpdatedByID,
			&category.Version,
		)
		if err != nil {
			return nil, types.PaginationMetadata{}, err
		}
		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, types.PaginationMetadata{}, err
	}

	metadata := common.CalculateMetadata(totalRecords, *f.Page, *f.PageSize)

	return categories, metadata, nil
}

func (r CategoryRepository) Update(category *data.Category) error {
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

	err := r.DBPOOL.QueryRow(ctx, query, args...).Scan(
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
		case errors.Is(err, pgx.ErrNoRows):
			return common.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r CategoryRepository) DeleteByID(id uuid.UUID) error {
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

func (r CategoryRepository) DeleteBySlug(slug string) error {
	query := `
	DELETE FROM categories
	WHERE slug = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := r.DBPOOL.Exec(ctx, query, slug)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected < 1 {
		return common.ErrRecordNotFound
	}

	return nil
}
