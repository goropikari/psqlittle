package backend

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import (
	"errors"
	"reflect"

	"github.com/goropikari/mysqlite2/core"
)

// DB is interface of DBMS
type DB interface {
	GetTable(string) (Table, error)
}

// Table is interface of table.
type Table interface {
	GetRows() []Row
	Copy() Table
	SetRows([]Row)
}

// Row is interface of row of table.
type Row interface {
	// GetValueByColName is used in ColRefNode when getting value
	GetValueByColName(core.ColName) core.Value
	GetValues() core.Values
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
func (db *Database) CreateTable(tableName string, cols Cols) error {
	if _, ok := db.Tables[tableName]; ok {
		return ErrTableAlreadyExists
	}

	ColNameIndexes := make(ColNameIndexes)
	for k, col := range cols {
		ColNameIndexes[col.ColName] = k
	}

	db.Tables[tableName] = &DBTable{
		Cols:           cols,
		Rows:           make(DBRows, 0),
		ColNameIndexes: ColNameIndexes,
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

// Col is type of column
type Col struct {
	ColName core.ColName
	ColType core.ColType
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

// DBRow is struct of row of table
type DBRow struct {
	Values core.Values
}

// DBRows is list of DBRow
type DBRows []*DBRow

// GetValueByColName gets value from row by ColName
func (r *DBRow) GetValueByColName(name core.ColName) core.Value {
	// TODO: implement
	return 1
}

// GetValues gets values from DBRow
func (r *DBRow) GetValues() core.Values {
	return r.Values
}

// Equal checks the equality of DBRow
func (r *DBRow) Equal(other *DBRow) bool {
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

// Copy copies DBRow
func (r *DBRow) Copy() *DBRow {
	vals := make(core.Values, len(r.Values))
	copy(vals, r.Values)
	return &DBRow{vals}
}

// Equal checks the equality of DBRows
func (r DBRows) Equal(other DBRows) bool {
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

// NotEqual checks the non-equality of DBRows
func (r DBRows) NotEqual(other DBRows) bool {
	return !r.Equal(other)
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
type ColNameIndexes map[core.ColName]int

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

// DBTable is struct for DBTable
type DBTable struct {
	Cols           Cols
	ColNameIndexes ColNameIndexes
	Rows           DBRows
}

// GetRows gets rows from given table
func (t *DBTable) GetRows() []Row {
	rows := make([]Row, 0, len(t.Rows))
	for _, row := range t.Rows {
		rows = append(rows, row)
	}

	return rows
}

// Copy copies DBTable
func (t *DBTable) Copy() Table {
	return &DBTable{
		Cols:           t.Cols.Copy(),
		ColNameIndexes: t.ColNameIndexes.Copy(),
		Rows:           t.Rows.Copy(),
	}
}

// SetRows replate rows
func (t *DBTable) SetRows(rows []Row) {
	dbRows := make([]*DBRow, 0, len(rows))
	for _, row := range rows {
		dbRows = append(dbRows, row.(*DBRow))
	}

	t.Rows = dbRows
}

// Equal checks the equality of DBTable
func (t DBTable) Equal(other DBTable) bool {
	return t.Cols.Equal(other.Cols) && t.Rows.Equal(other.Rows) && t.ColNameIndexes.Equal(other.ColNameIndexes)
}

// NotEqual checks the non-equality of DBTable
func (t DBTable) NotEqual(other DBTable) bool {
	return !t.Equal(other)
}

// Project is method to select columns of table.
func (t *DBTable) Project(names core.ColNames) (DBRows, error) {
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

func (t *DBTable) toIndex(names core.ColNames) ([]ColumnID, error) {
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
func (t *DBTable) Insert(targetCols Cols, inputValsList core.ValuesList) error {
	if targetCols == nil {
		targetCols = t.Cols
	}

	if err := t.validateInsert(targetCols, inputValsList); err != nil {
		return err
	}

	numCols := len(t.Cols)
	colNames := make(core.ColNames, 0, numCols)
	for _, col := range targetCols {
		colNames = append(colNames, col.ColName)
	}
	idxs, err := t.toIndex(colNames)
	if err != nil {
		return err
	}

	rows := make([]*DBRow, 0, len(inputValsList))
	for _, vals := range inputValsList {
		tvalues := make(core.Values, numCols)
		for vi := range idxs {
			tvalues[vi] = vals[vi]
		}
		rows = append(rows, &DBRow{Values: tvalues})
	}

	t.Rows = rows

	return nil
}

func (t *DBTable) validateInsert(cols Cols, valuesList core.ValuesList) error {
	// TODO: valuesList の各要素の長さが全部同じかチェックする
	for _, vals := range valuesList {
		if len(t.Cols) != len(vals) {
			return errors.New("invalid insert elements")
		}
	}

	// TODO: 型で validation かける

	return nil
}