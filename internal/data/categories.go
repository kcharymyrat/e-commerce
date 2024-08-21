package data

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID   `json:"id"`
	NameTk    string      `json:"name_tk"`
	NameEn    string      `json:"name_en"`
	NameRu    string      `json:"name_ru"`
	Parent    *Category   `json:"parent,omitempty"`
	Children  []*Category `json:"children,omitempty"`
	CreatedAt time.Time   `json:"created_at,omitempty"`
	UpdatedAt time.Time   `json:"updated_at,omitempty"`
	// CreatedBy *User `json:"-"`
	// UdatedBy  *User `json:"-"`
}
