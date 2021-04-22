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
		givenCols      Cols
		wantedCols     Cols
	}{
		{
			name:           "test create table",
			givenDB:        db,
			givenTableName: "hoge",
			givenCols: Cols{
				{ColName{"hoge", "id"}, integer},
				{ColName{"hoge", "name"}, varchar},
			},
			wantedCols: Cols{
				{ColName{"hoge", "id"}, integer},
				{ColName{"hoge", "name"}, varchar},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db.CreateTable(tt.givenTableName, tt.givenCols)

			actualCols := db.Tables[tt.givenTableName].Cols

			if !actualCols.Equal(tt.wantedCols) {
				t.Errorf("expected %v, actual %v", tt.wantedCols, actualCols)
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
		name       string
		expected   Rows
		givenTable Table
		givenCols  Cols
	}{
		{
			name: "test project columns",
			givenTable: Table{
				Cols: Cols{
					{ColName{"hoge", "id"}, integer},
					{ColName{"hoge", "name"}, varchar},
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
			givenCols: Cols{
				{ColName{"hoge", "id"}, integer},
				{ColName{"hoge", "name"}, varchar},
				{ColName{"hoge", "id"}, integer},
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
