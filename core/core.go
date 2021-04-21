package core

const NotFound = ColumnID(-1)

// ColName is column name
type ColName string

// ColNames is list of column names
type ColNames []ColName

type Value interface{}
type Values []Value

type Row struct {
	Values Values
}

type Rows []Row

func (r Row) Equal(other Row) bool {
	if other.Values == nil {
		return false
	}
	if len(r.Values) != len(other.Values) {
		return false
	}

	ok := true
	for i := 0; i < len(r.Values); i++ {
		if r.Values[i] != other.Values[i] {
			ok = false
		}
	}

	return ok
}

func (r Rows) Equal(other Rows) bool {
	if other == nil {
		return false
	}
	if len(r) != len(other) {
		return false
	}

	ok := true
	for i := 0; i < len(r); i++ {
		if !r[i].Equal(other[i]) {
			ok = false
		}
	}

	return ok
}

type ColumnID int64

func (r *Row) getByID(i ColumnID) Value {
	return r.Values[i]
}

type Table struct {
	ColNames ColNames
	Rows     Rows
}

type Adder interface {
	Add(v interface{})
}

type ValuesList []Values

type ResVals struct {
	Values Values
}

func (t *Table) Project(names ColNames) Rows {
	returnRows := make(Rows, 0, 10)
	idxs := t.ToIndex(names)

	for _, row := range t.Rows {
		returnRow := Row{}
		for _, i := range idxs {
			returnRow.Values = append(returnRow.Values, row.getByID(i))
		}
		returnRows = append(returnRows, returnRow)
	}

	return returnRows
}

func (t *Table) ToIndex(names ColNames) []ColumnID {
	idxs := make([]ColumnID, 0, 10)
	tbCols := t.ColNames
	for _, name := range names {
		ok := false
		for i, col := range tbCols {
			if name == col {
				idxs = append(idxs, ColumnID(i))
				ok = true
			}
		}
		if !ok {
			idxs = append(idxs, NotFound)
		}
	}

	return idxs
}
