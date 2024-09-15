package filters

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/validator"
)

type CreatedUpdatedByFilters struct {
	CreatedByIDs []uuid.UUID `json:"created_by_ids"`
	UpdatedByIDs []uuid.UUID `json:"updated_by_ids"`
}

func ValidateCreatedUpdatedByFilters(v *validator.Validator, f CreatedUpdatedByFilters) {

}
