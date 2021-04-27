//+build

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/goropikari/mysqlite2/backend"
	trans "github.com/goropikari/mysqlite2/translator"
)

func main() {
	// db := prepareDB()

	db := backend.NewDatabase()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("sql> ")
		query, err := reader.ReadString(';')
		if err != nil {
			fmt.Println(err)
			continue
		}
		if query == ".exit;" {
			os.Exit(0)
		}
		query = strings.Trim(query, " \n")
		if query == ";" {
			continue
		}

		raNode, err := trans.NewPGTranslator(query).Translate()
		if err != nil {
			fmt.Println(err)
			continue
		}
		res, err := raNode.Eval(db)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if res == nil {
			continue
		}
		recs := res.GetRecords()

		for k, rec := range recs {
			fmt.Printf("row %v: %v\n", k, rec)
		}
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
