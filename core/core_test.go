package core

import (
	"testing"
)

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
					ColName{"hoge", "id"},
					ColName{"piyo", "name"},
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
				ColName{"hoge", "id"},
				ColName{"piyo", "name"},
				ColName{"hoge", "id"},
			},
			expected: Rows{
				Row{
					Values: Values{"Hello", 1, "Hello"},
				},
				Row{
					Values: Values{"World", 1, "World"},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.givenTable.Project(tt.givenCols)
			if actual.Equal(tt.expected) {
				t.Errorf("given(%v): expected %v, actual %v", tt.givenTable, tt.expected, actual)
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
		givenSchema    TableSchema
		wantedSchema   TableSchema
	}{
		{
			name:           "test create table",
			givenDB:        db,
			givenTableName: "hoge",
			givenSchema: TableSchema{
				ColNames: ColNames{
					ColName{"hoge", "id"},
					ColName{"hoge", "name"},
				},
				ColTypes: ColTypes{
					ColType("int"),
					ColType("varchar"),
				},
			},
			wantedSchema: TableSchema{
				ColNames: ColNames{
					ColName{"hoge", "id"},
					ColName{"hoge", "name"},
				},
				ColTypes: ColTypes{
					ColType("int"),
					ColType("varchar"),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			db.CreateTable(tt.givenTableName, tt.givenSchema)

			actualSchema := db.Tables[tt.givenTableName].Schema

			if !actualSchema.Equal(tt.wantedSchema) {
				t.Errorf("expected %v, actual %v", tt.wantedSchema, actualSchema)
			}
		})
	}
}
