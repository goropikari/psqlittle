package core

// ColType is a type of column
type ColType int

const (
	Integer ColType = iota
	VarChar
)

// ColumnName is column name
type ColumnName struct {
	TableName string
	Name      string
}

// Copy copies ColumnName
func (cn ColumnName) Copy() ColumnName {
	return ColumnName{
		TableName: cn.TableName,
		Name:      cn.Name,
	}
}

// ColumnNames is list of ColumnName
type ColumnNames []ColumnName

// Copy copies ColumnNames
func (cn ColumnNames) Copy() ColumnNames {
	ColumnNames := make(ColumnNames, 0, len(cn))
	for _, name := range ColumnNames {
		ColumnNames = append(ColumnNames, name.Copy())
	}

	return ColumnNames
}

// Equal checks the equality of ColumnName
func (name ColumnName) Equal(other ColumnName) bool {
	return name.TableName == other.TableName && name.Name == other.Name
}

// Value is any type for column
type Value interface{}

// Values is list of Value
type Values []Value

// ValuesList is list of Values
type ValuesList []Values
