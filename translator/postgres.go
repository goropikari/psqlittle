package translator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/goropikari/mysqlite2/backend"
	"github.com/goropikari/mysqlite2/core"
	pg_query "github.com/pganalyze/pg_query_go/v2"
)

// Result is interface of query result
type Result interface {
	GetRecords() core.ValuesList
	GetColumns() []string
}

// QueryResult is result of query
type QueryResult struct {
	Columns []string
	Records core.ValuesList
}

// GetRecords gets records from query result
func (qr *QueryResult) GetRecords() core.ValuesList {
	return qr.Records
}

// GetColumns gets column name of records
func (qr *QueryResult) GetColumns() []string {
	return qr.Columns
}

// Statement is interface of query statement
type Statement interface {
	Eval(backend.DB) (Result, error)
}

// QueryStatement is statement of query
type QueryStatement struct {
	RANode RelationalAlgebraNode
}

// Eval evaluates QueryStatement
func (qs *QueryStatement) Eval(db backend.DB) (Result, error) {
	tb, err := qs.RANode.Eval(db)
	if err != nil {
		return nil, err
	}

	if tb == nil {
		return &QueryResult{}, nil
	}

	rows := tb.GetRows()
	cols := tb.GetColNames()

	recs := make(core.ValuesList, 0)
	names := make([]string, 0)

	for _, row := range rows {
		rec := make(core.Values, 0)
		for _, val := range row.GetValues() {
			switch val {
			case core.True:
				rec = append(rec, true)
			case core.False:
				rec = append(rec, false)
			case core.Null:
				rec = append(rec, nil)
			default:
				rec = append(rec, val)
			}
		}
		recs = append(recs, rec)
	}
	for _, col := range cols {
		names = append(names, col.Name)
	}

	return &QueryResult{
		Columns: names,
		Records: recs,
	}, nil
}

// Translator is an interface for translator of SQL parse
type Translator interface {
	Translate() RelationalAlgebraNode
}

// PGTranlator is translator for PostgreSQL syntax
type PGTranlator struct {
	query string
}

// NewPGTranslator is a constructor of PGTranlator
func NewPGTranslator(query string) *PGTranlator {
	return &PGTranlator{
		query: query,
	}
}

// Translate translates a postgres parse tree into RelationalAlgebraNode
func (pg *PGTranlator) Translate() (Statement, error) {
	result, err := pg_query.Parse(pg.query)
	if err != nil {
		return nil, err
	}

	var ra RelationalAlgebraNode
	stmt := result.Stmts[0].Stmt
	if node := stmt.GetSelectStmt(); node != nil {
		ra, err = pg.TranslateSelect(node)
	}
	if node := stmt.GetCreateStmt(); node != nil {
		ra, err = pg.TranslateCreateTable(node)
	}
	if node := stmt.GetDropStmt(); node != nil {
		ra, err = pg.TranslateDropTable(node)
	}
	if node := stmt.GetInsertStmt(); node != nil {
		ra, err = pg.TranslateInsert(node)
	}
	if node := stmt.GetUpdateStmt(); node != nil {
		ra, err = pg.TranslateUpdate(node)
	}
	if node := stmt.GetDeleteStmt(); node != nil {
		ra, err = pg.TranslateDelete(node)
	}

	if ra != nil {
		return &QueryStatement{
			RANode: ra,
		}, nil
	}

	return nil, errors.New("Don't support such query")
}

// TranslateDropTable translates sql parse tree into DropTableNode
func (pg *PGTranlator) TranslateDropTable(node *pg_query.DropStmt) (RelationalAlgebraNode, error) {
	tableList := node.GetObjects()
	tableNames := make([]string, 0)
	for _, tb := range tableList {
		tableNames = append(tableNames, tb.GetList().GetItems()[0].GetString_().GetStr())
	}

	return &DropTableNode{
		TableNames: tableNames,
	}, nil
}

// TranslateDelete translates sql parse tree into DeleteNode
func (pg *PGTranlator) TranslateDelete(node *pg_query.DeleteStmt) (RelationalAlgebraNode, error) {
	cond := constructExprNode(node.GetWhereClause())
	tableName := node.GetRelation().Relname

	return &DeleteNode{
		Condition: cond,
		TableName: tableName,
	}, nil
}

// TranslateUpdate translates sql parse tree into UpdateNode
func (pg *PGTranlator) TranslateUpdate(node *pg_query.UpdateStmt) (RelationalAlgebraNode, error) {
	cond := constructExprNode(node.GetWhereClause())
	tableName := strings.ToLower(node.GetRelation().Relname)
	targetColNames, resTargetNodes := interpreteUpdateTargetList(node.GetTargetList())

	return &UpdateNode{
		Condition:  cond,
		ColNames:   targetColNames,
		AssignExpr: resTargetNodes,
		TableName:  tableName,
	}, nil
}

// TranslateSelect translates postgres a select statement into ProjectionNode
func (pg *PGTranlator) TranslateSelect(pgtree *pg_query.SelectStmt) (RelationalAlgebraNode, error) {
	targetList := pgtree.GetTargetList()
	targetColNames, resTargetNodes := interpreteTargetList(targetList)

	table, err := pg.interpretFromClause(pgtree.GetFromClause())
	if err != nil {
		return nil, err
	}
	whereNode := constructWhereNode(pgtree.GetWhereClause(), table)
	orderByNode := constructOrderByNode(pgtree.GetSortClause(), whereNode)
	limitNode, err := constructLimitNode(pgtree.GetLimitCount(), orderByNode)
	if err != nil {
		return nil, err
	}

	return &ProjectionNode{
		TargetColNames: targetColNames,
		ResTargets:     resTargetNodes,
		RANode:         limitNode,
	}, nil
}

func constructLimitNode(limitNode *pg_query.Node, orderByNode RelationalAlgebraNode) (RelationalAlgebraNode, error) {
	if limitNode == nil {
		return orderByNode, nil
	}

	expr := constructExprNode(limitNode)
	var r backend.Row
	val, err := expr.Eval()(r)
	if err != nil {
		return nil, err
	}

	if reflect.ValueOf(val).Kind() != reflect.Int {
		return nil, errors.New("ERROR:  argument of LIMIT must not contain variables")
	}

	return &LimitNode{
		Count:  val.(int),
		RANode: orderByNode,
	}, nil
}

func constructOrderByNode(sortTree []*pg_query.Node, whereNode RelationalAlgebraNode) RelationalAlgebraNode {
	if len(sortTree) == 0 {
		return whereNode
	}
	sortKeys := make(core.ColumnNames, 0)
	sortDirs := make([]int, 0)
	sortKey := getColName(sortTree[0].GetSortBy().GetNode().GetColumnRef())
	sortDir := int(sortTree[0].GetSortBy().GetSortbyDir())
	// TODO: get all specified columns
	sortKeys = append(sortKeys, sortKey)
	sortDirs = append(sortDirs, sortDir)

	return &OrderByNode{
		SortKeys: sortKeys,
		SortDirs: sortDirs,
		RANode:   whereNode,
	}
}

func (pg *PGTranlator) interpretFromClause(fromTree []*pg_query.Node) (RelationalAlgebraNode, error) {
	tables := make([]RelationalAlgebraNode, 0, len(fromTree))

	for _, relation := range fromTree {
		if relation.GetRangeVar() != nil {
			tableName := relation.GetRangeVar().Relname
			alias := relation.GetRangeVar().Alias.GetAliasname()
			if alias == "" {
				tables = append(tables, &TableNode{TableName: tableName})
			} else {
				tables = append(tables, &RenameTableNode{
					Alias: alias,
					Table: &TableNode{
						TableName: tableName,
					},
				})
			}
		}
		if relation.GetRangeSubselect() != nil {
			subQueryTree := relation.GetRangeSubselect().GetSubquery().GetSelectStmt()
			alias := relation.GetRangeSubselect().Alias.GetAliasname()
			ra, err := pg.TranslateSelect(subQueryTree)
			if err != nil {
				return nil, err
			}
			tables = append(tables, &RenameTableNode{
				Alias: alias,
				Table: ra,
			})
		}
		if relation.GetJoinExpr() != nil {
			// Not Implemented
		}
	}

	table := crossJoinRA(tables)

	return table, nil
}

func crossJoinRA(ras []RelationalAlgebraNode) RelationalAlgebraNode {

	return &CrossJoinNode{
		RANodes: ras,
	}
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
			if colRef.GetFields()[0].GetAStar() != nil {
				resExprs = append(resExprs, ColWildcardNode{})
				names = append(names, core.ColumnName{Name: "*"})
			} else {
				colName := getColName(colRef)
				names = append(names, colName)
				resExprs = append(resExprs, ColRefNode{ColName: colName})
			}
		} else {
			// The column is an expression.
			// This column is not included in given table.
			names = append(names, core.ColumnName{})
			resExprs = append(resExprs, constructExprNode(val))
		}
	}

	return names, resExprs
}

func interpreteUpdateTargetList(targetList []*pg_query.Node) (core.ColumnNames, []ExpressionNode) {
	if targetList == nil {
		return nil, nil
	}

	names := make(core.ColumnNames, 0, len(targetList))
	resExprs := make([]ExpressionNode, 0, len(targetList))
	for _, target := range targetList {
		tableName := strings.ToLower(target.GetResTarget().GetName())
		colName := strings.ToLower(target.GetResTarget().GetIndirection()[0].GetString_().GetStr())
		val := constructExprNode(target.GetResTarget().GetVal())

		names = append(names, core.ColumnName{TableName: tableName, Name: colName})
		resExprs = append(resExprs, val)
	}

	return names, resExprs
}

// TranslateCreateTable translates sql parse tree into CreateTableNode
func (pg *PGTranlator) TranslateCreateTable(stmt *pg_query.CreateStmt) (RelationalAlgebraNode, error) {
	tableName := strings.ToLower(stmt.GetRelation().GetRelname())
	colDefs := prepareColDefs(stmt.GetTableElts(), tableName)

	return &CreateTableNode{
		TableName:  tableName,
		ColumnDefs: colDefs,
	}, nil
}

// TranslateInsert translates sql parse tree into InsertNode
func (pg *PGTranlator) TranslateInsert(stmt *pg_query.InsertStmt) (RelationalAlgebraNode, error) {
	tableName := strings.ToLower(stmt.GetRelation().GetRelname())
	rawValsLists := stmt.GetSelectStmt().GetSelectStmt().GetValuesLists()

	valsLists := make(core.ValuesList, 0, len(rawValsLists))
	for _, rawVals := range rawValsLists {
		items := rawVals.GetList().GetItems()
		vals := make(core.Values, 0, len(items))
		for _, item := range items {
			var r backend.Row
			val, err := constructExprNode(item).Eval()(r)
			if err != nil {
				return nil, err
			}
			vals = append(vals, val)
		}
		valsLists = append(valsLists, vals)
	}

	cols := stmt.GetCols()
	colNames := make(core.ColumnNames, 0, len(cols))
	for _, col := range cols {
		colNames = append(colNames, core.ColumnName{
			TableName: tableName,
			Name:      strings.ToLower(col.GetResTarget().GetName()),
		})
	}

	return &InsertNode{
		TableName:   tableName,
		ColumnNames: colNames,
		ValuesList:  valsLists,
	}, nil
}

func prepareColDefs(defNodes []*pg_query.Node, tableName string) core.Cols {
	colTyps := make(core.Cols, 0, len(defNodes))
	for _, defNode := range defNodes {
		def := defNode.GetColumnDef()
		name := def.GetColname()
		typ := mapGoType(def.GetTypeName().GetNames()[1].GetString_().GetStr())
		col := core.Col{
			ColName: core.ColumnName{
				TableName: strings.ToLower(tableName),
				Name:      strings.ToLower(name),
			},
			ColType: typ,
		}
		colTyps = append(colTyps, col)
	}

	return colTyps
}

func mapGoType(typ string) core.ColType {
	switch typ {
	case "int4":
		return core.Integer
	case "varchar":
		return core.VarChar
	}

	return core.Integer
}

func getColName(colRef *pg_query.ColumnRef) core.ColumnName {
	fields := colRef.GetFields()
	if len(fields) == 1 {
		// column is specified by column name
		colName := strings.ToLower(fields[0].GetString_().GetStr())
		return core.ColumnName{Name: colName}
	}
	if len(fields) == 2 {
		// column is specified by table name and column name
		tableName := strings.ToLower(fields[0].GetString_().GetStr())
		colName := strings.ToLower(fields[1].GetString_().GetStr())
		return core.ColumnName{TableName: tableName, Name: colName}
	}

	// Not Implemented
	// This columnRef includes schema name or something.
	fmt.Println("Not Implemented: This columnRef includes schema name or something.")
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
			return BoolConstNode{Bool: core.Null}
		}
	}
	if node.GetTypeCast() != nil {
		return interpretTypeCast(node.GetTypeCast())
	}
	if v := node.GetAExpr(); v != nil {
		return constructGetAExprNode(v)
	}
	if v := node.GetBoolExpr(); v != nil {
		return constructBoolExprNode(v)
	}
	if v := node.GetColumnRef(); v != nil {
		return constructColumnRef(v)
	}
	if v := node.GetNullTest(); v != nil {
		return constructNullTest(v)
	}
	if v := node.GetCaseExpr(); v != nil {
		return constructCaseNode(v)
	}

	// Not Implemented
	fmt.Println("Not Implemented")
	dummy := IntegerNode{Val: -1 << 60}
	return dummy
}

func constructCaseNode(node *pg_query.CaseExpr) ExpressionNode {
	var caseWhenExprs, caseResultExprs []ExpressionNode
	if arg := node.GetArg(); arg != nil {
		caseWhenExprs, caseResultExprs = constructCaseWithArgNode(node)
	} else {
		caseWhenExprs, caseResultExprs = constructCaseWithoutArgNode(node)
	}

	var defRes ExpressionNode
	if v := node.GetDefresult(); v != nil {
		defRes = constructExprNode(v)
	} else {
		defRes = &BoolConstNode{Bool: core.Null}
	}

	return &CaseNode{
		CaseWhenExprs:   caseWhenExprs,
		CaseResultExprs: caseResultExprs,
		DefaultResult:   defRes,
	}
}

func constructCaseWithoutArgNode(node *pg_query.CaseExpr) ([]ExpressionNode, []ExpressionNode) {
	caseWhenExprs := make([]ExpressionNode, 0)
	caseResultExprs := make([]ExpressionNode, 0)
	for _, caseWhen := range node.GetArgs() {
		caseWhenExprs = append(caseWhenExprs, constructExprNode(caseWhen.GetCaseWhen().GetExpr()))
		caseResultExprs = append(caseResultExprs,
			constructExprNode(caseWhen.GetCaseWhen().GetResult()))
	}

	return caseWhenExprs, caseResultExprs
}

func constructCaseWithArgNode(node *pg_query.CaseExpr) ([]ExpressionNode, []ExpressionNode) {
	arg := constructExprNode(node.GetArg())
	caseWhenExprs := make([]ExpressionNode, 0)
	caseResultExprs := make([]ExpressionNode, 0)
	for _, caseWhen := range node.GetArgs() {
		caseWhenExprs = append(caseWhenExprs,
			BinOpNode{
				Op:    EqualOp,
				Lexpr: arg,
				Rexpr: constructExprNode(caseWhen.GetCaseWhen().GetExpr()),
			})
		caseResultExprs = append(caseResultExprs,
			constructExprNode(caseWhen.GetCaseWhen().GetResult()))
	}

	return caseWhenExprs, caseResultExprs

}

func constructNullTest(node *pg_query.NullTest) ExpressionNode {
	expr := constructExprNode(node.GetArg())
	testtyp := node.GetNulltesttype()
	switch testtyp {
	case 1: // is null
		return &NullTestNode{
			TestType: EqualNull,
			Expr:     expr,
		}
	case 2: // is not null
		return &NullTestNode{
			TestType: NotEqualNull,
			Expr:     expr,
		}
	}

	return nil
}

func constructBoolExprNode(node *pg_query.BoolExpr) ExpressionNode {
	opType := node.GetBoolop()
	switch opType {
	case 1: // AND
		return &ANDNode{
			Lexpr: constructExprNode(node.GetArgs()[0]),
			Rexpr: constructExprNode(node.GetArgs()[1]),
		}
	case 2: // OR
		return &ORNode{
			Lexpr: constructExprNode(node.GetArgs()[0]),
			Rexpr: constructExprNode(node.GetArgs()[1]),
		}
	}

	return nil
}

func constructColumnRef(node *pg_query.ColumnRef) ExpressionNode {
	colName := getColName(node)
	return &ColRefNode{
		ColName: colName,
	}
}

func interpretTypeCast(c *pg_query.TypeCast) ExpressionNode {
	// Now, only support bool
	boolStr := c.GetArg().GetAConst().GetVal().GetString_().GetStr()
	if boolStr == "t" {
		return BoolConstNode{Bool: core.True}
	}
	return BoolConstNode{Bool: core.False}
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
	// ref: translator/const.go: MathOp
	// ref: translator/expression.go: func (e BinOpNode) Eval()

	switch op {
	case "=":
		return EqualOp
	case "!=", "<>":
		return NotEqualOp
	case "+":
		return Plus
	case "-":
		return Minus
	case "*":
		return Multiply
	case "/":
		return Divide
	case ">":
		return GT
	case "<":
		return LT
	case ">=":
		return GEQ
	case "<=":
		return LEQ
	case "||":
		return CONCAT
	}

	fmt.Println("Not Implemented math operator")

	return -1
}
