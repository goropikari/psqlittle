package core

// ColType is a type of column
type ColType int

const (
	Integer ColType = iota
	VarChar
)

// ColExpr is column name
type ColExpr struct {
	TableName string
	Name      string
}

// Copy copies ColExpr
func (cn ColExpr) Copy() ColExpr {
	return ColExpr{
		TableName: cn.TableName,
		Name:      cn.Name,
	}
}

// ColExprs is list of ColExpr
type ColExprs []ColExpr

// Copy copies ColExprs
func (cn ColExprs) Copy() ColExprs {
	ColExprs := make(ColExprs, 0, len(cn))
	for _, name := range ColExprs {
		ColExprs = append(ColExprs, name.Copy())
	}

	return ColExprs
}

// Equal checks the equality of ColExpr
func (name ColExpr) Equal(other ColExpr) bool {
	return name.TableName == other.TableName && name.Name == other.Name
}

// Value is any type for column
type Value interface{}

// Values is list of Value
type Values []Value

// ValuesList is list of Values
type ValuesList []Values
