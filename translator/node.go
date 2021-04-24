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
	ResTargets     []Expression
	TargetColNames core.ColumnNames
	Table          RelationalAlgebraNode
}

// Eval evaluates ProjectionNode
func (p *ProjectionNode) Eval(db backend.DB) (backend.Table, error) {
	tb, err := p.Table.Eval(db)
	if err != nil {
		return nil, err
	}
	newTable := tb.Copy()

	resFuncs := make([]func(backend.Row) core.Value, 0, len(p.ResTargets))
	for _, target := range p.ResTargets {
		resFuncs = append(resFuncs, target.Eval())
	}

	rows := newTable.GetRows()
	newRows := make([]backend.Row, 0, len(rows))
	for _, row := range rows {
		vals := make(core.Values, 0, len(p.ResTargets))
		for _, fn := range resFuncs {
			vals = append(vals, fn(row))
		}
		row.SetValues(vals)
		row.SetColNames(p.TargetColNames)
		newRows = append(newRows, row)
	}

	newTable.SetRows(newRows)
	newTable.SetColNames(p.TargetColNames)
	// TODO: implement SetCols if type validation is implemented
	// newTable.SetCols(cols)

	return newTable, nil
}

// WhereNode is Node of where clause
type WhereNode struct {
	Condition Expression
	Table     RelationalAlgebraNode
}

// Eval evaluate WhereNode
func (wn *WhereNode) Eval(db backend.DB) (backend.Table, error) {
	tb, err := wn.Table.Eval(db)
	if err != nil {
		return nil, err
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
