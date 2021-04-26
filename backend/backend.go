package backend

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import (
	"errors"

	"github.com/goropikari/mysqlite2/core"
)

// DB is interface of DBMS
type DB interface {
	GetTable(string) (Table, error)
	CreateTable(string, core.Cols) error
}

// Table is interface of table.
type Table interface {
	Copy() Table
	GetColNames() core.ColumnNames
	SetColNames(core.ColumnNames)
	GetRows() []Row
	SetRows([]Row)
	InsertValues(core.ColumnNames, core.ValuesList) error
	// GetCols() []Col
	// SetCols([]Col)
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

// DBRow is struct of row of table
type DBRow struct {
	ColNames core.ColumnNames
	Values   core.Values
}

// DBRows is list of DBRow
type DBRows []*DBRow

// GetValueByColName gets value from row by ColName
func (r *DBRow) GetValueByColName(name core.ColumnName) core.Value {
	for k, v := range r.ColNames {
		if v == name {
			return r.Values[k]
		}
	}
	return nil
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

// SetRows replate rows
func (t *DBTable) SetRows(rows []Row) {
	dbRows := make([]*DBRow, 0, len(rows))
	for _, row := range rows {
		dbRows = append(dbRows, row.(*DBRow))
	}

	t.Rows = dbRows
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

// Project is method to select columns of table.
func (t *DBTable) Project(names core.ColumnNames) (DBRows, error) {
	returnRows := make(DBRows, 0, 10)
	idxs, err := t.toIndex(names)
	if err != nil {
		return nil, err
	}

	for _, row := range t.Rows {
		returnRow := &DBRow{}
		for _, i := range idxs {
			returnRow.Values = append(returnRow.Values, row.getByID(i))
		}
		returnRows = append(returnRows, returnRow)
	}

	return returnRows, nil
}

// Rename renames table name
func (t *DBTable) Rename(tableName string) {
	for i := 0; i < len(t.Cols); i++ {
		col := t.Cols[i]
		col.ColName.TableName = tableName
		t.Cols[i] = col
	}
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
