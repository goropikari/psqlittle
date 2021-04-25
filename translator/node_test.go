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
		expected trans.BoolType
	}{
		{
			name: "True or True",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
			},
			expected: trans.True,
		},
		{
			name: "True or False",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
			},
			expected: trans.True,
		},
		{
			name: "False or True",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
			},
			expected: trans.True,
		},
		{
			name: "False or False",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
			},
			expected: trans.False,
		},
		{
			name: "True or Null",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
			},
			expected: trans.True,
		},
		{
			name: "Null or True",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
			},
			expected: trans.True,
		},
		{
			name: "False or Null",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
			},
			expected: trans.Null,
		},
		{
			name: "Null or False",
			node: trans.ORNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
			},
			expected: trans.Null,
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
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestANDNode(t *testing.T) {

	var tests = []struct {
		name     string
		node     trans.ExpressionNode
		expected trans.BoolType
	}{
		{
			name: "True and True",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
			},
			expected: trans.True,
		},
		{
			name: "True and False",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
			},
			expected: trans.False,
		},
		{
			name: "False and True",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
			},
			expected: trans.False,
		},
		{
			name: "False and False",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
			},
			expected: trans.False,
		},
		{
			name: "True and Null",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
			},
			expected: trans.Null,
		},
		{
			name: "Null or True",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
			},
			expected: trans.Null,
		},
		{
			name: "False and Null",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
			},
			expected: trans.False,
		},
		{
			name: "Null and False",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
			},
			expected: trans.False,
		},
		{
			name: "Null and not False",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
				Rexpr: trans.NotNode{
					Expr: trans.BoolConstNode{
						Bool: trans.False,
					},
				},
			},
			expected: trans.Null,
		},
		{
			name: "Null and not True",
			node: trans.ANDNode{
				Lexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
				Rexpr: trans.NotNode{
					Expr: trans.BoolConstNode{
						Bool: trans.True,
					},
				},
			},
			expected: trans.False,
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
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
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
					Bool: trans.Null,
				},
			},
			expected: trans.True,
		},
		{
			name: "Null is not Null",
			node: trans.NullTestNode{
				TestType: trans.NotEqualNull,
				Expr: trans.BoolConstNode{
					Bool: trans.Null,
				},
			},
			expected: trans.False,
		},
		{
			name: "0 = Null",
			node: trans.NullTestNode{
				TestType: trans.EqualNull,
				Expr: trans.IntegerNode{
					Val: 0,
				},
			},
			expected: trans.False,
		},
		{
			name: "1 = Null",
			node: trans.NullTestNode{
				TestType: trans.EqualNull,
				Expr: trans.IntegerNode{
					Val: 1,
				},
			},
			expected: trans.False,
		},
		{
			name: "2 = Null",
			node: trans.NullTestNode{
				TestType: trans.EqualNull,
				Expr: trans.IntegerNode{
					Val: 2,
				},
			},
			expected: trans.False,
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
			expected: trans.True,
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
			expected: trans.False,
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
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
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
			expected: trans.True,
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
			expected: trans.False,
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
			expected: trans.True,
		},
		{
			name: "1 = null",
			node: trans.BinOpNode{
				Op: trans.EqualOp,
				Lexpr: trans.IntegerNode{
					Val: 1,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
			},
			expected: trans.Null,
		},
		{
			name: "1 != null",
			node: trans.BinOpNode{
				Op: trans.NotEqualOp,
				Lexpr: trans.IntegerNode{
					Val: 1,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.Null,
				},
			},
			expected: trans.Null,
		},
		{
			name: "True = True",
			node: trans.BinOpNode{
				Op: trans.EqualOp,
				Lexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
			},
			expected: trans.True,
		},
		{
			name: "True = False",
			node: trans.BinOpNode{
				Op: trans.EqualOp,
				Lexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
			},
			expected: trans.False,
		},
		{
			name: "True != True",
			node: trans.BinOpNode{
				Op: trans.NotEqualOp,
				Lexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
			},
			expected: trans.False,
		},
		{
			name: "True != False",
			node: trans.BinOpNode{
				Op: trans.NotEqualOp,
				Lexpr: trans.BoolConstNode{
					Bool: trans.True,
				},
				Rexpr: trans.BoolConstNode{
					Bool: trans.False,
				},
			},
			expected: trans.True,
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
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
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
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestEvalWhereNode(t *testing.T) {

	cn1 := core.ColumnName{
		TableName: "hoge",
		Name:      "id",
	}

	var tests = []struct {
		name           string
		condnode       trans.ExpressionNode
		tableName      string
		givenName      core.ColumnName
		rowRes         []interface{}
		expectedRowNum int
	}{
		{
			name: "id = 123",
			condnode: trans.BinOpNode{
				Op: trans.EqualOp,
				Lexpr: trans.ColRefNode{
					ColName: cn1,
				},
				Rexpr: trans.IntegerNode{
					Val: 123,
				},
			},
			tableName:      "hoge",
			givenName:      cn1,
			rowRes:         []interface{}{123, 123, 456},
			expectedRowNum: 2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRows := []backend.Row{}
			for _, v := range tt.rowRes {
				row := mock.NewMockRow(ctrl)
				row.EXPECT().GetValueByColName(tt.givenName).Return(v).AnyTimes()
				mockRows = append(mockRows, row)
			}
			table := mock.NewMockTable(ctrl)
			table.EXPECT().GetRows().Return(mockRows).AnyTimes()

			count := 0
			spyTable := &SpyTable{
				Table:       table,
				ResultCount: &count,
			}
			db := mock.NewMockDB(ctrl)
			db.EXPECT().GetTable(tt.tableName).Return(spyTable, nil).AnyTimes()

			whereNode := trans.WhereNode{
				Condition: tt.condnode,
				Table: &trans.TableNode{
					TableName: tt.tableName,
				},
			}

			whereNode.Eval(db)
			if count != tt.expectedRowNum {
				t.Errorf("expected %v, actual %v", tt.expectedRowNum, count)
			}
		})
	}
}

func TestEvalProjectionNode(t *testing.T) {

	cn1 := core.ColumnName{
		TableName: "hoge",
		Name:      "id",
	}
	cn2 := core.ColumnName{
		TableName: "hoge",
		Name:      "name",
	}

	const someConst = 10

	var tests = []struct {
		name               string
		targetCols         core.ColumnNames
		resTargets         []trans.ExpressionNode
		table              trans.RelationalAlgebraNode
		tableName          string
		tableColNames      core.ColumnNames
		rowRes             core.ValuesList
		expectedRowNum     int
		expectedValuesList core.ValuesList
	}{
		{
			name:       "select id",
			targetCols: core.ColumnNames{cn2, cn1, core.ColumnName{}},
			table: &trans.TableNode{
				TableName: "hoge",
			},
			tableName: "hoge",
			resTargets: []trans.ExpressionNode{
				trans.ColRefNode{cn2},
				trans.ColRefNode{cn1},
				trans.IntegerNode{Val: someConst},
			},
			tableColNames: core.ColumnNames{cn1, cn2, core.ColumnName{}},
			rowRes: core.ValuesList{ // row mock response
				{123, "foo", someConst},
				{456, "bar", someConst},
				{789, "baz", someConst},
			},
			expectedRowNum: 3,
			expectedValuesList: core.ValuesList{
				{"foo", 123, someConst},
				{"bar", 456, someConst},
				{"baz", 789, someConst},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRows := []backend.Row{}
			for _, vals := range tt.rowRes {
				row := mock.NewMockRow(ctrl)
				for k, col := range tt.tableColNames {
					row.EXPECT().GetValueByColName(col).Return(vals[k]).AnyTimes()
				}
				row.EXPECT().SetValues(gomock.Any()).AnyTimes()
				mockRows = append(mockRows, &SpyRow{MockRow: row})
			}
			table := mock.NewMockTable(ctrl)
			table.EXPECT().GetRows().Return(mockRows).AnyTimes()

			count := 0
			spyTable := &SpyTable{
				Table:       table,
				ResultCount: &count,
			}
			db := mock.NewMockDB(ctrl)
			db.EXPECT().GetTable(tt.tableName).Return(spyTable, nil).AnyTimes()

			projectNode := trans.ProjectionNode{
				ResTargets:     tt.resTargets,
				TargetColNames: tt.targetCols,
				Table:          tt.table,
			}

			projectNode.Eval(db)

			if count != tt.expectedRowNum {
				t.Errorf("expected %v, actual %v", tt.expectedRowNum, count)
			}

			resValsList := make(core.ValuesList, 0, len(mockRows))
			for _, row := range mockRows {
				resValsList = append(resValsList, row.(*SpyRow).Values)
			}
			assert.Equal(t, tt.expectedValuesList, resValsList)
		})
	}
}

type SpyTable struct {
	Table       backend.Table
	ResultCount *int
	Values      core.ValuesList
}

func (s *SpyTable) Copy() backend.Table {
	return s
}

func (s *SpyTable) GetRows() []backend.Row {
	return s.Table.GetRows()
}

func (s *SpyTable) SetRows(rows []backend.Row) {
	*s.ResultCount = len(rows)
}

func (s *SpyTable) GetColNames() core.ColumnNames {
	return s.Table.GetColNames()
}

func (s *SpyTable) SetColNames(names core.ColumnNames) {
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

func (r *SpyRow) SetValues(values core.Values) {
	r.Values = values
}

func (r *SpyRow) SetColNames(names core.ColumnNames) {
	r.ColNames = names
}
