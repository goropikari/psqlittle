package backend

import (
	"testing"

	"github.com/goropikari/mysqlite2/core"
	"github.com/stretchr/testify/assert"
)

func TestTableCopy(t *testing.T) {

	var tests = []struct {
		name  string
		given Table
	}{
		{
			name: "test insert",
			given: Table{
				Cols: Cols{
					{core.ColName{"hoge", "id"}, core.Integer},
					{core.ColName{"hoge", "name"}, core.VarChar},
				},
				ColNameIndexes: ColNameIndexes{
					core.ColName{"hoge", "id"}:   0,
					core.ColName{"hoge", "name"}: 1,
				},
				Rows: []Row{
					{
						Values: core.Values{1, "Hello"},
					},
					{
						Values: core.Values{2, "World"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.given.Copy()
			actual.Cols[0].ColName.TableName = "piyo"
			actual.ColNameIndexes[core.ColName{"hoge", "id"}] = 1
			actual.Rows[0].Values[0] = "piyo"

			if !actual.NotEqual(tt.given) {
				t.Errorf("expected %v, actual %v", tt.given, actual)
			}
		})
	}
}

func TestCreate(t *testing.T) {

	db := NewDB()

	tests := []struct {
		name           string
		givenDB        *DB
		givenTableName string
		givenCols      Cols
		wantedTable    Table
	}{
		{
			name:           "test create table",
			givenDB:        db,
			givenTableName: "hoge",
			givenCols: Cols{
				{core.ColName{"hoge", "id"}, core.Integer},
				{core.ColName{"hoge", "name"}, core.VarChar},
			},
			wantedTable: Table{
				Cols: Cols{
					{core.ColName{"hoge", "id"}, core.Integer},
					{core.ColName{"hoge", "name"}, core.VarChar},
				},
				ColNameIndexes: ColNameIndexes{
					core.ColName{"hoge", "id"}:   0,
					core.ColName{"hoge", "name"}: 1,
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

	table := Table{
		Cols: Cols{
			{core.ColName{"hoge", "id"}, core.Integer},
			{core.ColName{"hoge", "name"}, core.VarChar},
		},
		ColNameIndexes: ColNameIndexes{
			core.ColName{"hoge", "id"}:   0,
			core.ColName{"hoge", "name"}: 1,
		},
		Rows: []Row{},
	}

	var tests = []struct {
		name          string
		expected      Table
		given         Table
		givenCols     Cols
		givenValsList core.ValuesList
	}{
		{
			name:  "test insert",
			given: table,
			expected: Table{
				Cols: Cols{
					{core.ColName{"hoge", "id"}, core.Integer},
					{core.ColName{"hoge", "name"}, core.VarChar},
				},
				ColNameIndexes: ColNameIndexes{
					core.ColName{"hoge", "id"}:   0,
					core.ColName{"hoge", "name"}: 1,
				},
				Rows: []Row{
					{
						core.Values{1, "taro"},
					},
					{
						core.Values{2, "hanako"},
					},
				},
			},
			givenCols: Cols{
				{core.ColName{"hoge", "id"}, core.Integer},
				{core.ColName{"hoge", "name"}, core.VarChar},
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
		expected      Rows
		givenTable    Table
		givenColNames core.ColNames
	}{
		{
			name: "test project columns",
			givenTable: Table{
				Cols: Cols{
					{core.ColName{"hoge", "id"}, core.Integer},
					{core.ColName{"hoge", "name"}, core.VarChar},
				},
				ColNameIndexes: ColNameIndexes{
					core.ColName{"hoge", "id"}:   0,
					core.ColName{"hoge", "name"}: 1,
				},
				Rows: Rows{
					Row{
						Values: core.Values{1, "Hello"},
					},
					Row{
						Values: core.Values{2, "World"},
					},
				},
			},
			givenColNames: core.ColNames{
				{"hoge", "id"},
				{"hoge", "name"},
				{"hoge", "id"},
			},
			expected: Rows{
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
