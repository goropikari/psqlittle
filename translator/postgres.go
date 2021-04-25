package translator

import (
	"strconv"

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
	targetList := pgtree.GetTargetList()
	targetColNames, resTargetNodes := interpreteTargetList(targetList)

	table := interpretFromClause(pgtree.GetFromClause())
	whereNode := constructWhereNode(pgtree.GetWhereClause(), table)

	return &ProjectionNode{
		TargetColNames: targetColNames,
		ResTargets:     resTargetNodes,
		Table:          whereNode,
	}, nil

	// return nil, nil
}

func interpretFromClause(fromTree []*pg_query.Node) RelationalAlgebraNode {
	tables := make([]RelationalAlgebraNode, 0, len(fromTree))

	for _, relation := range fromTree {
		if relation.GetRangeVar() != nil {
			tableName := relation.GetRangeVar().Relname
			alias := relation.GetRangeVar().Alias
			if alias == nil {
				tables = append(tables, &TableNode{TableName: tableName})
			} else {
				// Not Implemented
				// Construct Rename Node
			}
		}
		if relation.GetJoinExpr() != nil {
			// Not Implemented
		}
	}

	table := crossJoin(tables)

	return table
}

func crossJoin(tables []RelationalAlgebraNode) RelationalAlgebraNode {
	if len(tables) == 0 {
		return nil
	}
	return tables[0]
}

func constructWhereNode(whereTree *pg_query.Node, table RelationalAlgebraNode) RelationalAlgebraNode {
	cond := constructExprNode(whereTree)
	return &WhereNode{
		Condition: cond,
		Table:     table,
	}
}

func interpreteTargetList(targetList []*pg_query.Node) (core.ColumnNames, []ExpressionNode) {
	if targetList == nil {
		return nil, nil
	}

	names := make(core.ColumnNames, 0, len(targetList))
	resExprs := make([]ExpressionNode, 0, len(targetList))
	for _, target := range targetList {
		val := target.GetResTarget().GetVal()
		if colRef := val.GetColumnRef(); colRef != nil {
			colName := getColName(colRef)
			names = append(names, colName)
			resExprs = append(resExprs, ColRefNode{ColName: colName})
		} else {
			// This column is not included in given table.
			// The column is an expression.
			names = append(names, core.ColumnName{})
			resExprs = append(resExprs, constructExprNode(val))
		}
	}

	return names, resExprs
}

func getColName(colRef *pg_query.ColumnRef) core.ColumnName {
	fields := colRef.GetFields()
	if len(fields) == 1 {
		// column is specified by column name
		colName := fields[0].GetString_().GetStr()
		return core.ColumnName{Name: colName}
	}
	if len(fields) == 2 {
		// column is specified by table name and column name
		tableName := fields[0].GetString_().GetStr()
		colName := fields[1].GetString_().GetStr()
		return core.ColumnName{TableName: tableName, Name: colName}
	}

	// Not Implemented
	// The column is included schema name or something.
	return core.ColumnName{}
}

func constructExprNode(node *pg_query.Node) ExpressionNode {
	if node == nil {
		return nil
	}

	if node.GetAConst() != nil {
		val := node.GetAConst().GetVal()
		if val.GetInteger() != nil {
			return IntegerNode{Val: int(val.GetInteger().GetIval())}
		}
		if val.GetFloat() != nil {
			f, _ := strconv.ParseFloat(val.GetFloat().GetStr(), 64)
			return FloatNode{Val: f}
		}
		if val.GetString_() != nil {
			return StringNode{Val: val.GetString_().GetStr()}
		}
		if val.GetNull() != nil {
			return BoolConstNode{Bool: Null}
		}
	}
	if node.GetTypeCast() != nil {
		return interpretTypeCast(node.GetTypeCast())
	}
	if node.GetAExpr() != nil {
		return constructGetAExprNode(node.GetAExpr())
	}
	if node.GetBoolExpr() != nil {
		return nil
	}

	// Not Implemented
	dummy := IntegerNode{Val: -1 << 60}
	return dummy
}

func interpretTypeCast(c *pg_query.TypeCast) ExpressionNode {
	// Now, only support bool
	boolStr := c.GetArg().GetAConst().GetVal().GetString_().GetStr()
	if boolStr == "t" {
		return BoolConstNode{Bool: True}
	}
	return BoolConstNode{Bool: False}
}

func constructGetAExprNode(aExpr *pg_query.A_Expr) ExpressionNode {
	op := mathOperator(aExpr.GetName()[0].GetString_().GetStr())
	lexpr := constructExprNode(aExpr.GetLexpr())
	rexpr := constructExprNode(aExpr.GetRexpr())

	return &BinOpNode{
		Op:    op,
		Lexpr: lexpr,
		Rexpr: rexpr,
	}
}

func mathOperator(op string) MathOp {
	switch op {
	case "=":
		return EqualOp
	case "!=":
		return NotEqualOp
	}

	return -1
}
