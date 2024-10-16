package repository

import "fmt"

// TODO: Try to implement this, will field interface{} work with db ?
func argsAppender(
	field interface{},
	query *string,
	db_field_name string,
	argCounter *int,
	args []interface{},
	isFirst bool,
) {
	if field != nil {
		if isFirst {
			*query += fmt.Sprintf(" %s = $%d", db_field_name, argCounter)
		} else {
			*query += fmt.Sprintf(", %s = $%d", db_field_name, argCounter)
		}
		args = append(args, field)
		*argCounter++
	}
}
