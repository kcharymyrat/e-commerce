package filters

import (
	"fmt"
	"time"

	"github.com/kcharymyrat/e-commerce/internal/validator"
)

const DEPLOYMENT_TIME_STR = "2024-09-01T00:00:00Z"

type CreatedUpdatedAtFilters struct {
	CreatedAtFrom *time.Time `json:"created_at_from"`
	CreatedAtUpTo *time.Time `json:"created_at_up_to"`
	UpdatedAtFrom *time.Time `json:"updated_at_from"`
	UpdatedAtUpTo *time.Time `json:"updated_at_up_to"`
}

func ValidateCreatedUpdatedAtFilters(v *validator.Validator, f CreatedUpdatedAtFilters) {
	deployment_time, err := time.Parse(time.RFC3339, DEPLOYMENT_TIME_STR)
	if err != nil {
		panic("fix deployment_time, it is constant to be used")
	}
	if f.CreatedAtFrom != nil {
		v.Check(f.CreatedAtFrom.Before(time.Now()), "created_at_from", "can not be in the future")
		if f.CreatedAtUpTo != nil {
			v.Check(f.CreatedAtUpTo.After(deployment_time), "created_up_to", fmt.Sprintf("shall be after %v", deployment_time))
			v.Check(f.CreatedAtUpTo.After(*f.CreatedAtFrom), "created_at_up_to", "can not be before created_at_from")
		}
	}
	if f.UpdatedAtFrom != nil {
		v.Check(f.UpdatedAtFrom.Before(time.Now()), "updated_at_from", "can not be in the future")
		if f.UpdatedAtUpTo != nil {
			v.Check(f.UpdatedAtUpTo.After(deployment_time), "updated_up_to", fmt.Sprintf("shall be after %v", deployment_time))
			v.Check(f.UpdatedAtUpTo.After(*f.UpdatedAtFrom), "updated_at_up_to", "can not be before updated_at_from")
		}

	}
}
