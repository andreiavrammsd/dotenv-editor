package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv_Current(t *testing.T) {
	os.Clearenv()

	expected := []Variable{
		{
			Index:   1,
			Name:    "KEY",
			Value:   "VALUE",
			Deleted: false,
		},
		{
			Index:   2,
			Name:    "KEY2",
			Value:   "VALUE2",
			Deleted: false,
		},
	}
	for _, v := range expected {
		os.Setenv(v.Name, v.Value)
	}

	e, _ := New()
	vars := e.Current()

	assert.Equal(t, expected, vars)
}

func TestEnv_FromInput(t *testing.T) {
	input := `
		KEY=VALUE
		KEY2=VALUE2
	`

	expected := []Variable{
		{
			Index:   1,
			Name:    "KEY",
			Value:   "VALUE",
			Deleted: false,
		},
		{
			Index:   2,
			Name:    "KEY2",
			Value:   "VALUE2",
			Deleted: false,
		},
	}

	e, _ := New()
	vars := e.FromInput(input)

	assert.Equal(t, expected, vars)
}
