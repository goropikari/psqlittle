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

// Copy copies ColName
func (cn ColName) Copy() ColName {
	return ColName{
		TableName: cn.TableName,
		Name:      cn.Name,
	}
}

// ColNames is list of ColName
type ColNames []ColName

// Copy copies ColNames
func (cn ColNames) Copy() ColNames {
	colNames := make(ColNames, 0, len(cn))
	for _, name := range colNames {
		colNames = append(colNames, name.Copy())
	}

	return colNames
}

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
