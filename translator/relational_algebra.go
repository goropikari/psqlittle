package translator

//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

import (
	"errors"
	"fmt"

	"github.com/goropikari/psqlittle/backend"
	"github.com/goropikari/psqlittle/core"
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

// RenameTableNode is Node for renaming tabel
type RenameTableNode struct {
	Alias string
	Table RelationalAlgebraNode
}

// Eval evaluates RenameTableNode
func (rt *RenameTableNode) Eval(db backend.DB) (backend.Table, error) {
	if rt.Table == nil {
		return nil, errors.New("have to include table")
	}

	tb, err := rt.Table.Eval(db)
	if err != nil {
		return nil, err
	}

	newTable := tb.Copy()
	newTable.RenameTableName(rt.Alias)

	return newTable, nil
}

// ProjectionNode is Node of projection operation
type ProjectionNode struct {
	ResTargets     []ExpressionNode
	TargetColNames core.ColumnNames
	RANode         RelationalAlgebraNode
}

// Eval evaluates ProjectionNode
func (p *ProjectionNode) Eval(db backend.DB) (backend.Table, error) {
	if p.RANode == nil {
		return nil, nil
	}

	tb, err := p.RANode.Eval(db)
	if err != nil {
		return nil, err
	}
	if tb == nil {
		return p.makeEmptyTable()
	}
	newTable := tb.Copy()

	resFuncs := p.constructResFunc()

	if err := validateTargetColumn(newTable.GetColNames(), p.TargetColNames); err != nil {
		return nil, err
	}

	return newTable.Project(p.TargetColNames, resFuncs)
}

func validateTargetColumn(tbCols core.ColumnNames, targets core.ColumnNames) error {
	for _, tc := range targets {
		if (tc == core.ColumnName{Name: "*"}) {
			continue
		}
		if (tc == core.ColumnName{}) {
			// expresison
			continue
		}
		if !haveColumn(tc, tbCols) {
			return fmt.Errorf(`column "%v" does not exist`, makeColName(tc))
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

func makeColName(c core.ColumnName) string {
	if c.TableName == "" {
		return c.Name
	}
	return c.TableName + "." + c.Name
}

func (p *ProjectionNode) constructResFunc() []func(row backend.Row) (core.Value, error) {
	resFuncs := make([]func(backend.Row) (core.Value, error), 0, len(p.ResTargets))
	for _, target := range p.ResTargets {
		resFuncs = append(resFuncs, target.Eval())
	}

	return resFuncs
}

func (p *ProjectionNode) makeEmptyTable() (backend.Table, error) {
	resFuncs := p.constructResFunc()
	row := &EmptyTableRow{
		ColNames: p.TargetColNames,
		Values:   make(core.Values, 0),
	}

	for _, fn := range resFuncs {
		v, err := fn(row)
		if err != nil {
			return nil, err
		}
		row.Values = append(row.Values, v)
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

func (t *EmptyTable) GetName() string {
	return ""
}

func (t *EmptyTable) GetColNames() core.ColumnNames {
	return t.ColNames
}

func (t *EmptyTable) GetRows() []backend.Row {
	rows := make([]backend.Row, 0, len(t.Rows))
	for _, row := range t.Rows {
		rows = append(rows, row)
	}
	return rows
}

func (t *EmptyTable) GetCols() core.Cols {
	return nil
}

func (t *EmptyTable) RenameTableName(name string) {}

func (t *EmptyTable) InsertValues(cs core.ColumnNames, vs core.ValuesList) error { return nil }

func (t *EmptyTable) Project(cs core.ColumnNames, fns []func(backend.Row) (core.Value, error)) (backend.Table, error) {
	return nil, nil
}

func (t *EmptyTable) Where(fn func(backend.Row) (core.Value, error)) (backend.Table, error) {
	return nil, nil
}

func (t *EmptyTable) CrossJoin(backend.Table) (backend.Table, error) {
	return nil, nil
}

func (t *EmptyTable) OrderBy(ns core.ColumnNames, dirs []int) (backend.Table, error) {
	return nil, nil
}

func (t *EmptyTable) Limit(n int) (backend.Table, error) {
	return nil, nil
}

func (t *EmptyTable) Update(colNames core.ColumnNames, condFn func(backend.Row) (core.Value, error), assignValFns []func(backend.Row) (core.Value, error)) (backend.Table, error) {
	return nil, nil
}

func (t *EmptyTable) Delete(func(backend.Row) (core.Value, error)) (backend.Table, error) {
	return nil, nil
}

type EmptyTableRow struct {
	ColNames core.ColumnNames
	Values   core.Values
}

func (r *EmptyTableRow) GetValueByColName(core.ColumnName) (core.Value, error) {
	return nil, nil
}

func (r *EmptyTableRow) GetValues() core.Values {
	return r.Values
}

func (r *EmptyTableRow) GetColNames() core.ColumnNames {
	return nil
}

func (r *EmptyTableRow) UpdateValue(name core.ColumnName, val core.Value) {}

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
	condFunc := wn.Condition.Eval()

	return newTable.Where(condFunc)
}

// CrossJoinNode is a node of cross join.
type CrossJoinNode struct {
	RANodes []RelationalAlgebraNode
}

// Eval evaluates CrossJoinNode
func (c *CrossJoinNode) Eval(db backend.DB) (backend.Table, error) {
	tbs := make([]backend.Table, 0, len(c.RANodes))

	for _, ra := range c.RANodes {
		tb, err := ra.Eval(db)
		if err != nil {
			return nil, err
		}
		tbs = append(tbs, tb)
	}

	if err := validateTableName(tbs); err != nil {
		return nil, err
	}

	switch len(tbs) {
	case 0: // when table not specified. Only select expression case.
		return nil, nil
	case 1:
		return tbs[0], nil
	}

	var tb backend.Table
	var err error
	tb = tbs[0]
	for k, t := range tbs {
		if k == 0 {
			continue
		}
		tb, err = crossJoinTable(tb, t)
		if err != nil {
			return nil, err
		}
	}

	return tb, nil
}

func validateTableName(tbs []backend.Table) error {

	nm := make(map[string]int)
	for _, tb := range tbs {
		if name := tb.GetName(); name != "" {
			if _, ok := nm[name]; ok {
				return fmt.Errorf(`ERROR:  table name "%v" specified more than once`, name)
			}
			nm[name]++
		}
	}

	return nil
}

func crossJoinTable(tb1, tb2 backend.Table) (backend.Table, error) {
	if tb1 == nil || tb2 == nil {
		return nil, nil
	}

	tb, err := tb1.CrossJoin(tb2)
	if err != nil {
		return nil, err
	}

	return tb, nil
}

// OrderByNode is a Node for order by clause
type OrderByNode struct {
	SortKeys core.ColumnNames
	SortDirs []int
	RANode   RelationalAlgebraNode
}

// Eval evaluates OrderByNode
func (o *OrderByNode) Eval(db backend.DB) (backend.Table, error) {
	tb, err := o.RANode.Eval(db)
	if err != nil {
		return nil, err
	}

	tb.OrderBy(o.SortKeys, o.SortDirs)

	return tb, nil
}

// LimitNode is a Node for limit clause
type LimitNode struct {
	Count  int
	RANode RelationalAlgebraNode
}

// Eval evaluates LimitNode
func (l *LimitNode) Eval(db backend.DB) (backend.Table, error) {
	tb, err := l.RANode.Eval(db)
	if err != nil {
		return nil, err
	}

	tb, err = tb.Limit(l.Count)
	if err != nil {
		return nil, err
	}
	return tb, nil
}

// DropTableNode is a node of drop statement
type DropTableNode struct {
	TableNames []string
}

// Eval evaluates DropTableNode
func (d *DropTableNode) Eval(db backend.DB) (backend.Table, error) {
	for _, name := range d.TableNames {
		if err := db.DropTable(name); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// CreateTableNode is a node of create statement
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

// InsertNode is a node of create statement
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

// UpdateNode is a node of update statement
type UpdateNode struct {
	Condition  ExpressionNode
	ColNames   core.ColumnNames
	AssignExpr []ExpressionNode
	TableName  string
}

// Eval evaluates UpdateNode
func (u *UpdateNode) Eval(db backend.DB) (backend.Table, error) {
	var condFunc func(backend.Row) (core.Value, error)
	if u.Condition == nil {
		condFunc = func(row backend.Row) (core.Value, error) {
			return core.True, nil
		}
	} else {
		condFunc = u.Condition.Eval()
	}

	tb, err := db.GetTable(u.TableName)
	if err != nil {
		return nil, err
	}

	assignValFns := make([]func(backend.Row) (core.Value, error), 0)
	for _, expr := range u.AssignExpr {
		assignValFns = append(assignValFns, expr.Eval())
	}

	tb.Update(u.ColNames, condFunc, assignValFns)

	return nil, nil
}

// DeleteNode is a node of update statement
type DeleteNode struct {
	Condition ExpressionNode
	TableName string
}

// Eval evaluates DeleteNode
func (d *DeleteNode) Eval(db backend.DB) (backend.Table, error) {
	var condFunc func(backend.Row) (core.Value, error)
	if d.Condition == nil {
		condFunc = func(row backend.Row) (core.Value, error) {
			return core.True, nil
		}
	} else {
		condFunc = d.Condition.Eval()
	}

	tb, err := db.GetTable(d.TableName)
	if err != nil {
		return nil, err
	}

	return tb.Delete(condFunc)
}
