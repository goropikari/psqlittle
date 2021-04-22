package core

import (
	"errors"
	"reflect"
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

	ColNameIndexes := make(ColNameIndexes)
	for k, col := range Cols {
		ColNameIndexes[col.ColName] = k
	}

	db.Tables[tableName] = Table{
		Cols:           Cols,
		Rows:           make(Rows, 0),
		ColNameIndexes: ColNameIndexes,
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

// ColNames is list of ColName
type ColNames []ColName

// Equal checks the equality of ColName
func (name ColName) Equal(other ColName) bool {
	return name.TableName == other.TableName && name.Name == other.Name
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

// Copy copies Row
func (r Row) Copy() Row {
	vals := make(Values, len(r.Values))
	copy(vals, r.Values)
	return Row{vals}
}

// Equal checks the equality of Rows
func (r Rows) Equal(other Rows) bool {
	if len(r) == len(other) && len(r) == 0 {
		return true
	}
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

// NotEqual checks the non-equality of Rows
func (r Rows) NotEqual(other Rows) bool {
	return !r.Equal(other)
}

// Copy copies Rows
func (r Rows) Copy() Rows {
	rows := make(Rows, len(r))
	for k, row := range r {
		rows[k] = row.Copy()
	}

	return rows
}

// ColumnID is type of column id (index of column).
type ColumnID int

// getByID is method to get column value by ColumnID
func (r *Row) getByID(i ColumnID) Value {
	return r.Values[i]
}

// ColNameIndexes is map ColName to corresponding column index
type ColNameIndexes map[ColName]int

// Equal checks the equality of ColNameIndexes
func (c ColNameIndexes) Equal(other ColNameIndexes) bool {
	return reflect.DeepEqual(c, other)
}

// NotEqual checks the non-equality of ColNameIndexes
func (c ColNameIndexes) NotEqual(other ColNameIndexes) bool {
	return !c.Equal(other)
}

// Copy copies ColNameIndexes
func (c ColNameIndexes) Copy() ColNameIndexes {
	indexes := make(ColNameIndexes)

	for key, val := range c {
		indexes[key] = val
	}
	return indexes
}

// Table is struct for Table
type Table struct {
	Cols           Cols
	ColNameIndexes ColNameIndexes
	Rows           Rows
}

// Equal checks the equality of Table
func (t Table) Equal(other Table) bool {
	return t.Cols.Equal(other.Cols) && t.Rows.Equal(other.Rows) && t.ColNameIndexes.Equal(other.ColNameIndexes)
}

// NotEqual checks the non-equality of Table
func (t Table) NotEqual(other Table) bool {
	return !t.Equal(other)
}

// Project is method to select columns of table.
func (t *Table) Project(names ColNames) (Rows, error) {
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

// Rename renames table name
func (t *Table) Rename(tableName string) {
	for i := 0; i < len(t.Cols); i++ {
		col := t.Cols[i]
		col.ColName.TableName = tableName
		t.Cols[i] = col
	}
}

// Copy copies Table
func (t *Table) Copy() *Table {
	return &Table{
		Cols:           t.Cols.Copy(),
		ColNameIndexes: t.ColNameIndexes.Copy(),
		Rows:           t.Rows.Copy(),
	}
}

func (t *Table) toIndex(names ColNames) ([]ColumnID, error) {
	idxs := make([]ColumnID, 0, len(names))
	for _, name := range names {
		if val, ok := t.ColNameIndexes[name]; ok {
			idxs = append(idxs, ColumnID(val))
		} else {
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
	colNames := make(ColNames, 0, numCols)
	for _, col := range cols {
		colNames = append(colNames, col.ColName)
	}
	idxs, err := t.toIndex(colNames)
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
