package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {

	db := NewDB()

	tests := []struct {
		name           string
		givenDB        *DB
		givenTableName string
		givenColNames  ColNames
		wantedColNames ColNames
	}{
		{
			name:           "test create table",
			givenDB:        db,
			givenTableName: "hoge",
			givenColNames: ColNames{
				ColName{"hoge", "id", "int"},
				ColName{"hoge", "name", "varchar"},
			},
			wantedColNames: ColNames{
				ColName{"hoge", "id", "int"},
				ColName{"hoge", "name", "varchar"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db.CreateTable(tt.givenTableName, tt.givenColNames)

			actualColNames := db.Tables[tt.givenTableName].ColNames

			if !actualColNames.Equal(tt.wantedColNames) {
				t.Errorf("expected %v, actual %v", tt.wantedColNames, actualColNames)
			}
		})
	}
}

func TestInsert(t *testing.T) {

	table := Table{
		ColNames: []ColName{
			{"hoge", "id", "int"},
			{"hoge", "name", "varchar"},
		},
		Rows: []Row{},
	}

	var tests = []struct {
		name          string
		expected      Table
		given         Table
		givenColNames ColNames
		givenValsList ValuesList
	}{
		{
			name:  "test insert",
			given: table,
			expected: Table{
				ColNames: []ColName{
					{"hoge", "id", "int"},
					{"hoge", "name", "varchar"},
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
			givenColNames: []ColName{
				{"hoge", "id", "int"},
				{"hoge", "name", "varchar"},
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
			tt.given.Insert(tt.givenColNames, tt.givenValsList)

			if !tt.given.Equal(tt.expected) {
				t.Errorf("expected %s, actual %s", tt.expected, tt.given)
			}
		})
	}
}

func TestProject(t *testing.T) {

	tests := []struct {
		name       string
		expected   Rows
		givenTable Table
		givenCols  ColNames
	}{
		{
			name: "test project columns",
			givenTable: Table{
				ColNames: ColNames{
					{"hoge", "id", "int"},
					{"piyo", "name", "varchar"},
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
			givenCols: ColNames{
				{"hoge", "id", "int"},
				{"piyo", "name", "varchar"},
				{"hoge", "id", "int"},
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
			actual, err := tt.givenTable.Project(tt.givenCols)
			if actual.Equal(tt.expected) {
				t.Errorf("given(%v): expected %v, actual %v", tt.givenTable, tt.expected, actual)
			}
			assert.NoError(t, err)
		})
	}
}
