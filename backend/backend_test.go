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
				Cols: core.Cols{
					{
						ColName: core.ColumnName{TableName: "hoge", Name: "id"},
						ColType: core.Integer,
					},
					{
						ColName: core.ColumnName{TableName: "hoge", Name: "name"},
						ColType: core.VarChar,
					},
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
			actual.Rows[0].Values[0] = "piyo"

			assert.NotEqual(t, tt.given, actual, "these two tables should be different")
		})
	}
}

func TestCreate(t *testing.T) {

	db := NewDatabase()

	tests := []struct {
		name           string
		givenDB        *Database
		givenTableName string
		givenCols      core.Cols
		wantedTable    *DBTable
	}{
		{
			name:           "test create table",
			givenDB:        db,
			givenTableName: "hoge",
			givenCols: core.Cols{
				{
					ColName: core.ColumnName{TableName: "hoge", Name: "id"},
					ColType: core.Integer,
				},
				{
					ColName: core.ColumnName{TableName: "hoge", Name: "name"},
					ColType: core.VarChar,
				},
			},
			wantedTable: &DBTable{
				Name: "hoge",
				ColNames: core.ColumnNames{
					{TableName: "hoge", Name: "id"},
					{TableName: "hoge", Name: "name"},
				},
				Cols: core.Cols{
					{
						ColName: core.ColumnName{TableName: "hoge", Name: "id"},
						ColType: core.Integer,
					},
					{
						ColName: core.ColumnName{TableName: "hoge", Name: "name"},
						ColType: core.VarChar,
					},
				},
				Rows: make(DBRows, 0),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db.CreateTable(tt.givenTableName, tt.givenCols)

			actualTable := db.Tables[tt.givenTableName]

			assert.Equal(t, tt.wantedTable, actualTable)
		})
	}
}

func TestInsert(t *testing.T) {

	table := &DBTable{
		Cols: core.Cols{
			{
				ColName: core.ColumnName{TableName: "hoge", Name: "id"},
				ColType: core.Integer,
			},
			{
				ColName: core.ColumnName{TableName: "hoge", Name: "name"},
				ColType: core.VarChar,
			},
		},
		ColNames: core.ColumnNames{
			{TableName: "hoge", Name: "id"},
			{TableName: "hoge", Name: "name"},
		},
		Rows: DBRows{},
	}

	var tests = []struct {
		name          string
		expected      *DBTable
		given         *DBTable
		givenColNames core.ColumnNames
		givenValsList core.ValuesList
	}{
		{
			name:  "test insert",
			given: table,
			expected: &DBTable{
				Cols: core.Cols{
					{
						ColName: core.ColumnName{TableName: "hoge", Name: "id"},
						ColType: core.Integer,
					},
					{
						ColName: core.ColumnName{TableName: "hoge", Name: "name"},
						ColType: core.VarChar,
					},
				},
				ColNames: core.ColumnNames{
					{TableName: "hoge", Name: "id"},
					{TableName: "hoge", Name: "name"},
				},
				Rows: DBRows{
					{
						ColNames: core.ColumnNames{
							{TableName: "hoge", Name: "id"},
							{TableName: "hoge", Name: "name"},
						},
						Values: core.Values{1, "taro"},
					},
					{
						ColNames: core.ColumnNames{
							{TableName: "hoge", Name: "id"},
							{TableName: "hoge", Name: "name"},
						},
						Values: core.Values{2, "hanako"},
					},
				},
			},
			givenColNames: core.ColumnNames{
				{TableName: "hoge", Name: "id"},
				{TableName: "hoge", Name: "name"},
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
			tt.given.InsertValues(tt.givenColNames, tt.givenValsList)

			assert.Equal(t, tt.expected, tt.given)
		})
	}
}

func TestProject(t *testing.T) {

	cn1 := core.ColumnName{TableName: "hoge", Name: "id"}
	cn2 := core.ColumnName{TableName: "hoge", Name: "name"}

	fn1 := func(row Row) (core.Value, error) {
		return row.GetValueByColName(cn1)
	}
	fn2 := func(row Row) (core.Value, error) {
		return row.GetValueByColName(cn2)
	}

	tests := []struct {
		name          string
		expected      DBRows
		givenTable    DBTable
		givenColNames core.ColumnNames
		givenFncs     []func(row Row) (core.Value, error)
	}{
		{
			name: "test project columns",
			givenTable: DBTable{
				ColNames: core.ColumnNames{
					{TableName: "hoge", Name: "id"},
					{TableName: "hoge", Name: "name"},
				},
				Rows: DBRows{
					{
						ColNames: core.ColumnNames{
							{TableName: "hoge", Name: "id"},
							{TableName: "hoge", Name: "name"},
						},
						Values: core.Values{1, "Hello"},
					},
					{
						ColNames: core.ColumnNames{
							{TableName: "hoge", Name: "id"},
							{TableName: "hoge", Name: "name"},
						},
						Values: core.Values{2, "World"},
					},
				},
			},
			givenColNames: core.ColumnNames{
				{TableName: "hoge", Name: "id"},
				{TableName: "hoge", Name: "name"},
				{TableName: "hoge", Name: "id"},
			},
			givenFncs: []func(Row) (core.Value, error){fn1, fn2},
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
			actual, err := tt.givenTable.Project(tt.givenColNames, tt.givenFncs)

			assert.NotEqual(t, tt.expected, actual)
			assert.NoError(t, err)
		})
	}
}
