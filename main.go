//+build

package main

import (
	"fmt"

	"github.com/goropikari/mysqlite2/backend"
	"github.com/goropikari/mysqlite2/core"
	trans "github.com/goropikari/mysqlite2/translator"
)

func main() {
	db := prepareDB()

	// evaluate query
	query := "select hoge.name, hoge.id, *, 1000, 1.5, 'taro' from hoge"
	// query := "select true=true, 1, 1000"
	raNode, _ := trans.NewPGTranslator(query).Translate()
	fmt.Println("raNode: ", raNode)
	tb, _ := raNode.Eval(db)
	rows := tb.GetRows()

	for k, row := range rows {
		fmt.Printf("row %v: %v\n", k, row)
	}
}

func prepareDB() backend.DB {
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

	return db
}
