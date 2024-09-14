package data

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/internal/validator"
)

type Category struct {
	ID          uuid.UUID `json:"id"`
	Parent      uuid.UUID `json:"parent,omitempty"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	ImageUrl    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	CreatedByID uuid.UUID `json:"created_by_id"`
	UpdatedByID uuid.UUID `json:"updated_by_id"`
}

type CategoryPublic struct {
	ID          uuid.UUID `json:"id"`
	Parent      uuid.UUID `json:"parent,omitempty"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	ImageUrl    string    `json:"image_url"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	CreatedByID uuid.UUID `json:"-"`
	UpdatedByID uuid.UUID `json:"-"`
}

func ValidateCategory(v *validator.Validator, category *Category) {
	v.Check(category.Name != "", "name", "must be provided")
	v.Check(category.Slug != "", "slug", "must be provided")

	v.Check(len(category.Name) <= 500, "name", "must be more than 500 bytes long")
	v.Check(len(category.Slug) <= 500, "slug", "must be more than 500 bytes long")

	v.Check(category.UpdatedAt.After(time.Now()), "updated_at", "can not be in the future")
	v.Check(category.CreatedAt.After(category.UpdatedAt), "created_at", "can not be later than updated_at")
}

type CategoryModel struct {
	DBPOOL *pgxpool.Pool
}

func (c CategoryModel) Insert(category *Category) error {
	query := `
		INSERT INTO categories (name, parent, )
	`
}

func (c CategoryModel) Get(id uuid.UUID) (*Category, error) {
	return nil, nil
}

func (c CategoryModel) Update(category *Category) error {
	return nil
}

func (c CategoryModel) Delete(id uuid.UUID) error {
	return nil
}

// CREATE TABLE IF NOT EXISTS categories (
//     id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
//     parent uuid,
//     name varchar(50),
//     slug varchar(50),
//     description text,
//     image_url text NOT NULL,
//     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
//     updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
//     created_by_id uuid NOT NULL,
//     updated_by_id uuid NOT NULL,

//     CHECK (updated_at >= created_at)
// );

// CREATE TABLE IF NOT EXISTS products_categories (
//     id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
//     product_id uuid NOT NULL,
//     category_id uuid NOT NULL,

//     UNIQUE (product_id, category_id)
// );
