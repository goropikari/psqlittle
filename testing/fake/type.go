package fake

import (
	"math/rand"

	"github.com/goropikari/mysqlite2/core"
)

// ColName generates fake ColName
func ColName() core.ColExpr {
	return core.ColExpr{
		TableName: RandString(),
		Name:      RandString(),
	}
}

// Value generates fake Value
func Value() core.Value {
	return RandString()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandString generates a random string
func RandString() string {
	n := 15 // length of random string
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
