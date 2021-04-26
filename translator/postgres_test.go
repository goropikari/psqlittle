package translator_test

import (
	"testing"

	"github.com/goropikari/mysqlite2/core"
	trans "github.com/goropikari/mysqlite2/translator"
	"github.com/stretchr/testify/assert"
)

func TestTranslateSelect(t *testing.T) {
	var tests = []struct {
		name     string
		expected trans.RelationalAlgebraNode
		query    string
	}{
		{
			name: "test translator",
			expected: &trans.ProjectionNode{
				TargetColNames: core.ColumnNames{
					{TableName: "foo", Name: "id"},
					{TableName: "foo", Name: "name"},
				},
				ResTargets: []trans.ExpressionNode{
					trans.ColRefNode{core.ColumnName{TableName: "foo", Name: "id"}},
					trans.ColRefNode{core.ColumnName{TableName: "foo", Name: "name"}},
				},
				RANode: &trans.WhereNode{
					Condition: nil,
					Table: &trans.TableNode{
						TableName: "foo",
					},
				},
			},
			query: "SELECT foo.id, foo.name FROM foo",
		},
		{
			name: "test wildcard",
			expected: &trans.ProjectionNode{
				TargetColNames: core.ColumnNames{
					core.ColumnName{},
				},
				ResTargets: []trans.ExpressionNode{
					trans.ColWildcardNode{},
				},
				RANode: &trans.WhereNode{
					Condition: nil,
					Table: &trans.TableNode{
						TableName: "foo",
					},
				},
			},
			query: "SELECT * FROM foo",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			transl := trans.NewPGTranslator(tt.query)
			actual, _ := transl.Translate()

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestTranslateCreate(t *testing.T) {
	var tests = []struct {
		name      string
		tableName string
		expected  trans.RelationalAlgebraNode
		query     string
	}{
		{
			name:      "test translator",
			tableName: "foo",
			expected: &trans.CreateTableNode{
				TableName: "foo",
				ColumnDefs: core.Cols{
					core.Col{
						ColName: core.ColumnName{TableName: "foo", Name: "id"},
						ColType: core.Integer,
					},
					core.Col{
						ColName: core.ColumnName{TableName: "foo", Name: "name"},
						ColType: core.VarChar,
					},
				},
			},
			query: "CREATE TABLE foo (id int, name varchar(255))",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			transl := trans.NewPGTranslator(tt.query)
			actual, _ := transl.Translate()

			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestTranslateInsert(t *testing.T) {
	var tests = []struct {
		name      string
		tableName string
		expected  trans.RelationalAlgebraNode
		query     string
	}{
		{
			name:      "test insert",
			tableName: "foo",
			expected: &trans.InsertNode{
				TableName:   "foo",
				ColumnNames: core.ColumnNames{},
				ValuesList: core.ValuesList{
					core.Values{1, "mike"},
				},
			},
			query: "INSERT INTO foo values (1, 'mike')",
		},
		{
			name:      "test insert multi values",
			tableName: "foo",
			expected: &trans.InsertNode{
				TableName: "foo",
				ColumnNames: core.ColumnNames{
					{
						TableName: "foo",
						Name:      "id",
					},
					{
						TableName: "foo",
						Name:      "name",
					},
				},
				ValuesList: core.ValuesList{
					core.Values{1, "mike"},
					core.Values{100, "taro"},
				},
			},
			query: "INSERT INTO foo (id, name) values (1, 'mike'), (100, 'taro')",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			transl := trans.NewPGTranslator(tt.query)
			actual, _ := transl.Translate()

			assert.Equal(t, tt.expected, actual)
		})
	}
}
