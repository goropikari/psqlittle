package backend

import (
	"testing"

	"github.com/goropikari/mysqlite2/core"
	"github.com/stretchr/testify/assert"
)

func TestTableCopy(t *testing.T) {

	var tests = []struct {
		name  string
		given DBTable
	}{
		{
			name: "test insert",
			given: DBTable{
				Cols: Cols{
					{core.ColName{TableName: "hoge", Name: "id"}, core.Integer},
					{core.ColName{TableName: "hoge", Name: "name"}, core.VarChar},
				},
				ColNameIndexes: ColNameIndexes{
					core.ColName{TableName: "hoge", Name: "id"}:   0,
					core.ColName{TableName: "hoge", Name: "name"}: 1,
				},
				Rows: DBRows{
					&DBRow{
						Values: core.Values{1, "Hello"},
					},
					&DBRow{
						Values: core.Values{2, "World"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.given.Copy().(*DBTable) // convert Table to *DBTable
			actual.Cols[0].ColName.TableName = "piyo"
			actual.ColNameIndexes[core.ColName{TableName: "hoge", Name: "id"}] = 1
			actual.Rows[0].Values[0] = "piyo"

			if !actual.NotEqual(tt.given) {
				t.Errorf("expected %v, actual %v", tt.given, actual)
			}
		})
	}
}

func TestCreate(t *testing.T) {

	db := NewDatabase()

	tests := []struct {
		name           string
		givenDB        *Database
		givenTableName string
		givenCols      Cols
		wantedTable    DBTable
	}{
		{
			name:           "test create table",
			givenDB:        db,
			givenTableName: "hoge",
			givenCols: Cols{
				{core.ColName{TableName: "hoge", Name: "id"}, core.Integer},
				{core.ColName{TableName: "hoge", Name: "name"}, core.VarChar},
			},
			wantedTable: DBTable{
				Cols: Cols{
					{core.ColName{TableName: "hoge", Name: "id"}, core.Integer},
					{core.ColName{TableName: "hoge", Name: "name"}, core.VarChar},
				},
				ColNameIndexes: ColNameIndexes{
					core.ColName{TableName: "hoge", Name: "id"}:   0,
					core.ColName{TableName: "hoge", Name: "name"}: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db.CreateTable(tt.givenTableName, tt.givenCols)

			actualTable := db.Tables[tt.givenTableName]

			if !actualTable.Equal(tt.wantedTable) {
				t.Errorf("expected %v, actual %v", tt.wantedTable, actualTable)
			}
		})
	}
}

func TestInsert(t *testing.T) {

	table := DBTable{
		Cols: Cols{
			{core.ColName{TableName: "hoge", Name: "id"}, core.Integer},
			{core.ColName{TableName: "hoge", Name: "name"}, core.VarChar},
		},
		ColNameIndexes: ColNameIndexes{
			core.ColName{TableName: "hoge", Name: "id"}:   0,
			core.ColName{TableName: "hoge", Name: "name"}: 1,
		},
		Rows: DBRows{},
	}

	var tests = []struct {
		name          string
		expected      DBTable
		given         DBTable
		givenCols     Cols
		givenValsList core.ValuesList
	}{
		{
			name:  "test insert",
			given: table,
			expected: DBTable{
				Cols: Cols{
					{core.ColName{TableName: "hoge", Name: "id"}, core.Integer},
					{core.ColName{TableName: "hoge", Name: "name"}, core.VarChar},
				},
				ColNameIndexes: ColNameIndexes{
					core.ColName{TableName: "hoge", Name: "id"}:   0,
					core.ColName{TableName: "hoge", Name: "name"}: 1,
				},
				Rows: DBRows{
					{
						core.Values{1, "taro"},
					},
					{
						core.Values{2, "hanako"},
					},
				},
			},
			givenCols: Cols{
				{core.ColName{TableName: "hoge", Name: "id"}, core.Integer},
				{core.ColName{TableName: "hoge", Name: "name"}, core.VarChar},
			},
			givenValsList: []core.Values{
				{
					1,
					"taro",
				},
				{
					2,
					"hanako",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.given.Insert(tt.givenCols, tt.givenValsList)

			if !tt.given.Equal(tt.expected) {
				t.Errorf("expected %v, actual %v", tt.expected, tt.given)
			}
		})
	}
}

func TestProject(t *testing.T) {

	tests := []struct {
		name          string
		expected      DBRows
		givenTable    DBTable
		givenColNames core.ColNames
	}{
		{
			name: "test project columns",
			givenTable: DBTable{
				Cols: Cols{
					{core.ColName{TableName: "hoge", Name: "id"}, core.Integer},
					{core.ColName{TableName: "hoge", Name: "name"}, core.VarChar},
				},
				ColNameIndexes: ColNameIndexes{
					core.ColName{TableName: "hoge", Name: "id"}:   0,
					core.ColName{TableName: "hoge", Name: "name"}: 1,
				},
				Rows: DBRows{
					{
						Values: core.Values{1, "Hello"},
					},
					{
						Values: core.Values{2, "World"},
					},
				},
			},
			givenColNames: core.ColNames{
				{TableName: "hoge", Name: "id"},
				{TableName: "hoge", Name: "name"},
				{TableName: "hoge", Name: "id"},
			},
			expected: DBRows{
				{
					Values: core.Values{"Hello", 1, "Hello"},
				},
				{
					Values: core.Values{"World", 1, "World"},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.givenTable.Project(tt.givenColNames)
			if actual.Equal(tt.expected) {
				t.Errorf("given(%v): expected %v, actual %v", tt.givenTable, tt.expected, actual)
			}
			assert.NoError(t, err)
		})
	}
}