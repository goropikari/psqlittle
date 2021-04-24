package translator_test

import (
	"testing"

	trans "github.com/goropikari/mysqlite2/translator"
)

func TestTranslate(t *testing.T) {
	var tests = []struct {
		name     string
		expected trans.RelationalAlgebraNode
		query    string
	}{
		{
			name:     "test translator",
			expected: nil,
			query:    "SELECT foo.id, foo.name FROM foo",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			transl := trans.NewPGTranslator(tt.query)
			actual, _ := transl.Translate()

			if actual != tt.expected {
				t.Errorf("expected %s, actual %s", tt.expected, actual)
			}
		})
	}
}
