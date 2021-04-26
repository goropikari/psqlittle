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

// NullTestType is Null test type
type NullTestType int

const (
	// EqualNull corresponds to `IS NULL` operation
	EqualNull NullTestType = iota

	// NotEqualNull corresponds to `IS NOT NULL` operation
	NotEqualNull
)

// MathOp express SQL mathemathical operators
type MathOp int

const (
	// EqualOp is equal operator
	EqualOp MathOp = iota
	// ref: translator/expression.go: func (e BinOpNode) Eval()
	// ref: translator/postgres.go: func mathOperator()

	// NotEqualOp is not equal operator
	NotEqualOp

	Plus

	Minus

	Multiply

	Divide

	GT

	LT

	GEQ

	LEQ
)
