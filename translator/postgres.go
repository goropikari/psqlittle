package translator

import (
	"github.com/goropikari/mysqlite2/core"
	pg_query "github.com/pganalyze/pg_query_go/v2"
)

// Translator is an interface for translator of SQL parse
type Translator interface {
	Translate() RelationalAlgebraNode
}

// PGTranlator is translator for PostgreSQL syntax
type PGTranlator struct {
	Query string
}

// NewPGTranslator is a constructor of PGTranlator
func NewPGTranslator(query string) *PGTranlator {
	return &PGTranlator{
		Query: query,
	}
}

// Translate translates a postgres parse tree into RelationalAlgebraNode
func (pg *PGTranlator) Translate() (RelationalAlgebraNode, error) {
	result, err := pg_query.Parse(pg.Query)
	if err != nil {
		return nil, err
	}

	stmt := result.Stmts[0].Stmt
	if node := stmt.GetSelectStmt(); node != nil {
		return pg.TranslateSelect(node)
	}
	return nil, nil
}

// TranslateSelect translates postgres a select statement into ProjectionNode
func (pg *PGTranlator) TranslateSelect(pgtree *pg_query.SelectStmt) (RelationalAlgebraNode, error) {
	return nil, nil
	// targetColNames := retrieveColNameTargetList(pgtree.GetTargetList())
	// return &ProjectionNode{
	// 	TargetCols: targetColNames,
	// 	Table:      nil,
	// }, nil
}

func retrieveColNameTargetList(targetList []*pg_query.Node) core.ColNames {
	if targetList == nil {
		return nil
	}

	names := make(core.ColNames, 0, len(targetList))
	for _, target := range targetList {
		if target.GetResTarget().GetVal().GetColumnRef() != nil {
			fields := target.GetResTarget().GetVal().GetColumnRef().GetFields()
			if len(fields) == 2 { // column is specified by table name and column name
				tableName := fields[0].GetString_().GetStr()
				colName := fields[1].GetString_().GetStr()
				names = append(names, core.ColName{TableName: tableName, Name: colName})
			} else {
				// Not Implemented
			}
		}
		if target.GetResTarget().GetVal().GetAConst() != nil {
			// Not Implemented
		}
	}

	return names
}
