package filters

import (
	"github.com/google/uuid"
)

type CreatedUpdatedByFilters struct {
	CreatedByIDs []uuid.UUID `json:"created_by_ids" validate:"omitempty,dive,uuid"`
	UpdatedByIDs []uuid.UUID `json:"updated_by_ids" validate:"omitempty,dive,uuid"`
}
