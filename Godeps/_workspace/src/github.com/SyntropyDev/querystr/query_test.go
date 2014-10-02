package querystr_test

import (
	"errors"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/SyntropyDev/querystr"
)

const (
	tableName = "test_items"
)

type testStruct struct {
	sql  string
	v    url.Values
	args []interface{}
	err  error
}

type testItem struct {
	Id        int64
	Created   int64
	Updated   int64
	Deleted   bool
	Name      string
	Data      string
	TwoWord   string
	TruckId   int64
	Timestamp int64
}

var (
	testCases = []testStruct{
		{
			sql:  "select * from test_items where name = ? and timestamp > ? limit 1000 offset 0",
			v:    map[string][]string{"name": []string{"hello"}, "gt-timestamp": []string{"2"}},
			args: []interface{}{"hello", "2"},
			err:  nil,
		},
		{
			sql:  "select * from test_items where name in (?,?,?) limit 1000 offset 0",
			v:    map[string][]string{"in-name": []string{"john,brad,tim"}},
			args: []interface{}{"john", "brad", "tim"},
			err:  nil,
		},
		{
			sql:  "select * from test_items limit 1 offset 0",
			v:    map[string][]string{"q-limit": []string{"1"}},
			args: []interface{}{},
			err:  nil,
		},
		{
			sql:  "select * from test_items limit 1000 offset 1",
			v:    map[string][]string{"q-offset": []string{"1"}},
			args: []interface{}{},
			err:  nil,
		},
		{
			sql:  "select * from test_items order by timestamp desc limit 1000 offset 0",
			v:    map[string][]string{"q-order": []string{"desc-timestamp"}},
			args: []interface{}{},
			err:  nil,
		},
		{
			sql:  "",
			v:    map[string][]string{"q-order": []string{"desc-times"}},
			args: []interface{}{},
			err:  errors.New("should throw error"),
		},
		{
			sql:  "select * from test_items limit 1000 offset 0",
			v:    map[string][]string{"madeUp": []string{"should be ignored"}},
			args: []interface{}{},
			err:  nil,
		},
	}
)

func TestQueries(t *testing.T) {
	for _, testCase := range testCases {
		sql, args, err := querystr.Query(&testItem{}, tableName, testCase.v)

		if err == nil && testCase.err != nil || err != nil && testCase.err == nil {
			t.Fatal(err, "not equal", testCase.err)
		}

		if err != nil {
			continue
		}

		if strings.ToLower(sql) != strings.ToLower(testCase.sql) {
			t.Fatal(sql + " should equal " + testCase.sql)
		}

		for i, arg := range args {
			if !reflect.DeepEqual(arg, testCase.args[i]) {
				t.Fatal(arg, "should equal", testCase.args[i])
			}
		}

	}
}
