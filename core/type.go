package core

// ColType is a type of column
type ColType int

const (
	integer ColType = iota
	varchar
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
