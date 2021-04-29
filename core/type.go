package core

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

type WildcardType int

const (
	Wildcard WildcardType = iota
)

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

func (cn ColumnName) String() string {
	if cn.TableName == "" {
		return cn.Name
	}
	return cn.TableName + "." + cn.Name
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
	names := make(ColumnNames, 0, len(cn))
	for _, name := range cn {
		names = append(names, name.Copy())
	}

	return names
}

// Equal checks the equality of ColumnName
func (cn ColumnName) Equal(other ColumnName) bool {
	return cn.TableName == other.TableName && cn.Name == other.Name
}

// Equal checks the equality of ColumnNames
func (cn ColumnNames) Equal(others ColumnNames) bool {
	if len(cn) != len(others) {
		return false
	}

	if len(cn) == len(others) && len(cn) == 0 {
		return true
	}

	for k, name := range cn {
		if !name.Equal(others[k]) {
			return false
		}
	}

	return true
}

// Value is any type for column
type Value interface{}

// Values is list of Value
type Values []Value

// ValuesList is list of Values
type ValuesList []Values

// Col is type of column
type Col struct {
	ColName ColumnName
	ColType ColType
}

// Cols is list of Col
type Cols []Col

// Equal check the equality of Col
func (col Col) Equal(other Col) bool {
	return col.ColName.Equal(other.ColName) && col.ColType == other.ColType
}

// Equal checks the equality of Cols
func (cols Cols) Equal(others Cols) bool {
	for k, col := range cols {
		if !col.Equal(others[k]) {
			return false
		}
	}

	return true
}

// NotEqual checks the non-equality of Cols
func (cols Cols) NotEqual(others Cols) bool {
	return !cols.Equal(others)
}

// Copy copies Col.
func (col Col) Copy() Col {
	return Col{col.ColName, col.ColType}
}

// Copy copies Cols.
func (cols Cols) Copy() Cols {
	newCols := make(Cols, 0, len(cols))
	for _, col := range cols {
		newCols = append(newCols, col.Copy())
	}
	return newCols
}
