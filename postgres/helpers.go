package postgres

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/emanuelquerty/gymulty/domain"
)

func buildUserUpdateQuery(tenantID int, userID int, updates domain.UserUpdate) (string, []any) {
	updatesMap := map[string]any{
		"first_name": updates.FirstName,
		"last_name":  updates.LastName,
		"email":      updates.Email,
		"password":   updates.Password,
		"role":       updates.Role,
	}

	var builder strings.Builder
	var i int
	var columnValues []any
	for colName, colValue := range updatesMap {
		if reflect.ValueOf(colValue).IsNil() { // do not update fields with nil values
			continue
		}

		// we want builder to be i.e, "first_name=$1, role=$2, ..."
		i++
		builder.WriteString(colName)
		builder.WriteString("=$")
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString(", ")
		columnValues = append(columnValues, colValue) // appends the value of each column name
	}

	columnValues = append(columnValues, userID, tenantID)
	cols := builder.String()[:builder.Len()-2] // remove the trailing space and comma

	builder.Reset()
	builder.WriteString(" WHERE id=$")
	builder.WriteString(strconv.Itoa(i + 1))
	builder.WriteString(" AND tenant_id=$")
	builder.WriteString(strconv.Itoa(i + 2))
	builder.WriteString(" RETURNING *")
	cols += builder.String()

	query := "UPDATE users SET " + cols
	return query, columnValues
}
