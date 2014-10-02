package sqlutil

import (
	"github.com/coopernurse/gorp"
	"github.com/lann/squirrel"
)

func Select(s gorp.SqlExecutor, builder squirrel.SelectBuilder, src interface{}) error {
	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	_, err = s.Select(src, sql, args...)
	return err
}

func SelectOne(s gorp.SqlExecutor, builder squirrel.SelectBuilder, src interface{}) error {
	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	return s.SelectOne(src, sql, args...)
}

func SelectOneRelation(s gorp.SqlExecutor, tableName string, id interface{}, src interface{}) error {
	query := squirrel.Select("*").From(tableName).Where(squirrel.Eq{"ID": id})
	return SelectOne(s, query, src)
}
