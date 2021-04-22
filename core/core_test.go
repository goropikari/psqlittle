package core

import (
	"testing"

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
					{ColName{"hoge", "id"}, integer},
					{ColName{"hoge", "name"}, varchar},
				},
				ColNameIndexes: ColNameIndexes{
					ColName{"hoge", "id"}:   0,
					ColName{"hoge", "name"}: 1,
				},
				Rows: []Row{
					{
						Values: Values{1, "Hello"},
					},
					{
						Values: Values{2, "World"},
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
			actual.ColNameIndexes[ColName{"hoge", "id"}] = 1
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
				{ColName{"hoge", "id"}, integer},
				{ColName{"hoge", "name"}, varchar},
			},
			wantedTable: Table{
				Cols: Cols{
					{ColName{"hoge", "id"}, integer},
					{ColName{"hoge", "name"}, varchar},
				},
				ColNameIndexes: ColNameIndexes{
					ColName{"hoge", "id"}:   0,
					ColName{"hoge", "name"}: 1,
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
			{ColName{"hoge", "id"}, integer},
			{ColName{"hoge", "name"}, varchar},
		},
		ColNameIndexes: ColNameIndexes{
			ColName{"hoge", "id"}:   0,
			ColName{"hoge", "name"}: 1,
		},
		Rows: []Row{},
	}

	var tests = []struct {
		name          string
		expected      Table
		given         Table
		givenCols     Cols
		givenValsList ValuesList
	}{
		{
			name:  "test insert",
			given: table,
			expected: Table{
				Cols: Cols{
					{ColName{"hoge", "id"}, integer},
					{ColName{"hoge", "name"}, varchar},
				},
				ColNameIndexes: ColNameIndexes{
					ColName{"hoge", "id"}:   0,
					ColName{"hoge", "name"}: 1,
				},
				Rows: []Row{
					{
						Values{1, "taro"},
					},
					{
						Values{2, "hanako"},
					},
				},
			},
			givenCols: Cols{
				{ColName{"hoge", "id"}, integer},
				{ColName{"hoge", "name"}, varchar},
			},
			givenValsList: []Values{
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
		givenColNames ColNames
	}{
		{
			name: "test project columns",
			givenTable: Table{
				Cols: Cols{
					{ColName{"hoge", "id"}, integer},
					{ColName{"hoge", "name"}, varchar},
				},
				ColNameIndexes: ColNameIndexes{
					ColName{"hoge", "id"}:   0,
					ColName{"hoge", "name"}: 1,
				},
				Rows: Rows{
					Row{
						Values: Values{1, "Hello"},
					},
					Row{
						Values: Values{2, "World"},
					},
				},
			},
			givenColNames: ColNames{
				{"hoge", "id"},
				{"hoge", "name"},
				{"hoge", "id"},
			},
			expected: Rows{
				{
					Values: Values{"Hello", 1, "Hello"},
				},
				{
					Values: Values{"World", 1, "World"},
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
