package core

// ColType is a type of column
type ColType int

const (
	Integer ColType = iota
	VarChar
)

// ColName is column name
type ColName struct {
	TableName string
	Name      string
}

// ColNames is list of ColName
type ColNames []ColName

// Equal checks the equality of ColName
func (name ColName) Equal(other ColName) bool {
	return name.TableName == other.TableName && name.Name == other.Name
}

// Value is any type for column
type Value interface{}

// Values is list of Value
type Values []Value

// ValuesList is list of Values
type ValuesList []Values
