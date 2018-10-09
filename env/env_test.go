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

	vars := New().Current()

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

	vars := New().FromInput(input)

	assert.Equal(t, expected, vars)
}

func TestEnv_ToString(t *testing.T) {
	vars := []Variable{
		{
			Index:   1,
			Name:    "UNITY_HAS_2D_SUPPORT",
			NewName: "UNITY_HAS_2D_SUPPORT",
			Value:   "false",
			Comment: "comm",
			Deleted: false,
		},
		{
			Index:   2,
			Name:    "INSTANCE",
			NewName: "INSTANCE",
			Value:   "",
			Comment: "",
			Deleted: false,
		},
		{
			Index:   3,
			Name:    "GNOME_DESKTOP_SESSION_ID",
			NewName: "GNOME_DESKTOP_SESSION_ID",
			Value:   "this-is-deprecated",
			Comment: "",
			Deleted: false,
		},
		{
			Index:   5,
			Name:    "NEWVAR",
			NewName: "NEWVAR",
			Value:   "NEWVAL",
			Comment: "NEWCOMM",
			Deleted: true,
		},
	}

	expected := `UNITY_HAS_2D_SUPPORT=false # comm

INSTANCE=

GNOME_DESKTOP_SESSION_ID=this-is-deprecated

`

	out := New().ToString(vars)
	assert.Equal(t, expected, out)
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*Env)(nil), New())
}
