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
					ColName("hoge"),
					ColName("piyo"),
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
				ColName("piyo"),
				ColName("hoge"),
				ColName("piyo"),
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
