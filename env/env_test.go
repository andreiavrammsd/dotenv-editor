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

#test=asdasd
  
UNITY_HAS_2D_SUPPORT=false  # comm
INSTANCE=   

x
GNOME_DESKTOP_SESSION_ID=this-is-deprecated # 

	`

	expected := []Variable{
		{
			Index:   1,
			Name:    "UNITY_HAS_2D_SUPPORT",
			Value:   "false",
			Comment: "comm",
			Deleted: false,
		},
		{
			Index:   2,
			Name:    "INSTANCE",
			Value:   "",
			Comment: "",
			Deleted: false,
		},
		{
			Index:   3,
			Name:    "GNOME_DESKTOP_SESSION_ID",
			Value:   "this-is-deprecated",
			Comment: "",
			Deleted: false,
		},
	}

	e, _ := New()
	vars := e.FromInput(input)

	assert.Equal(t, expected, vars)
}

func TestNew(t *testing.T) {
	e, err := New()
	assert.Implements(t, (*Env)(nil), e)
	assert.NoError(t, err)
}

func TestNewWithError(t *testing.T) {
	linePattern = func() string {
		return "^([a-zA-Z_]"
	}
	e, err := New()
	assert.Nil(t, e)
	assert.Error(t, err)
}
