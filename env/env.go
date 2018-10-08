package env

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

const (
	pattern = `^([a-zA-Z_]{1,}[a-zA-Z0-9_]{0,})=([^#\n\r]*)#?(.*)`
	eol     = "\n"
)

// Env manages env variables as files (dotenv)
type Env interface {
	Current() []Variable
	FromInput(input string) []Variable
	Sync(src string, vars map[string]Variable) string
	ToFile(vars map[string]Variable) string
}

type env struct {
	reg *regexp.Regexp
}

// Current creates variables list based on current machin
func (e *env) Current() []Variable {
	current := os.Environ()
	return e.getVarsFromList(current)
}

// FromInput creates variables list from given string input
func (e *env) FromInput(input string) []Variable {
	data := strings.Split(input, eol)
	return e.getVarsFromList(data)
}

// Sync updates a dotenv file content with variables map
func (*env) Sync(src string, vars map[string]Variable) string {
	lines := strings.Split(src, eol)

	for i, line := range lines {
		for key, v := range vars {
			linePrefix := fmt.Sprintf("%s=", key)

			if strings.HasPrefix(line, linePrefix) {
				if v.Deleted {
					lines = append(lines[:i], lines[i+1:]...)
				} else {
					lines[i] = v.ToLine()
				}

				delete(vars, key)
				break
			}
		}
	}

	for _, v := range vars {
		lines = append(lines, v.ToLine())
	}

	return strings.Join(lines, eol)
}

// ToFile converts variables map to string
func (*env) ToFile(vars map[string]Variable) string {
	list := make(variablesList, 0, len(vars))

	for _, v := range vars {
		if !v.Deleted {
			list = append(list, v)
		}
	}

	sort.Sort(list)

	out := ""
	for _, v := range list {
		out += v.ToLine()
	}

	return out
}

func (e *env) getVarsFromList(list []string) []Variable {
	vars := make([]Variable, 0, len(list))
	index := 0

	for _, l := range list {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}

		data := e.reg.FindStringSubmatch(l)
		if len(data) != 4 {
			continue
		}

		index++

		v := Variable{}
		v.Index = index
		v.Name = data[1]
		v.Value = data[2]
		v.Comment = strings.TrimLeftFunc(data[3], func(r rune) bool {
			return r == ' ' || r == '#'
		})

		vars = append(vars, v)
	}

	return vars
}

// New sets up env package
func New() (Env, error) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &env{reg: reg}, nil
}
