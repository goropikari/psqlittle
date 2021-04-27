package integration_test

import (
	"fmt"
	"testing"

	"github.com/goropikari/mysqlite2/backend"
	"github.com/goropikari/mysqlite2/core"
	trans "github.com/goropikari/mysqlite2/translator"
	"github.com/stretchr/testify/assert"
)

func TestSelectQuery(t *testing.T) {

	tests := []struct {
		name     string
		query    string
		expected trans.Result
	}{
		{
			name:  "select *",
			query: "select * from hoge",
			expected: &trans.QueryResult{
				Columns: []string{"id", "cid", "name"},
				Records: core.ValuesList{
					{123, 1000, "taro"},
					{456, 500, "hanako"},
					{789, nil, "mike"},
				},
			},
		},
		{
			name:  "rename table",
			query: "select h.id, h.name from hoge as h",
			expected: &trans.QueryResult{
				Columns: []string{"id", "name"},
				Records: core.ValuesList{
					{123, "taro"},
					{456, "hanako"},
					{789, "mike"},
				},
			},
		},
		{
			name:  "specify column name",
			query: "select hoge.name, hoge.id from hoge",
			expected: &trans.QueryResult{
				Columns: []string{"name", "id"},
				Records: core.ValuesList{
					{"taro", 123},
					{"hanako", 456},
					{"mike", 789},
				},
			},
		},
		{
			name:  "wildcard and specify column name",
			query: "select *, hoge.name, hoge.id from hoge",
			expected: &trans.QueryResult{
				Columns: []string{"id", "cid", "name", "name", "id"},
				Records: core.ValuesList{
					{123, 1000, "taro", "taro", 123},
					{456, 500, "hanako", "hanako", 456},
					{789, nil, "mike", "mike", 789},
				},
			},
		},
		{
			name:  "simple where",
			query: "select hoge.name from hoge where hoge.id > 123",
			expected: &trans.QueryResult{
				Columns: []string{"name"},
				Records: core.ValuesList{
					{"hanako"},
					{"mike"},
				},
			},
		},
		{
			name:  "complex condition",
			query: "select hoge.name from hoge where hoge.id > 123 and hoge.cid < 1000 or hoge.name = 'hanako'",
			expected: &trans.QueryResult{
				Columns: []string{"name"},
				Records: core.ValuesList{
					{"hanako"},
				},
			},
		},
		{
			name:  "basic and",
			query: "select true and true, true and false, false and true, false and false",
			expected: &trans.QueryResult{
				Columns: []string{"", "", "", ""},
				Records: core.ValuesList{
					{true, false, false, false},
				},
			},
		},
		{
			name:  "null and",
			query: "select null and null, true and null, null and true, false and null, null and false",
			expected: &trans.QueryResult{
				Columns: []string{"", "", "", "", ""},
				Records: core.ValuesList{
					{nil, nil, nil, false, false},
				},
			},
		},
		{
			name:  "basic or",
			query: "select true or true, true or false, false or true, false or false",
			expected: &trans.QueryResult{
				Columns: []string{"", "", "", ""},
				Records: core.ValuesList{
					{true, true, true, false},
				},
			},
		},
		{
			name:  "null or",
			query: "select null or null, true or null, null or true, false or null, null or false",
			expected: &trans.QueryResult{
				Columns: []string{"", "", "", "", ""},
				Records: core.ValuesList{
					{nil, true, true, nil, nil},
				},
			},
		},
		{
			name:  "null is null",
			query: "select null is null",
			expected: &trans.QueryResult{
				Columns: []string{""},
				Records: core.ValuesList{
					{true},
				},
			},
		},
		{
			name:  "null is not null",
			query: "select null is not null",
			expected: &trans.QueryResult{
				Columns: []string{""},
				Records: core.ValuesList{
					{false},
				},
			},
		},
		{
			name:  "null = null, null != null",
			query: "select null = null, null != null",
			expected: &trans.QueryResult{
				Columns: []string{"", ""},
				Records: core.ValuesList{
					{nil, nil},
				},
			},
		},
		{
			name:  "string concat",
			query: "select 'hoge' || 'piyo'",
			expected: &trans.QueryResult{
				Columns: []string{""},
				Records: core.ValuesList{
					{"hogepiyo"},
				},
			},
		},
		{
			name:  "case expression",
			query: "select case when hoge.name = 'taro' then 'TARO' else 'OTHER' end from hoge",
			expected: &trans.QueryResult{
				Columns: []string{""},
				Records: core.ValuesList{
					{"TARO"},
					{"OTHER"},
					{"OTHER"},
				},
			},
		},
		{
			name:  "int arithmetic",
			query: "select 1+15, 5-3, 2*20, 50/2",
			expected: &trans.QueryResult{
				Columns: []string{"", "", "", ""},
				Records: core.ValuesList{
					{16, 2, 40, 25},
				},
			},
		},
		{
			name:  "float arithmetic",
			query: "select 1.0+1.5, 2.1-0.1, 2.0*20.7, 50.0/4",
			expected: &trans.QueryResult{
				Columns: []string{"", "", "", ""},
				Records: core.ValuesList{
					{2.5, 2.0, 41.4, 12.5},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			db := prepareDB()
			raNode, _ := trans.NewPGTranslator(tt.query).Translate()
			actual, _ := raNode.Eval(db)

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func prepareDB() backend.DB {
	db := backend.NewDatabase()

	query := "create table hoge (id int, cid int, name varchar(255))"
	raNode, _ := trans.NewPGTranslator(query).Translate()
	raNode.Eval(db)
	if _, err := raNode.Eval(db); err != nil {
		fmt.Println("error:", err)
	}

	query = "insert into hoge (name, cid, id) values ('taro', 1000, 123), ('hanako', 500, 456), ('mike', null, 789)"
	raNode, _ = trans.NewPGTranslator(query).Translate()
	raNode.Eval(db)

	return db
}
