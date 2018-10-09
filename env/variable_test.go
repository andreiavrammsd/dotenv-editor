package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var lineNewName = "new name"
var toLineTests = []struct {
	v        Variable
	expected string
}{
	{
		Variable{
			1,
			"name",
			&lineNewName,
			"value",
			"",
			nil,
		},
		"new name=value\n\n",
	},
	{
		Variable{
			2,
			"k",
			nil,
			"v",
			"single line comment",
			nil,
		},
		"k=v # single line comment\n\n",
	},
	{
		Variable{
			3,
			"k",
			nil,
			"   ",
			"multi \n\n line \n comment",
			nil,
		},
		"k=    # multi  line  comment\n\n",
	},
}

func TestVariable_ToLine(t *testing.T) {
	for _, test := range toLineTests {
		assert.Equal(t, test.expected, test.v.ToLine())
	}
}
