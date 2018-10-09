package env

import (
	"fmt"
	"strings"
)

// Variable model
type Variable struct {
	Index   int    `json:"index"`
	Name    string `json:"name"`
	NewName string `json:"new_name,omitempty"`
	Value   string `json:"value"`
	Comment string `json:"comment"`
	Deleted bool   `json:"deleted"`
}

// ToLine generates the string representation of a variable in a file
func (v Variable) ToLine() string {
	c := strings.TrimSpace(v.Comment)
	if c != "" {
		c = fmt.Sprintf(" # %s", strings.Replace(c, "\n", "", -1))
	}

	name := v.Name
	if v.NewName != "" {
		name = v.NewName
	}

	return fmt.Sprintf("%s=%s%s\n\n", name, v.Value, c)
}
