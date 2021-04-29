package translator_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/goropikari/mysqlite2/backend"
	"github.com/goropikari/mysqlite2/backend/mock"
	"github.com/goropikari/mysqlite2/core"
	"github.com/goropikari/mysqlite2/testing/fake"
	trans "github.com/goropikari/mysqlite2/translator"
	"github.com/stretchr/testify/assert"
)

func TestORNode(t *testing.T) {

	var tests = []struct {
		name     string
		node     trans.ExpressionNode
		expected core.BoolType
	}{
		{
			name: "True or True",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.True,
				},
			},
			expected: core.True,
		},
		{
			name: "True or False",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.False,
				},
			},
			expected: core.True,
		},
		{
			name: "False or True",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.True,
				},
			},
			expected: core.True,
		},
		{
			name: "False or False",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.False,
				},
			},
			expected: core.False,
		},
		{
			name: "True or Null",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
			},
			expected: core.True,
		},
		{
			name: "Null or True",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.True,
				},
			},
			expected: core.True,
		},
		{
			name: "False or Null",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
			},
			expected: core.Null,
		},
		{
			name: "Null or False",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.False,
				},
			},
			expected: core.Null,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	row := mock.NewMockRow(ctrl)
	row.EXPECT().GetValueByColName(gomock.Any()).Return(gomock.Any()).AnyTimes()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.node.Eval()(row)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestANDNode(t *testing.T) {

	var tests = []struct {
		name     string
		node     trans.ExpressionNode
		expected core.BoolType
	}{
		{
			name: "True and True",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.True,
				},
			},
			expected: core.True,
		},
		{
			name: "True and False",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.False,
				},
			},
			expected: core.False,
		},
		{
			name: "False and True",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.True,
				},
			},
			expected: core.False,
		},
		{
			name: "False and False",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.False,
				},
			},
			expected: core.False,
		},
		{
			name: "True and Null",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
			},
			expected: core.Null,
		},
		{
			name: "Null or True",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.True,
				},
			},
			expected: core.Null,
		},
		{
			name: "False and Null",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
			},
			expected: core.False,
		},
		{
			name: "Null and False",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.False,
				},
			},
			expected: core.False,
		},
		{
			name: "Null and not False",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
				Rexpr: trans.NotNode{
					Expr: trans.BoolConstNode{
						Bool: core.False,
					},
				},
			},
			expected: core.Null,
		},
		{
			name: "Null and not True",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
				Rexpr: trans.NotNode{
					Expr: trans.BoolConstNode{
						Bool: core.True,
					},
				},
			},
			expected: core.False,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	row := mock.NewMockRow(ctrl)
	row.EXPECT().GetValueByColName(gomock.Any()).Return(gomock.Any()).AnyTimes()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.node.Eval()(row)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestNullTestNode(t *testing.T) {

	var tests = []struct {
		name     string
		node     trans.ExpressionNode
		rowRes   core.Value
		expected interface{}
	}{
		{
			name: "Null is Null",
			node: trans.NullTestNode{
				TestType: trans.EqualNull,
				Expr: trans.BoolConstNode{
					Bool: core.Null,
				},
			},
			expected: core.True,
		},
		{
			name: "Null is not Null",
			node: trans.NullTestNode{
				TestType: trans.NotEqualNull,
				Expr: trans.BoolConstNode{
					Bool: core.Null,
				},
			},
			expected: core.False,
		},
		{
			name: "0 = Null",
			node: trans.NullTestNode{
				TestType: trans.EqualNull,
				Expr: trans.IntegerNode{
					Val: 0,
				},
			},
			expected: core.False,
		},
		{
			name: "1 = Null",
			node: trans.NullTestNode{
				TestType: trans.EqualNull,
				Expr: trans.IntegerNode{
					Val: 1,
				},
			},
			expected: core.False,
		},
		{
			name: "2 = Null",
			node: trans.NullTestNode{
				TestType: trans.EqualNull,
				Expr: trans.IntegerNode{
					Val: 2,
				},
			},
			expected: core.False,
		},
		{
			name: "id is null (if id's value is null)",
			node: trans.NullTestNode{
				TestType: trans.EqualNull,
				Expr: trans.ColRefNode{
					ColName: core.ColumnName{
						TableName: "hoge",
						Name:      "id",
					},
				},
			},
			rowRes:   nil,
			expected: core.True,
		},
		{
			name: "id is null (if id's value is 1)",
			node: trans.NullTestNode{
				TestType: trans.EqualNull,
				Expr: trans.ColRefNode{
					ColName: core.ColumnName{
						TableName: "hoge",
						Name:      "id",
					},
				},
			},
			rowRes:   1,
			expected: core.False,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			row := mock.NewMockRow(ctrl)
			row.EXPECT().GetValueByColName(gomock.Any()).Return(tt.rowRes).AnyTimes()

			actual := tt.node.Eval()(row)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestBinOpNode(t *testing.T) {

	var tests = []struct {
		name     string
		node     trans.ExpressionNode
		expected interface{}
	}{
		{
			name: "1 = 1",
			node: trans.BinOpNode{
				Op: trans.EqualOp,
				Lexpr: trans.IntegerNode{
					Val: 1,
				},
				Rexpr: trans.IntegerNode{
					Val: 1,
				},
			},
			expected: core.True,
		},
		{
			name: "1 = 2",
			node: trans.BinOpNode{
				Op: trans.EqualOp,
				Lexpr: trans.IntegerNode{
					Val: 1,
				},
				Rexpr: trans.IntegerNode{
					Val: 2,
				},
			},
			expected: core.False,
		},
		{
			name: "1 != 2",
			node: trans.BinOpNode{
				Op: trans.NotEqualOp,
				Lexpr: trans.IntegerNode{
					Val: 1,
				},
				Rexpr: trans.IntegerNode{
					Val: 2,
				},
			},
			expected: core.True,
		},
		{
			name: "1 = null",
			node: trans.BinOpNode{
				Op: trans.EqualOp,
				Lexpr: trans.IntegerNode{
					Val: 1,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
			},
			expected: core.Null,
		},
		{
			name: "1 != null",
			node: trans.BinOpNode{
				Op: trans.NotEqualOp,
				Lexpr: trans.IntegerNode{
					Val: 1,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.Null,
				},
			},
			expected: core.Null,
		},
		{
			name: "True = True",
			node: trans.BinOpNode{
				Op: trans.EqualOp,
				Lexpr: trans.BoolConstNode{
					Bool: core.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.True,
				},
			},
			expected: core.True,
		},
		{
			name: "True = False",
			node: trans.BinOpNode{
				Op: trans.EqualOp,
				Lexpr: trans.BoolConstNode{
					Bool: core.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.False,
				},
			},
			expected: core.False,
		},
		{
			name: "True != True",
			node: trans.BinOpNode{
				Op: trans.NotEqualOp,
				Lexpr: trans.BoolConstNode{
					Bool: core.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.True,
				},
			},
			expected: core.False,
		},
		{
			name: "True != False",
			node: trans.BinOpNode{
				Op: trans.NotEqualOp,
				Lexpr: trans.BoolConstNode{
					Bool: core.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: core.False,
				},
			},
			expected: core.True,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	row := mock.NewMockRow(ctrl)
	row.EXPECT().GetValueByColName(gomock.Any()).Return(gomock.Any()).AnyTimes()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.node.Eval()(row)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestEvalColRefNode(t *testing.T) {

	n1 := fake.ColName()

	var tests = []struct {
		name      string
		node      trans.ExpressionNode
		givenName core.ColumnName
		expected  interface{}
	}{
		{
			name: "Get row's value",
			node: trans.ColRefNode{
				ColName: n1,
			},
			givenName: n1,
			expected:  fake.Value(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			row := mock.NewMockRow(ctrl)
			row.EXPECT().GetValueByColName(tt.givenName).Return(tt.expected).AnyTimes()

			actual := tt.node.Eval()(row)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

// FIX this where test
// func TestEvalWhereNode(t *testing.T) {
//
// 	cn1 := core.ColumnName{
// 		TableName: "hoge",
// 		Name:      "id",
// 	}
//
// 	var tests = []struct {
// 		name           string
// 		condnode       trans.ExpressionNode
// 		tableName      string
// 		givenName      core.ColumnName
// 		rowRes         []interface{}
// 		expectedRowNum int
// 	}{
// 		{
// 			name: "id = 123",
// 			condnode: trans.BinOpNode{
// 				Op: trans.EqualOp,
// 				Lexpr: trans.ColRefNode{
// 					ColName: cn1,
// 				},
// 				Rexpr: trans.IntegerNode{
// 					Val: 123,
// 				},
// 			},
// 			tableName:      "hoge",
// 			givenName:      cn1,
// 			rowRes:         []interface{}{123, 123, 456},
// 			expectedRowNum: 2,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
//
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			mockRows := []backend.Row{}
// 			for _, v := range tt.rowRes {
// 				row := mock.NewMockRow(ctrl)
// 				row.EXPECT().GetValueByColName(tt.givenName).Return(v).AnyTimes()
// 				mockRows = append(mockRows, row)
// 			}
// 			table := mock.NewMockTable(ctrl)
// 			table.EXPECT().GetRows().Return(mockRows).AnyTimes()
//
// 			count := 0
// 			spyTable := &SpyTable{
// 				Table:       table,
// 				ResultCount: &count,
// 			}
// 			db := mock.NewMockDB(ctrl)
// 			db.EXPECT().GetTable(tt.tableName).Return(spyTable, nil).AnyTimes()
//
// 			whereNode := trans.WhereNode{
// 				Condition: tt.condnode,
// 				Table: &trans.TableNode{
// 					TableName: tt.tableName,
// 				},
// 			}
//
// 			whereNode.Eval(db)
// 		})
// 	}
// }

// TODO: Add test
// ProjectionNode

type SpyTable struct {
	Table       backend.Table
	ResultCount *int
	Values      core.ValuesList
}

func (s *SpyTable) Copy() backend.Table {
	return s
}

func (s *SpyTable) GetName() string {
	return ""
}

func (s *SpyTable) GetRows() []backend.Row {
	return s.Table.GetRows()
}

func (s *SpyTable) GetColNames() core.ColumnNames {
	return s.Table.GetColNames()
}

func (s *SpyTable) GetCols() core.Cols {
	return nil
}

func (s *SpyTable) InsertValues(cs core.ColumnNames, vs core.ValuesList) error {
	return nil
}

func (s *SpyTable) RenameTableName(name string) {}

func (s *SpyTable) Project(cs core.ColumnNames, fns []func(backend.Row) core.Value) (backend.Table, error) {
	return nil, nil
}

func (s *SpyTable) Where(fn func(backend.Row) core.Value) (backend.Table, error) {
	return nil, nil
}

func (s *SpyTable) CrossJoin(backend.Table) (backend.Table, error) {
	return nil, nil
}

func (t *SpyTable) Update(colNames core.ColumnNames, condFn func(backend.Row) core.Value, assignValFns []func(backend.Row) core.Value) (backend.Table, error) {
	return nil, nil
}

func (s *SpyTable) Delete(fn func(backend.Row) core.Value) (backend.Table, error) {
	return nil, nil
}

type SpyRow struct {
	MockRow  backend.Row
	Values   core.Values
	ColNames core.ColumnNames
}

func (r *SpyRow) GetValueByColName(name core.ColumnName) core.Value {
	return r.MockRow.GetValueByColName(name)
}

func (r *SpyRow) GetValues() core.Values {
	return r.MockRow.GetValues()
}

func (r *SpyRow) SetColNames(names core.ColumnNames) {
	r.ColNames = names
}
func (r *SpyRow) UpdateValue(name core.ColumnName, val core.Value) {}
