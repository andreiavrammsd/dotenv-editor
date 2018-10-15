package env

import (
	"fmt"
	"os"
	"regexp"
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
	Sync(src string, vars []Variable) string
	ToString(vars []Variable) string
}

type env struct {
	reg *regexp.Regexp
}

// Current creates variables list based on current machine
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
func (*env) Sync(src string, vars []Variable) string {
	lines := strings.Split(src, eol)

	for i, line := range lines {
		for j, v := range vars {
			linePrefix := fmt.Sprintf("%s=", v.Name)

			if strings.HasPrefix(line, linePrefix) {
				if v.IsDeleted() {
					lines = append(lines[:i], lines[i+1:]...)
				} else {
					lines[i] = v.ToString()
				}

				vars = append(vars[:j], vars[j+1:]...)
				break
			}
		}
	}

	for _, v := range vars {
		lines = append(lines, v.ToString())
	}

	return strings.Join(lines, eol)
}

// ToString converts variables list to string
func (*env) ToString(vars []Variable) (out string) {
	for _, v := range vars {
		if !v.IsDeleted() {
			out += fmt.Sprintf("%s%s", v.ToString(), eol)
		}
	}
	return
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
		v.Value = strings.TrimSpace(data[2])
		v.Comment = strings.TrimLeftFunc(data[3], func(r rune) bool {
			return r == ' ' || r == '#'
		})

		vars = append(vars, v)
	}

	return vars
}

// New sets up env package
func New() Env {
	reg := regexp.MustCompile(pattern)
	return &env{reg: reg}
}
