package query

import (
	"fmt"

	"github.com/go-pg/pg/v10/orm"
)

func RowToJson(q *orm.Query, table, asColumn string) {
	q.ColumnExpr(fmt.Sprintf("row_to_json(%s) as %s", table, asColumn))
}

func OrderBy(q *orm.Query, ods []string) {
	for _, o := range ods {
		order := o
		if o[0] == '-' && len(o) > 1 {
			order = fmt.Sprintf("%s DESC", o[1:])
		}
		q.Order(order)
	}
}
