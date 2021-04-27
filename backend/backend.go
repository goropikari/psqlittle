package backend

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import (
	"errors"
	"fmt"

	"github.com/goropikari/mysqlite2/core"
)

// DB is interface of DBMS
type DB interface {
	GetTable(string) (Table, error)
	CreateTable(string, core.Cols) error
	DropTable(string) error
}

// Table is interface of table.
type Table interface {
	Copy() Table
	GetColNames() core.ColumnNames
	GetRows() []Row
	InsertValues(core.ColumnNames, core.ValuesList) error
	UpdateTableName(string)
	Project(core.ColumnNames, []func(Row) core.Value) (Table, error)
	Where(func(Row) core.Value) (Table, error)
	Update(core.ColumnNames, func(Row) core.Value, []func(Row) core.Value) (Table, error)
	Delete(func(Row) core.Value) (Table, error)
}

// Row is interface of row of table.
type Row interface {
	// GetValueByColName is used in ColRefNode when getting value
	GetValueByColName(core.ColumnName) core.Value
	GetValues() core.Values
	SetValues(core.Values)
	SetColNames(core.ColumnNames)
	UpdateValue(core.ColumnName, core.Value)
}

// Database is struct for Database
type Database struct {
	Tables map[string]*DBTable
}

// NewDatabase is constructor of Database
func NewDatabase() *Database {
	return &Database{
		Tables: make(map[string]*DBTable),
	}
}

// CreateTable is method to create table
func (db *Database) CreateTable(tableName string, cols core.Cols) error {
	if _, ok := db.Tables[tableName]; ok {
		return ErrTableAlreadyExists
	}

	colNames := make(core.ColumnNames, 0, len(cols))
	for _, col := range cols {
		colNames = append(colNames, col.ColName)
	}

	db.Tables[tableName] = &DBTable{
		ColNames: colNames,
		Cols:     cols,
		Rows:     make(DBRows, 0),
	}
	return nil
}

// GetTable gets table from DB
func (db *Database) GetTable(tableName string) (Table, error) {
	if _, ok := db.Tables[tableName]; !ok {
		return nil, ErrTableNotFound
	}

	tb := db.Tables[tableName]
	return tb, nil
}

// DropTable drop table from DB
func (db *Database) DropTable(tableName string) error {
	if _, ok := db.Tables[tableName]; ok {
		delete(db.Tables, tableName)
		return nil
	}
	return ErrTableNotFound
}

// DBRow is struct of row of table
type DBRow struct {
	ColNames core.ColumnNames
	Values   core.Values
}

// DBRows is list of DBRow
type DBRows []*DBRow

type ErrColumnNotFound int

const (
	ColumnNotFound ErrColumnNotFound = iota
)

// GetValueByColName gets value from row by ColName
func (r *DBRow) GetValueByColName(name core.ColumnName) core.Value {
	for k, v := range r.ColNames {
		if v == name {
			return r.Values[k]
		}
	}
	return ColumnNotFound
}

// GetValues gets values from DBRow
func (r *DBRow) GetValues() core.Values {
	return r.Values
}

// SetValues sets vals into row
func (r *DBRow) SetValues(vals core.Values) {
	r.Values = vals
}

// SetColNames sets column names into row
func (r *DBRow) SetColNames(names core.ColumnNames) {
	r.ColNames = names
}

// UpdateValue updates value by specifing column name
func (r *DBRow) UpdateValue(name core.ColumnName, val core.Value) {
	for k, colName := range r.ColNames {
		// fmt.Println("colName:", colName, "givenName:", name)
		if colName == name {
			r.Values[k] = val
		}
	}
}

// Copy copies DBRow
func (r *DBRow) Copy() *DBRow {
	vals := make(core.Values, len(r.Values))
	copy(vals, r.Values)
	names := make(core.ColumnNames, len(r.ColNames))
	copy(names, r.ColNames)
	return &DBRow{
		ColNames: names,
		Values:   vals,
	}
}

// Copy copies DBRows
func (r DBRows) Copy() DBRows {
	rows := make(DBRows, len(r))
	for k, row := range r {
		rows[k] = row.Copy()
	}

	return rows
}

// ColumnID is type of column id (index of column).
type ColumnID int

// getByID is method to get column value by ColumnID
func (r *DBRow) getByID(i ColumnID) core.Value {
	return r.Values[i]
}

// ColNameIndexes is map ColName to corresponding column index
type ColNameIndexes map[core.ColumnName]int

// Copy copies ColNameIndexes
func (c ColNameIndexes) Copy() ColNameIndexes {
	indexes := make(ColNameIndexes)

	for key, val := range c {
		indexes[key] = val
	}
	return indexes
}

// DBTable is struct for DBTable
type DBTable struct {
	ColNames core.ColumnNames
	Cols     core.Cols
	Rows     DBRows
}

// Copy copies DBTable
func (t *DBTable) Copy() Table {
	tb := &DBTable{
		ColNames: t.ColNames.Copy(),
		Cols:     t.Cols.Copy(),
		Rows:     t.Rows.Copy(),
	}
	return tb
}

// GetColNames return column names of table
func (t *DBTable) GetColNames() core.ColumnNames {
	return t.ColNames
}

// SetColNames sets ColNames in Table
func (t *DBTable) SetColNames(names core.ColumnNames) {
	t.ColNames = names
}

// GetRows gets rows from given table
func (t *DBTable) GetRows() []Row {
	// ref: https://stackoverflow.com/a/12994852
	rows := make([]Row, 0, len(t.Rows))
	for _, row := range t.Rows {
		rows = append(rows, row)
	}

	return rows
}

// InsertValues inserts values into the table
func (t *DBTable) InsertValues(names core.ColumnNames, valsList core.ValuesList) error {
	if len(names) == 0 {
		names = t.GetColNames()
	}
	colNames := t.GetColNames()

	err := t.validateInsert(names, valsList)
	if err != nil {
		return err
	}

	numCols := len(colNames)
	indexes := make([]int, 0)
	for _, name := range names {
		for k, v := range colNames {
			if name == v {
				indexes = append(indexes, k)
			}
		}
	}

	for _, vals := range valsList {
		row := &DBRow{ColNames: colNames, Values: make(core.Values, numCols)}
		for vi, ci := range indexes {
			row.Values[ci] = vals[vi]
		}
		t.Rows = append(t.Rows, row)
	}

	return nil
}

func (t *DBTable) validateInsert(names core.ColumnNames, valuesList core.ValuesList) error {
	for _, vals := range valuesList {
		if len(names) != len(vals) {
			return errors.New("invalid insert elements")
		}
	}

	// TODO: 型で validation かける

	return nil
}

// UpdateTableName updates table name
func (t *DBTable) UpdateTableName(name string) {
	for i := 0; i < len(t.ColNames); i++ {
		t.ColNames[i].TableName = name
	}
	for i := 0; i < len(t.Cols); i++ {
		t.Cols[i].ColName.TableName = name
	}

	for i := 0; i < len(t.Rows); i++ {
		for j := 0; j < len(t.Rows[i].ColNames); j++ {
			t.Rows[i].ColNames[j].TableName = name
		}
	}
}

// Project is method to select columns of table.
func (t *DBTable) Project(TargetColNames core.ColumnNames, resFuncs []func(Row) core.Value) (Table, error) {
	rows := t.GetRows()
	newRows := make(DBRows, 0, len(rows))
	for _, row := range t.Rows {
		colNames := make(core.ColumnNames, 0)
		vals := make(core.Values, 0)
		for k, fn := range resFuncs {
			if v := fn(row); v != core.Wildcard {
				if v == ColumnNotFound {
					return nil, fmt.Errorf("column %v is not found", TargetColNames[k])
				}
				vals = append(vals, v)
				colNames = append(colNames, TargetColNames[k])
			} else { // column wildcard
				// Add values
				for _, val := range row.GetValues() {
					if val == nil {
						// Fix me: nil should be converted
						// when the value is inserted.
						vals = append(vals, core.Null)
					} else {
						vals = append(vals, val)
					}
				}

				// Add columns
				for _, name := range t.GetColNames() {
					colNames = append(colNames, name)
				}
			}
		}
		row.SetValues(vals)
		row.SetColNames(colNames)
		newRows = append(newRows, row)
	}

	t.Rows = newRows
	t.SetColNames(TargetColNames)
	// TODO: implement SetCols if type validation is implemented
	// newTable.SetCols(cols)

	return t, nil
}

// Where filters rows by given where conditions
func (t *DBTable) Where(condFn func(Row) core.Value) (Table, error) {
	srcRows := t.Rows
	rows := make([]*DBRow, 0)
	for _, row := range srcRows {
		if condFn(row) == core.True {
			rows = append(rows, row)
		}
	}
	t.Rows = rows

	return t, nil
}

// Update updates records
func (t *DBTable) Update(colNames core.ColumnNames, condFn func(Row) core.Value, assignValFns []func(Row) core.Value) (Table, error) {
	rows := t.Rows
	for _, row := range rows {
		if condFn(row) == core.True {
			for k, name := range colNames {
				row.UpdateValue(name, assignValFns[k](row))
			}
		}
	}

	return nil, nil
}

func (t *DBTable) Delete(condFn func(Row) core.Value) (Table, error) {
	updatedRows := make([]*DBRow, 0)
	for _, row := range t.Rows {
		if condFn(row) == core.True {
			continue
		} else {
			updatedRows = append(updatedRows, row)
		}
	}

	t.Rows = updatedRows
	return nil, nil
}

func (t *DBTable) toIndex(names core.ColumnNames) ([]ColumnID, error) {
	idxs := make([]ColumnID, 0, len(names))
	rawNames := t.GetColNames()
	for _, name := range names {
		for k, rawName := range rawNames {
			if name.Equal(rawName) {
				idxs = append(idxs, ColumnID(k))
			} else {
				return nil, ErrIndexNotFound
			}
		}
	}

	return idxs, nil
}
