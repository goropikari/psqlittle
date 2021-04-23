package translator

import "github.com/goropikari/mysqlite2/core"

// Row is interface of row of table.
type Row interface {
	GetValueByColName(core.ColName) core.Value
}

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
	Equal MathOp = iota
	NotEqual
)

// Expr is interface of boolean expression
type Expr interface {
	Eval() func(row Row) core.Value
}

// BoolConstNode is expression of boolean const
type BoolConstNode struct {
	Bool core.Value
}

// Eval evaluates BoolConstNode
func (b BoolConstNode) Eval() func(Row) core.Value {
	return func(row Row) core.Value {
		return b.Bool
	}
}

// IntegerNode is expression of integer
type IntegerNode struct {
	Val int
}

// Eval evaluates IntegerNode
func (i IntegerNode) Eval() func(Row) core.Value {
	return func(row Row) core.Value {
		return i.Val
	}
}

// ColRefNode is expression of integer
type ColRefNode struct {
	ColName core.ColName
}

// Eval evaluates ColRefNode
func (i ColRefNode) Eval() func(Row) core.Value {
	return func(row Row) core.Value {
		return row.GetValueByColName(i.ColName)
	}
}

// NotNode is expression of Not
type NotNode struct {
	Expr Expr
}

// Eval evaluates NotNode
func (nn NotNode) Eval() func(Row) core.Value {
	return func(row Row) core.Value {
		return Not(nn.Expr.Eval()(row))
	}
}

// ORNode is expression of OR
type ORNode struct {
	Lexpr Expr
	Rexpr Expr
}

// Eval evaluates ORNode
func (orn ORNode) Eval() func(Row) core.Value {
	return func(row Row) core.Value {
		return OR(orn.Lexpr.Eval()(row), orn.Rexpr.Eval()(row))
	}
}

// ANDNode is expression of AND
type ANDNode struct {
	Lexpr Expr
	Rexpr Expr
}

// Eval evaluates ANDNode
func (andn ANDNode) Eval() func(Row) core.Value {
	return func(row Row) core.Value {
		return AND(andn.Lexpr.Eval()(row), andn.Rexpr.Eval()(row))
	}
}

// BinOpNode is expression of BinOpNode
type BinOpNode struct {
	Op    MathOp
	Lexpr Expr
	Rexpr Expr
}

// Eval evaluates BinOpNode
func (e BinOpNode) Eval() func(Row) core.Value {
	return func(row Row) core.Value {
		l := e.Lexpr.Eval()(row)
		r := e.Rexpr.Eval()(row)
		if l == Null || r == Null {
			return Null
		}

		truth := false
		switch e.Op {
		case Equal:
			truth = l == r
		case NotEqual:
			truth = l != r
		}

		if truth {
			return True
		}
		return False
	}
}

// Not negates x
func Not(x core.Value) core.Value {
	if x == Null {
		return Null
	}
	if x == True {
		return False
	}
	if x == False {
		return True
	}

	// Shouldn't reach this
	return Null
}

// OR calculates x OR y
func OR(x, y core.Value) core.Value {
	// memo:
	// True or Null -> True
	// Null or True -> True
	// False or Null -> Null
	// Null or False -> Null
	if x == True || y == True {
		return True
	}
	if x == Null || y == Null {
		return Null
	}
	if x == False && y == False {
		return False
	}
	return True
}

// AND calculates x AND y
func AND(x, y core.Value) core.Value {
	// memo:
	// True and Null -> Null
	// Null and True -> Null
	// False and Null -> False
	// Null and False -> False
	if x == True && y == True {
		return True
	}
	if x == False || y == False {
		return False
	}
	return Null
}
