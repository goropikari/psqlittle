package translator_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	trans "github.com/goropikari/mysqlite2/translator"
	"github.com/goropikari/mysqlite2/translator/mock"
)

func TestOR(t *testing.T) {

	var tests = []struct {
		name     string
		node     trans.Expr
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

func TestAND(t *testing.T) {

	var tests = []struct {
		name     string
		node     trans.Expr
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

func TestBinOp(t *testing.T) {

	var tests = []struct {
		name     string
		node     trans.Expr
		expected interface{}
	}{
		{
			name: "1 = 1",
			node: trans.BinOpNode{
				Op: trans.Equal,
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
				Op: trans.Equal,
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
				Op: trans.NotEqual,
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
				Op: trans.Equal,
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
				Op: trans.NotEqual,
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
				Op: trans.Equal,
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
				Op: trans.Equal,
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
				Op: trans.NotEqual,
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
				Op: trans.NotEqual,
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
