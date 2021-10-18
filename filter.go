package query

import (
	"fmt"
	"reflect"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/shopspring/decimal"
	"github.com/zikosw/bun-query/internal"
)

func IsEmpty(v interface{}) (bool, error) {
	if v == nil {
		return true, nil
	}

	switch vt := v.(type) {
	case int:
		return vt == 0, nil
	case int64:
		return vt == 0, nil
	case float32:
		return vt == 0, nil
	case float64:
		return vt == 0, nil
	case string:
		return vt == "", nil
	case decimal.Decimal:
		return vt.IsZero(), nil
	case time.Time:
		return vt.IsZero(), nil
	default:
		return false, fmt.Errorf("unknown typed")
	}
}

// where =,>=,<=,like. limit, offset
func Filter(q *orm.Query, opts interface{}) error {
	sqlTag := "sql"
	pgOpTag := "pg_op"
	limitOp := "limit"
	offsetOp := "offset"
	likeOp := "like"

	typ := reflect.TypeOf(opts) //.Elem()
	val := reflect.ValueOf(opts)

	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("binding element must be a struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		field := val.Field(i)

		var col string
		if col = typeField.Tag.Get(sqlTag); col == "" {
			col = internal.Underscore(typeField.Name)
		}
		if col == "-" {
			continue
		}

		operator := "="
		if tag := typeField.Tag.Get(pgOpTag); tag != "" {
			operator = tag
		}

		if operator == limitOp {
			q.Limit(int(field.Int()))
			continue
		}
		if operator == offsetOp {
			q.Offset(int(field.Int()))
			continue
		}

		v := field.Interface()
		if empty, err := IsEmpty(v); err != nil {
			return err
		} else if !empty {
			if operator == likeOp {
				// q.Where(fmt.Sprintf("%s LIKE '%%'||?||'%%'", col), v)
				q.Where(fmt.Sprintf("%s LIKE '%%%s%%'", col, v))
			} else {
				q.Where(fmt.Sprintf("%s%s?", col, operator), v)
			}
			fmt.Println("filter,", col, operator, v)
		}
	}
	return nil
}
