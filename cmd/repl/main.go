//+build

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/goropikari/psqlittle/backend"
	trans "github.com/goropikari/psqlittle/translator"
)

func main() {
	db, path := setupDB()
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
			// DDL
			writeLog(path, query)
			continue
		}
		recs := res.GetRecords()

		for k, rec := range recs {
			fmt.Printf("row %v: %v\n", k, rec)
		}
	}
}

func writeLog(path, query string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(query); err != nil {
		log.Println(err)
	}
}

func setupDB() (backend.DB, string) {
	path := getEnvWithDefault("DB_DATA_PATH", "data.db")

	db := backend.NewDatabase()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return db, path
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	ss := strings.Split(string(bytes), ";")
	for _, s := range ss {
		if strings.Trim(s, " \n") == "" {
			continue
		}
		raNode, _ := trans.NewPGTranslator(s).Translate()
		_, err := raNode.Eval(db)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return db, path
}

func getEnvWithDefault(key string, d string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return d
}
