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
				Table: &trans.WhereNode{
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
				Table: &trans.WhereNode{
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
