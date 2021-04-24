package translator

import (
	"github.com/goropikari/mysqlite2/backend"
	"github.com/goropikari/mysqlite2/core"
)

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

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
		val := row.GetValueByColName(i.ColName)
		if val != nil {
			return val
		}
		return Null
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

// NullTestNode is expression of `IS (NOT) NULL`
type NullTestNode struct {
	TestType NullTestType
	Expr     WhereExpr
}

// Eval evaluates NullTestNode
func (n NullTestNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		val := n.Expr.Eval()(row)
		truth := False
		if val == Null {
			truth = True
		}
		if n.TestType == EqualNull {
			return truth
		}
		return Not(truth)
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
