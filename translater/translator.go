package translator

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
	Eval() func(i int) interface{}
}

// BoolConstNode is expression of boolean const
type BoolConstNode struct {
	Bool interface{}
}

// Eval evaluates BoolConstNode
func (b BoolConstNode) Eval() func(int) interface{} {
	return func(x int) interface{} {
		return b.Bool
	}
}

// IntergerNode is expression of integer
type IntergerNode struct {
	Val int
}

// Eval evaluates IntegerNode
func (i IntergerNode) Eval() func(int) interface{} {
	return func(x int) interface{} {
		return i.Val
	}
}

// NotNode is expression of Not
type NotNode struct {
	Expr Expr
}

// Eval evaluates NotNode
func (nn NotNode) Eval() func(int) interface{} {
	return func(i int) interface{} {
		return Not(nn.Expr.Eval()(i))
	}
}

// ORNode is expression of OR
type ORNode struct {
	Lexpr Expr
	Rexpr Expr
}

// Eval evaluates ORNode
func (orn ORNode) Eval() func(int) interface{} {
	return func(i int) interface{} {
		return OR(orn.Lexpr.Eval()(i), orn.Rexpr.Eval()(i))
	}
}

// ANDNode is expression of AND
type ANDNode struct {
	Lexpr Expr
	Rexpr Expr
}

// Eval evaluates ANDNode
func (andn ANDNode) Eval() func(int) interface{} {
	return func(i int) interface{} {
		return AND(andn.Lexpr.Eval()(i), andn.Rexpr.Eval()(i))
	}
}

// BinOpNode is expression of BinOpNode
type BinOpNode struct {
	Op    MathOp
	Lexpr Expr
	Rexpr Expr
}

// Eval evaluates BinOpNode
func (e BinOpNode) Eval() func(int) interface{} {
	return func(i int) interface{} {
		l := e.Lexpr.Eval()(i)
		r := e.Rexpr.Eval()(i)
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
func Not(x interface{}) interface{} {
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
func OR(x, y interface{}) interface{} {
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
func AND(x, y interface{}) interface{} {
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
