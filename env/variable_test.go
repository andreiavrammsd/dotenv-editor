package env

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
		"new name=value",
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
		"k=v # single line comment",
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
		"k=    # multi  line  comment",
	},
}

func TestVariable_ToLine(t *testing.T) {
	Convey("Given a variable", t, func() {
		Convey("It should be converted to its string representation", func() {
			for _, test := range toLineTests {
				line := test.v.ToString()
				So(line, ShouldEqual, test.expected)
			}
		})
	})
}
