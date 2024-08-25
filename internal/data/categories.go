package data

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/validator"
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

func ValidateCategory(v *validator.Validator, category *Category) {
	v.Check(category.NameTk != "", "name_tk", "must be provided")
	v.Check(category.NameRu != "", "name_ru", "must be provided")
	v.Check(category.NameEn != "", "name_en", "must be provided")

	v.Check(len(category.NameEn) <= 500, "name_en", "must be more than 500 bytes long")
	v.Check(len(category.NameEn) <= 500, "name_en", "must be more than 500 bytes long")
	v.Check(len(category.NameEn) <= 500, "name_en", "must be more than 500 bytes long")

	fmt.Println(category.CreatedAt)
	fmt.Println(category.UpdatedAt)

	v.Check(category.UpdatedAt.After(time.Now()), "updated_at", "can not be in the future")
	v.Check(category.CreatedAt.After(category.UpdatedAt), "created_at", "can not be later than updated_at")
}
