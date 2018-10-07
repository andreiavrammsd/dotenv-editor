package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var toLineTests = []struct {
	v        Variable
	expected string
}{
	{
		Variable{1, "name", "value", "", false},
		"name=value\n\n",
	},
	{
		Variable{2, "k", "v", "single line comment", false},
		"k=v # single line comment\n\n",
	},
	{
		Variable{3, "k", "   ", "multi \n\n line \n comment", false},
		"k=    # multi  line  comment\n\n",
	},
}

func TestVariable_ToLine(t *testing.T) {
	for _, test := range toLineTests {
		assert.Equal(t, test.expected, test.v.ToLine())
	}
}
