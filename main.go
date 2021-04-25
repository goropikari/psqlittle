//+build

package main

import (
	"fmt"

	"github.com/goropikari/mysqlite2/backend"
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
	db := backend.NewDatabase()

	query := "create table hoge (id int, name varchar(255))"
	raNode, _ := trans.NewPGTranslator(query).Translate()
	raNode.Eval(db)
	if _, err := raNode.Eval(db); err != nil {
		fmt.Println("error:", err)
	}

	query = "insert into hoge (name, id) values ('taro', 9876)"
	raNode, _ = trans.NewPGTranslator(query).Translate()
	raNode.Eval(db)
	query = "insert into hoge (name, id) values ('hanako', 12343)"
	raNode, _ = trans.NewPGTranslator(query).Translate()
	raNode.Eval(db)
	query = "insert into hoge (name, id) values ('mike', 7893)"
	raNode, _ = trans.NewPGTranslator(query).Translate()
	raNode.Eval(db)

	return db
}
