package querystr

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/lann/squirrel"
)

type Operator string

const (
	Lt  Operator = "lt"
	Lte Operator = "lte"
	Gt  Operator = "gt"
	Gte Operator = "gte"
	In  Operator = "in"
)

type Order string

const (
	ASC  Order = "asc"
	DESC Order = "desc"
)

const (
	KeyOrder  = "q-order"
	KeyLimit  = "q-limit"
	KeyOffset = "q-offset"
)

func Query(src interface{}, tableName string, values url.Values) (sql string, args []interface{}, err error) {
	builder := squirrel.Select("*").From(tableName)

	// add where clause
	for key, value := range values {
		nBuilder, err := whereValueForKey(builder, src, key, value[0])
		if err != nil {
			return "", []interface{}{}, err
		}
		builder = nBuilder
	}

	// add order if it exists
	if oVal := values.Get(KeyOrder); oVal != "" {
		field, order, err := orderFromValue(src, oVal)
		if err != nil {
			return "", []interface{}{}, err
		}
		builder = builder.OrderBy(field + " " + string(order))
	}

	// add limit and offset
	builder = builder.Limit(uintFromKey(values, KeyLimit, 1000))
	builder = builder.Offset(uintFromKey(values, KeyOffset, 0))

	return builder.ToSql()
}

func uintFromKey(values url.Values, key string, d uint64) uint64 {
	v := values.Get(key)
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil || i < 0 {
		return d
	}
	return uint64(i)
}

func whereValueForKey(builder squirrel.SelectBuilder, src interface{}, key, value string) (squirrel.SelectBuilder, error) {
	keyParts := strings.Split(key, "-")
	sqlString := ""

	switch len(keyParts) {
	case 1:
		if !structHasField(src, keyParts[0]) {
			return builder, nil
		}
		sqlString = keyParts[0] + " = ?"
	case 2:
		op := Operator(keyParts[0])
		field := keyParts[1]
		if !structHasField(src, field) {
			return builder, nil
		}
		switch op {
		case Lt:
			sqlString = field + " < ?"
		case Lte:
			sqlString = field + " <= ?"
		case Gt:
			sqlString = field + " > ?"
		case Gte:
			sqlString = field + " >= ?"
		case In:
			values := strings.Split(value, ",")
			builder = builder.Where(squirrel.Eq{field: values})
			return builder, nil
		default:
			return builder, nil
		}
	default:
		return builder, nil
	}

	builder = builder.Where(sqlString, value)
	return builder, nil
}

func orderFromValue(src interface{}, value string) (string, Order, error) {
	parts := strings.Split(value, "-")

	// return error if not in the format asc-field
	if len(parts) != 2 {
		return "", ASC, errors.New("mware: q-order format invalid.")
	}

	// return error field isn't in struct
	if !structHasField(src, parts[1]) {
		message := fmt.Sprintf("mware: q-order field, %v, isn't present in model", parts[1])
		return "", ASC, errors.New(message)
	}

	// return error if ordering is invalid
	order := Order(parts[0])
	if order != ASC && order != DESC {
		message := fmt.Sprintf("mware: q-order ordering, %v, is invalid.  Must use asc or desc.", parts[0])
		return "", ASC, errors.New(message)
	}

	return parts[1], order, nil
}

func structHasField(src interface{}, key string) bool {
	objT := reflect.TypeOf(src).Elem()
	for i := 0; i < objT.NumField(); i++ {
		field := strings.ToLower(objT.Field(i).Name)
		if field == strings.ToLower(key) {
			return true
		}
	}
	return false
}
