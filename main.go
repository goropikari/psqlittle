//+build

package main

import (
	"fmt"

	"github.com/goropikari/mysqlite2/backend"
	"github.com/goropikari/mysqlite2/core"
	trans "github.com/goropikari/mysqlite2/translator"
)

func main() {

	tableName := "hoge"
	db := backend.NewDatabase()
	cols := backend.Cols{
		{
			ColName: core.ColumnName{
				TableName: tableName,
				Name:      "id",
			},
			ColType: core.Integer,
		},
		{
			ColName: core.ColumnName{
				TableName: tableName,
				Name:      "name",
			},
			ColType: core.VarChar,
		},
	}

	colNames := make(core.ColumnNames, 0, len(cols))
	for _, col := range cols {
		colNames = append(colNames, col.ColName)
	}

	vals := core.ValuesList{
		core.Values{1, "Hello"},
		core.Values{1, "Hello"},
		core.Values{1, "Hello"},
		core.Values{nil, "World"},
	}

	db.CreateTable("hoge", cols)
	table, _ := db.GetTable("hoge")
	err := table.(*backend.DBTable).Insert(colNames, vals)
	if err != nil {
		fmt.Println(err)
	}
	table, _ = db.GetTable("hoge")
	fmt.Println("after insert:", table.(*backend.DBTable).Rows[0])

	whereCond := trans.BinOpNode{
		Op: trans.EqualOp,
		Lexpr: trans.ColRefNode{
			ColName: core.ColumnName{
				TableName: tableName,
				Name:      "id",
			},
		},
		Rexpr: trans.IntegerNode{
			Val: 1,
		},
		// Rexpr: trans.BoolConstNode{
		// 	Bool: trans.Null,
		// },
	}

	// whereCond := trans.NullTestNode{
	// 	TestType: trans.EqualNull,
	// 	Expr: trans.ColRefNode{
	// 		core.ColumnName{
	// 			TableName: "hoge",
	// 			Name:      "id",
	// 		},
	// 	},
	// }

	whereNode := trans.WhereNode{
		Condition: whereCond,
		Table: &trans.TableNode{
			TableName: tableName,
		},
	}

	projectNode := &trans.ProjectionNode{
		TargetCols: core.ColumnNames{
			{TableName: "hoge", Name: "name"},
			{TableName: "hoge", Name: "id"},
		},
		Table: &whereNode,
	}

	// tb, _ := whereNode.Eval(db)
	tb, _ := projectNode.Eval(db)
	fmt.Println(tb)
	rows := tb.GetRows()
	for k, row := range rows {
		fmt.Printf("row %v: %v\n", k, row)
	}
}
