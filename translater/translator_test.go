package translator

import "testing"

func TestOR(t *testing.T) {

	var tests = []struct {
		name     string
		node     Expr
		expected BoolType
	}{
		{
			name: "True or True",
			node: ORNode{
				Lexpr: BoolConstNode{
					Bool: True,
				},
				Rexpr: BoolConstNode{
					Bool: True,
				},
			},
			expected: True,
		},
		{
			name: "True or False",
			node: ORNode{
				Lexpr: BoolConstNode{
					Bool: True,
				},
				Rexpr: BoolConstNode{
					Bool: False,
				},
			},
			expected: True,
		},
		{
			name: "False or True",
			node: ORNode{
				Lexpr: BoolConstNode{
					Bool: False,
				},
				Rexpr: BoolConstNode{
					Bool: True,
				},
			},
			expected: True,
		},
		{
			name: "False or False",
			node: ORNode{
				Lexpr: BoolConstNode{
					Bool: False,
				},
				Rexpr: BoolConstNode{
					Bool: False,
				},
			},
			expected: False,
		},
		{
			name: "True or Null",
			node: ORNode{
				Lexpr: BoolConstNode{
					Bool: True,
				},
				Rexpr: BoolConstNode{
					Bool: Null,
				},
			},
			expected: True,
		},
		{
			name: "Null or True",
			node: ORNode{
				Lexpr: BoolConstNode{
					Bool: Null,
				},
				Rexpr: BoolConstNode{
					Bool: True,
				},
			},
			expected: True,
		},
		{
			name: "False or Null",
			node: ORNode{
				Lexpr: BoolConstNode{
					Bool: False,
				},
				Rexpr: BoolConstNode{
					Bool: Null,
				},
			},
			expected: Null,
		},
		{
			name: "Null or False",
			node: ORNode{
				Lexpr: BoolConstNode{
					Bool: Null,
				},
				Rexpr: BoolConstNode{
					Bool: False,
				},
			},
			expected: Null,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.node.Eval()(1)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestAND(t *testing.T) {

	var tests = []struct {
		name     string
		node     Expr
		expected BoolType
	}{
		{
			name: "True and True",
			node: ANDNode{
				Lexpr: BoolConstNode{
					Bool: True,
				},
				Rexpr: BoolConstNode{
					Bool: True,
				},
			},
			expected: True,
		},
		{
			name: "True and False",
			node: ANDNode{
				Lexpr: BoolConstNode{
					Bool: True,
				},
				Rexpr: BoolConstNode{
					Bool: False,
				},
			},
			expected: False,
		},
		{
			name: "False and True",
			node: ANDNode{
				Lexpr: BoolConstNode{
					Bool: False,
				},
				Rexpr: BoolConstNode{
					Bool: True,
				},
			},
			expected: False,
		},
		{
			name: "False and False",
			node: ANDNode{
				Lexpr: BoolConstNode{
					Bool: False,
				},
				Rexpr: BoolConstNode{
					Bool: False,
				},
			},
			expected: False,
		},
		{
			name: "True and Null",
			node: ANDNode{
				Lexpr: BoolConstNode{
					Bool: True,
				},
				Rexpr: BoolConstNode{
					Bool: Null,
				},
			},
			expected: Null,
		},
		{
			name: "Null or True",
			node: ANDNode{
				Lexpr: BoolConstNode{
					Bool: Null,
				},
				Rexpr: BoolConstNode{
					Bool: True,
				},
			},
			expected: Null,
		},
		{
			name: "False and Null",
			node: ANDNode{
				Lexpr: BoolConstNode{
					Bool: False,
				},
				Rexpr: BoolConstNode{
					Bool: Null,
				},
			},
			expected: False,
		},
		{
			name: "Null and False",
			node: ANDNode{
				Lexpr: BoolConstNode{
					Bool: Null,
				},
				Rexpr: BoolConstNode{
					Bool: False,
				},
			},
			expected: False,
		},
		{
			name: "Null and not False",
			node: ANDNode{
				Lexpr: BoolConstNode{
					Bool: Null,
				},
				Rexpr: NotNode{
					Expr: BoolConstNode{
						Bool: False,
					},
				},
			},
			expected: Null,
		},
		{
			name: "Null and not True",
			node: ANDNode{
				Lexpr: BoolConstNode{
					Bool: Null,
				},
				Rexpr: NotNode{
					Expr: BoolConstNode{
						Bool: True,
					},
				},
			},
			expected: False,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.node.Eval()(1)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestBinOp(t *testing.T) {

	var tests = []struct {
		name     string
		node     Expr
		expected interface{}
	}{
		{
			name: "1 = 1",
			node: BinOpNode{
				Op: Equal,
				Lexpr: IntergerNode{
					Val: 1,
				},
				Rexpr: IntergerNode{
					Val: 1,
				},
			},
			expected: True,
		},
		{
			name: "1 = 2",
			node: BinOpNode{
				Op: Equal,
				Lexpr: IntergerNode{
					Val: 1,
				},
				Rexpr: IntergerNode{
					Val: 2,
				},
			},
			expected: False,
		},
		{
			name: "1 != 2",
			node: BinOpNode{
				Op: NotEqual,
				Lexpr: IntergerNode{
					Val: 1,
				},
				Rexpr: IntergerNode{
					Val: 2,
				},
			},
			expected: True,
		},
		{
			name: "1 = null",
			node: BinOpNode{
				Op: Equal,
				Lexpr: IntergerNode{
					Val: 1,
				},
				Rexpr: BoolConstNode{
					Bool: Null,
				},
			},
			expected: Null,
		},
		{
			name: "1 != null",
			node: BinOpNode{
				Op: NotEqual,
				Lexpr: IntergerNode{
					Val: 1,
				},
				Rexpr: BoolConstNode{
					Bool: Null,
				},
			},
			expected: Null,
		},
		{
			name: "True = True",
			node: BinOpNode{
				Op: Equal,
				Lexpr: BoolConstNode{
					Bool: True,
				},
				Rexpr: BoolConstNode{
					Bool: True,
				},
			},
			expected: True,
		},
		{
			name: "True = False",
			node: BinOpNode{
				Op: Equal,
				Lexpr: BoolConstNode{
					Bool: True,
				},
				Rexpr: BoolConstNode{
					Bool: False,
				},
			},
			expected: False,
		},
		{
			name: "True != True",
			node: BinOpNode{
				Op: NotEqual,
				Lexpr: BoolConstNode{
					Bool: True,
				},
				Rexpr: BoolConstNode{
					Bool: True,
				},
			},
			expected: False,
		},
		{
			name: "True != False",
			node: BinOpNode{
				Op: NotEqual,
				Lexpr: BoolConstNode{
					Bool: True,
				},
				Rexpr: BoolConstNode{
					Bool: False,
				},
			},
			expected: True,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.node.Eval()(1)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
