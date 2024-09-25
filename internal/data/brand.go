package data

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Brand struct {
	ID          uuid.UUID `json:"id"`
	LogoUrl     string    `json:"logo_url"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedById uuid.UUID `json:"created_by_id"`
	UpdatedById uuid.UUID `json:"updated_by_id"`
}

type BrandModel struct {
	DBPOOL *pgxpool.Pool
}

// CREATE TABLE IF NOT EXISTS brands (
//     id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
//     logo_url text NOT NULL UNIQUE,
//     title varchar(50) NOT NULL UNIQUE,
//     slug varchar(50) NOT NULL UNIQUE,
//     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
//     updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
//     created_by_id uuid NOT NULL,
//     updated_by_id uuid NOT NULL,

//     CHECK (updated_at >= created_at)
// )

// CREATE TABLE IF NOT EXISTS products_brands (
//     id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
//     product_id uuid NOT NULL,
//     brand_id uuid NOT NULL,

//     UNIQUE (product_id, brand_id)
//
