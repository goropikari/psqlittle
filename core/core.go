package core

import (
	"errors"
	"reflect"
)

const NotFound = ColumnID(-1)

var (
	TableAlreadyExistsError = errors.New("The table already exists")
)

// ColName is column name
type ColName struct {
	TableName string
	Name      string
}

// ColNames is list of column names
type ColNames []ColName

func (name ColName) Equal(other ColName) bool {
	return name.TableName == other.TableName && name.Name == other.Name
}

func (names ColNames) Equal(others ColNames) bool {
	for k, name := range names {
		if !name.Equal(others[k]) {
			return false
		}
	}

	return true
}

type Value interface{}
type Values []Value
type ValuesList []Values

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

type ColType string
type ColTypes []ColType

type TableSchema struct {
	ColNames ColNames
	ColTypes ColTypes
}

func (schema TableSchema) Equal(other TableSchema) bool {
	return reflect.DeepEqual(schema.ColNames, other.ColNames) && reflect.DeepEqual(schema.ColTypes, other.ColTypes)
}

type Table struct {
	ColNames ColNames
	Rows     Rows
	Schema   TableSchema
}

func (t Table) Equal(other Table) bool {
	return t.ColNames.Equal(other.ColNames) && t.Rows.Equal(other.Rows) && t.Schema.Equal(other.Schema)
}

func (t *Table) validateInsert(cols ColNames, valuesList ValuesList) error {
	// TODO: valuesList の各要素の長さが全部同じかチェックする
	if len(t.ColNames) != len(valuesList[0]) {
		return errors.New("invalid insert elements")
	}

	// TODO: 型で validation かける

	return nil
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
			if name.Equal(col) {
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

type DB struct {
	Tables map[string]Table
}

func NewDB() *DB {
	db := &DB{
		Tables: make(map[string]Table),
	}
	return db
}

func (db *DB) CreateTable(tableName string, schema TableSchema) error {
	if _, ok := db.Tables[tableName]; ok {
		return TableAlreadyExistsError
	}

	db.Tables[tableName] = Table{
		ColNames: make(ColNames, 0),
		Rows:     make(Rows, 0),
		Schema:   schema,
	}
	return nil
}

func (table *Table) Insert(cols ColNames, inputValsList ValuesList) error {
	if cols == nil {
		cols = table.ColNames
	}

	if err := table.validateInsert(cols, inputValsList); err != nil {
		return err
	}

	numCols := len(table.ColNames)
	idxs := table.ToIndex(cols)

	for _, vals := range inputValsList {
		tvalues := make(Values, numCols)
		for vi := range idxs {
			tvalues[vi] = vals[vi]
		}

	}

	return nil
}
