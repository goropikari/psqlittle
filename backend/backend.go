package backend

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import (
	"errors"
	"fmt"
	"sort"

	"github.com/goropikari/psqlittle/core"
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
	GetName() string
	GetColNames() core.ColumnNames
	GetRows() []Row
	GetCols() core.Cols
	InsertValues(core.ColumnNames, core.ValuesList) error
	RenameTableName(string)
	Project(core.ColumnNames, []func(Row) (core.Value, error)) (Table, error)
	Where(func(Row) (core.Value, error)) (Table, error)
	CrossJoin(Table) (Table, error)
	OrderBy(core.ColumnNames, []int) (Table, error)
	Limit(int) (Table, error)
	Update(core.ColumnNames, func(Row) (core.Value, error), []func(Row) (core.Value, error)) (Table, error)
	Delete(func(Row) (core.Value, error)) (Table, error)
}

// Row is interface of row of table.
type Row interface {
	// GetValueByColName is used in ColRefNode when getting value
	GetValueByColName(core.ColumnName) (core.Value, error)
	GetValues() core.Values
	GetColNames() core.ColumnNames
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
		return fmt.Errorf(`ERROR:  relation %v already exist`, tableName)
	}

	colNames := make(core.ColumnNames, 0, len(cols))
	for _, col := range cols {
		colNames = append(colNames, col.ColName)
	}

	db.Tables[tableName] = &DBTable{
		Name:     tableName,
		ColNames: colNames,
		Cols:     cols,
		Rows:     make(DBRows, 0),
	}
	return nil
}

// GetTable gets table from DB
func (db *Database) GetTable(tableName string) (Table, error) {
	if _, ok := db.Tables[tableName]; !ok {
		return nil, fmt.Errorf(`ERROR:  relation "%v" does not exist`, tableName)
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
	return fmt.Errorf(`ERROR: relation "%v" does not exist`, tableName)
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
func (r *DBRow) GetValueByColName(name core.ColumnName) (core.Value, error) {
	for k, v := range r.ColNames {
		if v == name {
			return r.Values[k], nil
		}
	}
	return nil, fmt.Errorf(`ERROR:  column "%v" does not exist`, name.String())
}

// GetValues gets values from DBRow
func (r *DBRow) GetValues() core.Values {
	return r.Values
}

// GetColNames gets column names from DBRow
func (r *DBRow) GetColNames() core.ColumnNames {
	return r.ColNames
}

// UpdateValue updates value by specifing column name
func (r *DBRow) UpdateValue(name core.ColumnName, val core.Value) {
	for k, colName := range r.ColNames {
		if colName.Name == name.Name {
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
	Name     string
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

// GetName return table name
func (t *DBTable) GetName() string {
	return t.Name
}

// GetColNames return column names of table
func (t *DBTable) GetColNames() core.ColumnNames {
	return t.ColNames
}

// GetCols return column names of table
func (t *DBTable) GetCols() core.Cols {
	return t.Cols
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

// RenameTableName updates table name
func (t *DBTable) RenameTableName(name string) {
	t.Name = name

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
func (t *DBTable) Project(TargetColNames core.ColumnNames, resFuncs []func(Row) (core.Value, error)) (Table, error) {
	rows := t.GetRows()
	if len(rows) == 0 {
		return t, nil
	}
	newRows := make(DBRows, 0, len(rows))
	for _, row := range t.Rows {
		colNames := make(core.ColumnNames, 0)
		vals := make(core.Values, 0)
		for k, fn := range resFuncs {
			v, err := fn(row)
			if err != nil {
				return nil, err
			}
			if v != core.Wildcard {
				if v == ColumnNotFound {
					return nil, fmt.Errorf(`ERROR:  column "%v" does not exist`, TargetColNames[k])
				}
				vals = append(vals, v)
				colNames = append(colNames, TargetColNames[k])
			} else {
				// column wildcard
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
		row.Values = vals
		row.ColNames = colNames
		newRows = append(newRows, row)
	}

	t.Rows = newRows

	tbColNames := make(core.ColumnNames, 0)
	for _, name := range newRows[0].ColNames {
		tbColNames = append(tbColNames, name)
	}
	t.ColNames = tbColNames
	// TODO: implement SetCols if type validation is implemented
	// newTable.SetCols(cols)

	return t, nil
}

// Where filters rows by given where conditions
func (t *DBTable) Where(condFn func(Row) (core.Value, error)) (Table, error) {
	srcRows := t.Rows
	rows := make([]*DBRow, 0)
	for _, row := range srcRows {
		v, err := condFn(row)
		if err != nil {
			return nil, err
		}
		if v == core.True {
			rows = append(rows, row)
		}
	}
	t.Rows = rows

	return t, nil
}

// CrossJoin took cross join given tables
func (t *DBTable) CrossJoin(rtb Table) (Table, error) {
	ns := uniteColNames(t.GetColNames(), rtb.GetColNames())
	cols := uniteCols(t.GetCols(), rtb.GetCols())

	rows := make([]*DBRow, 0)
	rs1 := t.GetRows()
	rs2 := rtb.GetRows()
	for _, r1 := range rs1 {
		for _, r2 := range rs2 {
			rows = append(rows, uniteRow(r1, r2).(*DBRow))
		}
	}

	return &DBTable{
		ColNames: ns,
		Cols:     cols,
		Rows:     rows,
	}, nil
}

func uniteRow(r1, r2 Row) Row {
	vals := make(core.Values, 0)
	for _, v := range r1.GetValues() {
		vals = append(vals, v)
	}
	for _, v := range r2.GetValues() {
		vals = append(vals, v)
	}

	cols := make(core.ColumnNames, 0)
	for _, c := range r1.GetColNames() {
		cols = append(cols, c)
	}
	for _, c := range r2.GetColNames() {
		cols = append(cols, c)
	}

	return &DBRow{
		ColNames: cols,
		Values:   vals,
	}
}

func uniteColNames(lcs, rcs core.ColumnNames) core.ColumnNames {
	ns := make(core.ColumnNames, 0)
	for _, c := range lcs {
		ns = append(ns, c)
	}
	for _, c := range rcs {
		ns = append(ns, c)
	}

	return ns
}

func uniteCols(l, r core.Cols) core.Cols {
	cs := make(core.Cols, 0)
	for _, c := range l {
		cs = append(cs, c)
	}
	for _, c := range l {
		cs = append(cs, c)
	}

	return cs
}

// OrderBy sorts rows by given column names
func (t *DBTable) OrderBy(cols core.ColumnNames, sortDirs []int) (Table, error) {
	if err := validateOrderByColumn(t.ColNames, cols); err != nil {
		return nil, err
	}

	rows := t.Rows
	name := cols[0]
	sortDir := sortDirs[0]
	sort.Slice(rows, func(i, j int) bool {
		l, _ := rows[i].GetValueByColName(name)
		r, _ := rows[j].GetValueByColName(name)

		return core.LessForSort(l, r, sortDir)
	})

	t.Rows = rows

	return t, nil
}

func validateOrderByColumn(tbCols, targets core.ColumnNames) error {
	for _, tc := range targets {
		if (tc == core.ColumnName{}) {
			// expresison
			continue
		}
		if !haveColumn(tc, tbCols) {
			return fmt.Errorf(`column "%v" does not exist`, tc.String())
		}
	}

	return nil
}

func haveColumn(c core.ColumnName, cs core.ColumnNames) bool {
	for _, col := range cs {
		if c == col {
			return true
		}
	}

	return false
}

// Limit selects limited number of record
func (t *DBTable) Limit(N int) (Table, error) {
	if len(t.Rows) <= N {
		return t, nil
	}
	oldRows := t.GetRows()
	newRows := make([]*DBRow, 0)
	for i := 0; i < N; i++ {
		row := oldRows[i]
		newRows = append(newRows,
			&DBRow{
				ColNames: row.GetColNames(),
				Values:   row.GetValues(),
			})
	}

	return &DBTable{
		ColNames: t.GetColNames(),
		Cols:     t.GetCols(),
		Rows:     newRows,
	}, nil
}

// Update updates records
func (t *DBTable) Update(colNames core.ColumnNames, condFn func(Row) (core.Value, error), assignValFns []func(Row) (core.Value, error)) (Table, error) {
	rows := t.Rows
	for _, row := range rows {
		a, err := condFn(row)
		if err != nil {
			return nil, err
		}
		if a == core.True {
			for k, name := range colNames {
				v, err := assignValFns[k](row)
				if err != nil {
					return nil, err
				}
				row.UpdateValue(name, v)
			}
		}
	}

	return nil, nil
}

func (t *DBTable) Delete(condFn func(Row) (core.Value, error)) (Table, error) {
	updatedRows := make([]*DBRow, 0)
	for _, row := range t.Rows {
		v, err := condFn(row)
		if err != nil {
			return nil, err
		}
		if v == core.True {
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
