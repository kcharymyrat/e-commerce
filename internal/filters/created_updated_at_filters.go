package filters

import (
	"time"
)

type CreatedUpdatedAtFilter struct {
	CreatedAtFrom *time.Time `json:"created_at_from"`
	CreatedAtUpTo *time.Time `json:"created_at_up_to" validate:"omitempty,gtfield=CreatedAtFrom"`
	UpdatedAtFrom *time.Time `json:"updated_at_from"`
	UpdatedAtUpTo *time.Time `json:"updated_at_up_to" validate:"omitempty,gtfield=CreatedAtFrom"`
}
