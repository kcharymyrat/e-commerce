package filters

import (
	"fmt"
	"time"
)

type CreatedUpdatedAtFilter struct {
	CreatedAtFrom *time.Time `json:"created_at_from"`
	CreatedAtUpTo *time.Time `json:"created_at_up_to" validate:"omitempty,gtfield=CreatedAtFrom"`
	UpdatedAtFrom *time.Time `json:"updated_at_from"`
	UpdatedAtUpTo *time.Time `json:"updated_at_up_to" validate:"omitempty,gtfield=CreatedAtFrom"`
}

func AddCreatedUpdateAtFilterToSQL(
	f *CreatedUpdatedAtFilter, query *string, argCounter *int, args []interface{},
) {
	if f.CreatedAtFrom != nil {
		*query += fmt.Sprintf(" AND created_at >= $%d", *argCounter)
		args = append(args, *f.CreatedAtFrom)
		*argCounter++
	}

	if f.CreatedAtUpTo != nil {
		*query += fmt.Sprintf(" AND created_at <= $%d", *argCounter)
		args = append(args, *f.CreatedAtUpTo)
		*argCounter++
	}

	if f.UpdatedAtFrom != nil {
		*query += fmt.Sprintf(" AND updated_at >= $%d", *argCounter)
		args = append(args, *f.UpdatedAtFrom)
		*argCounter++
	}

	if f.UpdatedAtUpTo != nil {
		*query += fmt.Sprintf(" AND updated_at <= $%d", *argCounter)
		args = append(args, *f.UpdatedAtUpTo)
		*argCounter++
	}
}
