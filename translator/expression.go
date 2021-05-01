package translator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/goropikari/psqlittle/backend"
	"github.com/goropikari/psqlittle/core"
)

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

// ExpressionNode is interface of boolean expression
type ExpressionNode interface {
	Eval() func(row backend.Row) (core.Value, error)
}

// BoolConstNode is expression of boolean const
type BoolConstNode struct {
	Bool core.BoolType
}

// Eval evaluates BoolConstNode
func (b BoolConstNode) Eval() func(backend.Row) (core.Value, error) {
	return func(row backend.Row) (core.Value, error) {
		return b.Bool, nil
	}
}

// IntegerNode is expression of integer
type IntegerNode struct {
	Val int
}

// Eval evaluates IntegerNode
func (i IntegerNode) Eval() func(backend.Row) (core.Value, error) {
	return func(row backend.Row) (core.Value, error) {
		return i.Val, nil
	}
}

// FloatNode is expression of integer
type FloatNode struct {
	Val float64
}

// Eval evaluates FloatNode
func (f FloatNode) Eval() func(backend.Row) (core.Value, error) {
	return func(row backend.Row) (core.Value, error) {
		return f.Val, nil
	}
}

// StringNode is expression of integer
type StringNode struct {
	Val string
}

// Eval evaluates StringNode
func (s StringNode) Eval() func(backend.Row) (core.Value, error) {
	return func(row backend.Row) (core.Value, error) {
		return s.Val, nil
	}
}

// ColRefNode is expression of integer
type ColRefNode struct {
	ColName core.ColumnName
}

// Eval evaluates ColRefNode
func (n ColRefNode) Eval() func(backend.Row) (core.Value, error) {
	return func(row backend.Row) (core.Value, error) {
		val, err := row.GetValueByColName(n.ColName)
		if err != nil {
			return nil, err
		}
		if val == nil {
			return core.Null, nil
		}
		return val, nil
	}
}

// ColWildcardNode is expression of integer
type ColWildcardNode struct{}

// Eval evaluates ColWildcardNode
func (n ColWildcardNode) Eval() func(backend.Row) (core.Value, error) {
	return func(backend.Row) (core.Value, error) {
		return core.Wildcard, nil
	}
}

// NotNode is expression of Not
type NotNode struct {
	Expr ExpressionNode
}

// Eval evaluates NotNode
func (nn NotNode) Eval() func(backend.Row) (core.Value, error) {
	return func(row backend.Row) (core.Value, error) {
		v, err := nn.Expr.Eval()(row)
		if err != nil {
			return nil, err
		}
		return core.Not(v), nil
	}
}

// ORNode is expression of OR
type ORNode struct {
	Lexpr ExpressionNode
	Rexpr ExpressionNode
}

// Eval evaluates ORNode
func (orn ORNode) Eval() func(backend.Row) (core.Value, error) {
	return func(row backend.Row) (core.Value, error) {
		l, err := orn.Lexpr.Eval()(row)
		if err != nil {
			return nil, err
		}
		r, err := orn.Rexpr.Eval()(row)
		if err != nil {
			return nil, err
		}

		return core.OR(l, r), nil
	}
}

// ANDNode is expression of AND
type ANDNode struct {
	Lexpr ExpressionNode
	Rexpr ExpressionNode
}

// Eval evaluates ANDNode
func (andn ANDNode) Eval() func(backend.Row) (core.Value, error) {
	return func(row backend.Row) (core.Value, error) {
		l, err := andn.Lexpr.Eval()(row)
		if err != nil {
			return nil, err
		}
		r, err := andn.Rexpr.Eval()(row)
		if err != nil {
			return nil, err
		}
		return core.AND(l, r), nil
	}
}

// NullTestNode is expression of `IS (NOT) NULL`
type NullTestNode struct {
	TestType NullTestType
	Expr     ExpressionNode
}

// Eval evaluates NullTestNode
func (n NullTestNode) Eval() func(backend.Row) (core.Value, error) {
	return func(row backend.Row) (core.Value, error) {
		val, err := n.Expr.Eval()(row)
		if err != nil {
			return nil, err
		}
		// val is null
		if n.TestType == EqualNull {
			if val == core.Null {
				return core.True, nil
			}
			return core.False, nil
		}

		// val is not null
		if val == core.Null {
			return core.False, nil
		}
		return core.True, nil
	}
}

// CaseNode is expression of CaseNode
type CaseNode struct {
	CaseWhenExprs   []ExpressionNode
	CaseResultExprs []ExpressionNode
	DefaultResult   ExpressionNode
}

// Eval evaluates CaseNode
func (c *CaseNode) Eval() func(backend.Row) (core.Value, error) {
	return func(row backend.Row) (core.Value, error) {
		for k, expr := range c.CaseWhenExprs {
			v, err := expr.Eval()(row)
			if err != nil {
				return nil, err
			}
			if v == core.True {
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
func (e BinOpNode) Eval() func(backend.Row) (core.Value, error) {
	// ref: translator/const.go: MathOp
	// ref: translator/postgres.go: func mathOperator()
	return func(row backend.Row) (core.Value, error) {
		l, err := e.Lexpr.Eval()(row)
		if err != nil {
			return nil, err
		}
		r, err := e.Rexpr.Eval()(row)
		if err != nil {
			return nil, err
		}
		if l == core.Null || r == core.Null {
			return core.Null, nil
		}

		switch e.Op {
		case EqualOp:
			return toSQLBool(l == r), nil
		case NotEqualOp:
			return toSQLBool(l != r), nil
		case CONCAT:
			lStr := fmt.Sprintf("%v", l)
			rStr := fmt.Sprintf("%v", r)
			return lStr + rStr, nil
		}

		if reflect.ValueOf(l).Kind() == reflect.Int {
			if reflect.ValueOf(r).Kind() == reflect.Int {
				return compIntInt(e.Op, l, r), nil
			}
			return compIntFloat(e.Op, l, r), nil
		}

		if reflect.ValueOf(l).Kind() == reflect.Float64 {
			if reflect.ValueOf(r).Kind() == reflect.Float64 {
				return compFloatFloat(e.Op, l, r), nil
			}
			return compFloatInt(e.Op, l, r), nil
		}
		if reflect.ValueOf(l).Kind() == reflect.String && reflect.ValueOf(r).Kind() == reflect.String {
			return compStrStr(e.Op, l, r), nil
		}

		return core.Null, errors.New("Not Implemented")
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
