package core

import (
	"errors"
)

var (
	// ErrTableAlreadyExists occures when creating table exists
	ErrTableAlreadyExists = errors.New("the table already exists")

	// ErrIndexNotFound occurs when a table doesn't contain given column.
	ErrIndexNotFound = errors.New("there is no index corresponding column name")
)

// DB is struct for DB
type DB struct {
	Tables map[string]Table
}

// NewDB is constructor of DB
func NewDB() *DB {
	db := &DB{
		Tables: make(map[string]Table),
	}
	return db
}

// CreateTable is method to create table
func (db *DB) CreateTable(tableName string, Cols Cols) error {
	if _, ok := db.Tables[tableName]; ok {
		return ErrTableAlreadyExists
	}

	db.Tables[tableName] = Table{
		Cols: Cols,
		Rows: make(Rows, 0),
	}
	return nil
}

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

// Col is type of column
type Col struct {
	ColName ColName
	ColType ColType
}

// Cols is list of column names
type Cols []Col

// Equal check the equality of Col
func (col Col) Equal(other Col) bool {
	return col.ColName.Equal(other.ColName) && col.ColType == other.ColType
}

// Equal checks the equality of ColName
func (name ColName) Equal(other ColName) bool {
	return name.TableName == other.TableName && name.Name == other.Name
}

// Equal checks the equality of Cols
func (names Cols) Equal(others Cols) bool {
	for k, name := range names {
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

// Row is struct of row of table
type Row struct {
	Values Values
}

// Rows is list of Row
type Rows []Row

// Equal checks the equality of Row
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

// Equal checks the equality of Rows
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

// ColumnID is type of column id (index of column).
type ColumnID int64

// getByID is method to get column value by ColumnID
func (r *Row) getByID(i ColumnID) Value {
	return r.Values[i]
}

// Table is struct for Table
type Table struct {
	Cols           Cols
	ColNameIndexes map[ColName]int64
	Rows           Rows
}

// Equal checks the equality of Table
func (t Table) Equal(other Table) bool {
	return t.Cols.Equal(other.Cols) && t.Rows.Equal(other.Rows)
}

// Project is method to select columns of table.
func (t *Table) Project(names Cols) (Rows, error) {
	returnRows := make(Rows, 0, 10)
	idxs, err := t.toIndex(names)
	if err != nil {
		return nil, err
	}

	for _, row := range t.Rows {
		returnRow := Row{}
		for _, i := range idxs {
			returnRow.Values = append(returnRow.Values, row.getByID(i))
		}
		returnRows = append(returnRows, returnRow)
	}

	return returnRows, nil
}

func (t *Table) toIndex(names Cols) ([]ColumnID, error) {
	idxs := make([]ColumnID, 0, 10)
	tbCols := t.Cols
	for _, name := range names {
		ok := false
		for i, col := range tbCols {
			if name.Equal(col) {
				idxs = append(idxs, ColumnID(i))
				ok = true
			}
		}

		if !ok {
			return nil, ErrIndexNotFound
		}
	}

	return idxs, nil
}

// Insert is method to insert record into table.
func (t *Table) Insert(cols Cols, inputValsList ValuesList) error {
	if cols == nil {
		cols = t.Cols
	}

	if err := t.validateInsert(cols, inputValsList); err != nil {
		return err
	}

	numCols := len(t.Cols)
	idxs, err := t.toIndex(cols)
	if err != nil {
		return err
	}

	for _, vals := range inputValsList {
		tvalues := make(Values, numCols)
		for vi := range idxs {
			tvalues[vi] = vals[vi]
		}
		t.Rows = append(t.Rows, Row{Values: tvalues})
	}

	return nil
}

func (t *Table) validateInsert(cols Cols, valuesList ValuesList) error {
	// TODO: valuesList の各要素の長さが全部同じかチェックする
	if len(t.Cols) != len(valuesList[0]) {
		return errors.New("invalid insert elements")
	}

	// TODO: 型で validation かける

	return nil
}
