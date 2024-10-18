package filters

import (
	"fmt"

	"github.com/google/uuid"
)

type CreatedUpdatedByFilter struct {
	CreatedByIDs []uuid.UUID `json:"created_by_ids" validate:"omitempty,dive,uuid"`
	UpdatedByIDs []uuid.UUID `json:"updated_by_ids" validate:"omitempty,dive,uuid"`
}

func AddCreatedUpdateByFilterToSQL(
	f *CreatedUpdatedByFilter, query *string, argCounter *int, args []interface{},
) {
	if len(f.CreatedByIDs) > 0 {
		*query += fmt.Sprintf(" AND created_by = ANY($%d)", *argCounter)
		args = append(args, f.CreatedByIDs)
		*argCounter++
	}

	if len(f.UpdatedByIDs) > 0 {
		*query += fmt.Sprintf(" AND updated_by_id = ANY($%d)", *argCounter)
		args = append(args, f.UpdatedByIDs)
		*argCounter++
	}
}
