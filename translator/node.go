package translator

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import (
	"github.com/goropikari/mysqlite2/backend"
	"github.com/goropikari/mysqlite2/core"
)

// RelationalAlgebraNode is interface of RelationalAlgebraNode
type RelationalAlgebraNode interface {
	Eval(backend.DB) (backend.Table, error)
}

// TableNode is Node of table
type TableNode struct {
	TableName string
}

// Eval evaluates TableNode
func (t *TableNode) Eval(db backend.DB) (backend.Table, error) {
	tb, err := db.GetTable(t.TableName)
	if err != nil {
		return nil, err
	}

	return tb, err
}

// ProjectionNode is Node of projection operation
type ProjectionNode struct {
	ResTargets     []ExpressionNode
	TargetColNames core.ColumnNames
	Table          RelationalAlgebraNode
}

// Eval evaluates ProjectionNode
func (p *ProjectionNode) Eval(db backend.DB) (backend.Table, error) {
	if p.Table == nil {
		return nil, nil
	}

	tb, err := p.Table.Eval(db)
	if err != nil {
		return nil, err
	}
	if tb == nil {
		return p.evalEmptyTable()
	}
	newTable := tb.Copy()

	resFuncs := p.constructResFunc()

	rows := newTable.GetRows()
	newRows := make([]backend.Row, 0, len(rows))
	for _, row := range rows {
		colNames := make(core.ColumnNames, 0)
		vals := make(core.Values, 0)
		for k, fn := range resFuncs {
			if v := fn(row); v != Wildcard {
				vals = append(vals, v)
				colNames = append(colNames, p.TargetColNames[k])
			} else { // column wildcard
				// Add values
				for _, val := range row.GetValues() {
					if val == nil {
						// Fix me: nil should be converted
						// when the value is inserted.
						vals = append(vals, Null)
					} else {
						vals = append(vals, val)
					}
				}

				// Add columns
				for _, name := range newTable.GetColNames() {
					colNames = append(colNames, name)
				}
			}
		}
		row.SetValues(vals)
		row.SetColNames(colNames)
		newRows = append(newRows, row)
	}

	newTable.SetRows(newRows)
	newTable.SetColNames(p.TargetColNames)
	// TODO: implement SetCols if type validation is implemented
	// newTable.SetCols(cols)

	return newTable, nil
}

func (p *ProjectionNode) constructResFunc() []func(row backend.Row) core.Value {
	resFuncs := make([]func(backend.Row) core.Value, 0, len(p.ResTargets))
	for _, target := range p.ResTargets {
		resFuncs = append(resFuncs, target.Eval())
	}

	return resFuncs
}

func (p *ProjectionNode) evalEmptyTable() (backend.Table, error) {
	resFuncs := p.constructResFunc()
	row := &EmptyTableRow{
		ColNames: p.TargetColNames,
		Values:   make(core.Values, 0),
	}

	for _, fn := range resFuncs {
		row.Values = append(row.Values, fn(row))
	}

	return &EmptyTable{
		ColNames: p.TargetColNames,
		Rows:     []*EmptyTableRow{row},
	}, nil
}

type EmptyTable struct {
	ColNames core.ColumnNames
	Rows     []*EmptyTableRow
}

func (t *EmptyTable) Copy() backend.Table {
	return t
}

func (t *EmptyTable) GetColNames() core.ColumnNames {
	return t.ColNames
}

func (t *EmptyTable) SetColNames(names core.ColumnNames) {}

func (t *EmptyTable) GetRows() []backend.Row {
	rows := make([]backend.Row, 0, len(t.Rows))
	for _, row := range t.Rows {
		rows = append(rows, row)
	}
	return rows
}

func (t *EmptyTable) SetRows(rows []backend.Row) {}

type EmptyTableRow struct {
	ColNames core.ColumnNames
	Values   core.Values
}

func (t *EmptyTable) InsertValues(cs core.ColumnNames, vs core.ValuesList) error { return nil }

func (r *EmptyTableRow) GetValueByColName(core.ColumnName) core.Value {
	return nil
}

func (r *EmptyTableRow) GetValues() core.Values {
	return r.Values
}

func (r *EmptyTableRow) SetValues(vals core.Values)         {}
func (r *EmptyTableRow) SetColNames(names core.ColumnNames) {}

// WhereNode is Node of where clause
type WhereNode struct {
	Condition ExpressionNode
	Table     RelationalAlgebraNode
}

// Eval evaluate WhereNode
func (wn *WhereNode) Eval(db backend.DB) (backend.Table, error) {
	if wn.Table == nil {
		return backend.Table(nil), nil
	}

	tb, err := wn.Table.Eval(db)
	if err != nil {
		return nil, err
	}

	if wn.Condition == nil {
		return tb, nil
	}

	newTable := tb.Copy()
	srcRows := newTable.GetRows()
	condFunc := wn.Condition.Eval()

	rows := make([]backend.Row, 0, len(srcRows))
	for _, row := range srcRows {
		if condFunc(row) == True {
			rows = append(rows, row)
		}
	}

	newTable.SetRows(rows)

	return newTable, nil
}

// CreateTableNode is a node of create statements
type CreateTableNode struct {
	TableName  string
	ColumnDefs core.Cols
}

// Eval evaluates CreateTableNode
func (c *CreateTableNode) Eval(db backend.DB) (backend.Table, error) {
	if err := db.CreateTable(c.TableName, c.ColumnDefs); err != nil {
		return nil, err
	}
	return nil, nil
}

// InsertNode is a node of create statements
type InsertNode struct {
	TableName   string
	ColumnNames core.ColumnNames
	ValuesList  core.ValuesList
}

// Eval evaluates CreateTableNode
func (c *InsertNode) Eval(db backend.DB) (backend.Table, error) {
	tb, err := db.GetTable(c.TableName)
	if err != nil {
		return nil, err
	}
	tb.InsertValues(c.ColumnNames, c.ValuesList)

	return nil, nil
}
