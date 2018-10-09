package env

import (
	"fmt"
	"strings"
)

// Variable model
type Variable struct {
	Index   int     `json:"index"`
	Name    string  `json:"name"`
	NewName *string `json:"new_name,omitempty"`
	Value   string  `json:"value"`
	Comment string  `json:"comment"`
	Deleted *bool   `json:"deleted,omitempty"`
}

// IsDeleted indicates the variable is marked to be deleted
func (v Variable) IsDeleted() bool {
	return v.Deleted != nil && *v.Deleted
}

// IsRenamed indicates the variable was given a new name
func (v Variable) IsRenamed() bool {
	return v.NewName != nil && *v.NewName != ""
}

// ToString generates the string representation of a variable in a file
func (v Variable) ToString() string {
	c := strings.TrimSpace(v.Comment)
	if c != "" {
		c = fmt.Sprintf(" # %s", strings.Replace(c, eol, "", -1))
	}

	name := v.Name
	if v.IsRenamed() {
		name = *v.NewName
	}

	return fmt.Sprintf("%s=%s%s", name, v.Value, c)
}
