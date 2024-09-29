package data

import (
	"time"

	"github.com/google/uuid"
)

type Language struct {
	ID          uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	Code        string    `json:"code" db:"code" validate:"required,max=10"`
	Name        string    `json:"name" db:"name" validate:"required,max=50"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" validate:"required,gtefield=CreatedAt"`
	CreatedById uuid.UUID `json:"created_by_id" db:"created_by_id" validate:"required,uuid"`
	UpdatedById uuid.UUID `json:"updated_by_id" db:"updated_by_id" validate:"required,uuid"`
	Version     int       `json:"version" db:"version" validate:"required,number,min=1"`
}

// CREATE TABLE IF NOT EXISTS languages (
//     id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
//     code varchar(10) NOT NULL UNIQUE,
//     name varchar(50) NOT NULL,
//     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
//     updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
//     created_by_id uuid NOT NULL,
//     updated_by_id uuid NOT NULL,
//     version integer NOT NULL DEFAULT 1,

//     CHECK (updated_at >= created_at)
// );
