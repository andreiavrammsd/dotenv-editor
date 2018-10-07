package env

import (
	"fmt"
	"strings"
)

type Variable struct {
	Index   int    `json:"index"`
	Name    string `json:"name"`
	Value   string `json:"value"`
	Comment string `json:"comment"`
	Deleted bool   `json:"deleted"`
}

type variablesList []Variable

// Len is part of sort.Interface.
func (d variablesList) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d variablesList) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d variablesList) Less(i, j int) bool {
	return d[i].Index < d[j].Index
}

func (v Variable) ToLine() (line string) {
	c := strings.TrimSpace(v.Comment)
	if c != "" {
		c = fmt.Sprintf(" # %s", strings.Replace(c, "\n", "", -1))
	}

	line += fmt.Sprintf("%s=%s%s\n\n", v.Name, v.Value, c)

	return
}
