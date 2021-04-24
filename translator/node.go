package translator

import (
	"github.com/goropikari/mysqlite2/backend"
	"github.com/goropikari/mysqlite2/core"
)

// BoolType express SQL boolean including Null
type BoolType int

const (
	// True is true of BoolType
	True BoolType = iota

	// False is false of BoolType
	False

	// Null is null of BoolType
	Null
)

// MathOp express SQL mathemathical operators
type MathOp int

const (
	// EqualOp is equal operator
	EqualOp MathOp = iota

	// NotEqualOp is not equal operator
	NotEqualOp
)

// RelationalAlgebraNode is interface of RelationalAlgebraNode
type RelationalAlgebraNode interface {
	Eval(backend.DB) backend.Table
}

// WhereNode is Node of where clause
type WhereNode struct {
	Condition WhereExpr
	Table     TableNode
}

// TableNode is Node of table
type TableNode struct {
	TableName string
}

// Eval evaluates TableNode
func (t *TableNode) Eval(db backend.DB) (backend.Table, error) {
	tb, err := db.GetTable(t.TableName)
	if err != nil {
		return nil, err
	}

	return tb, err
}

// Eval evaluate WhereNode
func (wn *WhereNode) Eval(db backend.DB) (backend.Table, error) {
	newTable, err := wn.Table.Eval(db)
	if err != nil {
		return nil, err
	}

	srcRows := newTable.Copy().GetRows()
	condFunc := wn.Condition.Eval()

	rows := make([]backend.Row, 0, len(srcRows))
	for _, row := range srcRows {
		if condFunc(row) == True {
			rows = append(rows, row)
		}
	}

	newTable.SetRows(rows)

	return newTable, nil
}

// WhereExpr is interface of boolean expression
type WhereExpr interface {
	Eval() func(row backend.Row) core.Value
}

// BoolConstNode is expression of boolean const
type BoolConstNode struct {
	Bool core.Value
}

// Eval evaluates BoolConstNode
func (b BoolConstNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		return b.Bool
	}
}

// IntegerNode is expression of integer
type IntegerNode struct {
	Val int
}

// Eval evaluates IntegerNode
func (i IntegerNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		return i.Val
	}
}

// ColRefNode is expression of integer
type ColRefNode struct {
	ColName core.ColName
}

// Eval evaluates ColRefNode
func (i ColRefNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		return row.GetValueByColName(i.ColName)
	}
}

// NotNode is expression of Not
type NotNode struct {
	Expr WhereExpr
}

// Eval evaluates NotNode
func (nn NotNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		return Not(nn.Expr.Eval()(row))
	}
}

// ORNode is expression of OR
type ORNode struct {
	Lexpr WhereExpr
	Rexpr WhereExpr
}

// Eval evaluates ORNode
func (orn ORNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		return OR(orn.Lexpr.Eval()(row), orn.Rexpr.Eval()(row))
	}
}

// ANDNode is expression of AND
type ANDNode struct {
	Lexpr WhereExpr
	Rexpr WhereExpr
}

// Eval evaluates ANDNode
func (andn ANDNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		return AND(andn.Lexpr.Eval()(row), andn.Rexpr.Eval()(row))
	}
}

// BinOpNode is expression of BinOpNode
type BinOpNode struct {
	Op    MathOp
	Lexpr WhereExpr
	Rexpr WhereExpr
}

// Eval evaluates BinOpNode
func (e BinOpNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		l := e.Lexpr.Eval()(row)
		r := e.Rexpr.Eval()(row)
		if l == Null || r == Null {
			return Null
		}

		truth := false
		switch e.Op {
		case EqualOp:
			truth = l == r
		case NotEqualOp:
			truth = l != r
		}

		if truth {
			return True
		}
		return False
	}
}
