package translator

import (
	"fmt"
	"reflect"

	"github.com/goropikari/mysqlite2/backend"
	"github.com/goropikari/mysqlite2/core"
)

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

// ExpressionNode is interface of boolean expression
type ExpressionNode interface {
	Eval() func(row backend.Row) core.Value
}

// BoolConstNode is expression of boolean const
type BoolConstNode struct {
	Bool core.BoolType
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

// FloatNode is expression of integer
type FloatNode struct {
	Val float64
}

// Eval evaluates FloatNode
func (f FloatNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		return f.Val
	}
}

// StringNode is expression of integer
type StringNode struct {
	Val string
}

// Eval evaluates StringNode
func (s StringNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		return s.Val
	}
}

// ColRefNode is expression of integer
type ColRefNode struct {
	ColName core.ColumnName
}

// Eval evaluates ColRefNode
func (n ColRefNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		val := row.GetValueByColName(n.ColName)
		if val == nil {
			return core.Null
		}
		return val
	}
}

// ColWildcardNode is expression of integer
type ColWildcardNode struct{}

// Eval evaluates ColWildcardNode
func (n ColWildcardNode) Eval() func(backend.Row) core.Value {
	return func(backend.Row) core.Value {
		return core.Wildcard
	}
}

// NotNode is expression of Not
type NotNode struct {
	Expr ExpressionNode
}

// Eval evaluates NotNode
func (nn NotNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		return core.Not(nn.Expr.Eval()(row))
	}
}

// ORNode is expression of OR
type ORNode struct {
	Lexpr ExpressionNode
	Rexpr ExpressionNode
}

// Eval evaluates ORNode
func (orn ORNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		return core.OR(orn.Lexpr.Eval()(row), orn.Rexpr.Eval()(row))
	}
}

// ANDNode is expression of AND
type ANDNode struct {
	Lexpr ExpressionNode
	Rexpr ExpressionNode
}

// Eval evaluates ANDNode
func (andn ANDNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		return core.AND(andn.Lexpr.Eval()(row), andn.Rexpr.Eval()(row))
	}
}

// NullTestNode is expression of `IS (NOT) NULL`
type NullTestNode struct {
	TestType NullTestType
	Expr     ExpressionNode
}

// Eval evaluates NullTestNode
func (n NullTestNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		val := n.Expr.Eval()(row)
		truth := core.False
		if val == core.Null {
			truth = core.True
		}
		if n.TestType == EqualNull {
			return truth
		}
		return core.Not(truth)
	}
}

// CaseNode is expression of CaseNode
type CaseNode struct {
	CaseWhenExprs   []ExpressionNode
	CaseResultExprs []ExpressionNode
	DefaultResult   ExpressionNode
}

// Eval evaluates CaseNode
func (c *CaseNode) Eval() func(backend.Row) core.Value {
	return func(row backend.Row) core.Value {
		for k, expr := range c.CaseWhenExprs {
			if expr.Eval()(row) == core.True {
				return c.CaseResultExprs[k].Eval()(row)
			}
		}

		return c.DefaultResult.Eval()(row)
	}
}

// BinOpNode is expression of BinOpNode
type BinOpNode struct {
	Op    MathOp
	Lexpr ExpressionNode
	Rexpr ExpressionNode
}

// Eval evaluates BinOpNode
func (e BinOpNode) Eval() func(backend.Row) core.Value {
	// ref: translator/const.go: MathOp
	// ref: translator/postgres.go: func mathOperator()
	return func(row backend.Row) core.Value {
		l := e.Lexpr.Eval()(row)
		r := e.Rexpr.Eval()(row)
		if l == core.Null || r == core.Null {
			return core.Null
		}

		switch e.Op {
		case EqualOp:
			return toSQLBool(l == r)
		case NotEqualOp:
			return toSQLBool(l != r)
		case CONCAT:
			lStr := fmt.Sprintf("%v", l)
			rStr := fmt.Sprintf("%v", r)
			return lStr + rStr
		}

		if reflect.ValueOf(l).Kind() == reflect.Int {
			if reflect.ValueOf(r).Kind() == reflect.Int {
				return compIntInt(e.Op, l, r)
			}
			return compIntFloat(e.Op, l, r)
		}

		if reflect.ValueOf(l).Kind() == reflect.Float64 {
			if reflect.ValueOf(r).Kind() == reflect.Float64 {
				return compFloatFloat(e.Op, l, r)
			}
			return compFloatInt(e.Op, l, r)
		}
		if reflect.ValueOf(l).Kind() == reflect.String && reflect.ValueOf(r).Kind() == reflect.String {
			return compStrStr(e.Op, l, r)
		}

		fmt.Println("Not Implemented")
		return core.Null
	}
}

func compIntInt(op MathOp, l core.Value, r core.Value) core.Value {
	switch op {
	case Plus:
		return l.(int) + r.(int)
	case Minus:
		return l.(int) - r.(int)
	case Multiply:
		return l.(int) * r.(int)
	case Divide:
		return l.(int) / r.(int)
	case GT:
		return toSQLBool(l.(int) > r.(int))
	case LT:
		return toSQLBool(l.(int) < r.(int))
	case GEQ:
		return toSQLBool(l.(int) >= r.(int))
	case LEQ:
		return toSQLBool(l.(int) <= r.(int))
	}

	return nil
}

func compIntFloat(op MathOp, l core.Value, r core.Value) core.Value {
	switch op {
	case Plus:
		return float64(l.(int)) + r.(float64)
	case Minus:
		return float64(l.(int)) - r.(float64)
	case Multiply:
		return float64(l.(int)) * r.(float64)
	case Divide:
		return float64(l.(int)) / r.(float64)
	case GT:
		fmt.Println(float64(l.(int)), r.(float64), float64(l.(int)) > r.(float64))
		return toSQLBool(float64(l.(int)) > r.(float64))
	case LT:
		return toSQLBool(float64(l.(int)) < r.(float64))
	case GEQ:
		return toSQLBool(float64(l.(int)) >= r.(float64))
	case LEQ:
		return toSQLBool(float64(l.(int)) <= r.(float64))
	}

	return nil
}

func compFloatInt(op MathOp, l core.Value, r core.Value) core.Value {
	switch op {
	case Plus:
		return l.(float64) + float64(r.(int))
	case Minus:
		return l.(float64) - float64(r.(int))
	case Multiply:
		return l.(float64) * float64(r.(int))
	case Divide:
		return l.(float64) / float64(r.(int))
	case GT:
		return toSQLBool(l.(float64) > float64(r.(int)))
	case LT:
		return toSQLBool(l.(float64) < float64(r.(int)))
	case GEQ:
		return toSQLBool(l.(float64) >= float64(r.(int)))
	case LEQ:
		return toSQLBool(l.(float64) <= float64(r.(int)))
	}

	return nil
}

func compFloatFloat(op MathOp, l core.Value, r core.Value) core.Value {
	switch op {
	case Plus:
		return l.(float64) + r.(float64)
	case Minus:
		return l.(float64) - r.(float64)
	case Multiply:
		return l.(float64) * r.(float64)
	case Divide:
		return l.(float64) / r.(float64)
	case GT:
		return toSQLBool(l.(float64) > r.(float64))
	case LT:
		return toSQLBool(l.(float64) < r.(float64))
	case GEQ:
		return toSQLBool(l.(float64) >= r.(float64))
	case LEQ:
		return toSQLBool(l.(float64) <= r.(float64))
	}

	return nil
}

func compStrStr(op MathOp, l core.Value, r core.Value) core.Value {
	switch op {
	case Plus:
		return l.(string) + r.(string)
	case GT:
		return toSQLBool(l.(string) > r.(string))
	case LT:
		return toSQLBool(l.(string) < r.(string))
	case GEQ:
		return toSQLBool(l.(string) >= r.(string))
	case LEQ:
		return toSQLBool(l.(string) <= r.(string))
	}

	return nil
}

func toSQLBool(b bool) core.BoolType {
	if b {
		return core.True
	}
	return core.False
}
