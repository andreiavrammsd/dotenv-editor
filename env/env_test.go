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
			NewName: nil,
			Value:   "VALUE",
			Deleted: nil,
		},
		{
			Index:   2,
			Name:    "KEY2",
			NewName: nil,
			Value:   "VALUE2",
			Deleted: nil,
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
			NewName: nil,
			Value:   "false",
			Comment: "comm",
			Deleted: nil,
		},
		{
			Index:   2,
			Name:    "INSTANCE",
			NewName: nil,
			Value:   "",
			Comment: "",
			Deleted: nil,
		},
		{
			Index:   3,
			Name:    "GNOME_DESKTOP_SESSION_ID",
			NewName: nil,
			Value:   "this-is-deprecated",
			Comment: "",
			Deleted: nil,
		},
	}

	vars := New().FromInput(input)

	assert.Equal(t, expected, vars)
}

func TestEnv_ToString(t *testing.T) {
	newName1 := "UNITY_HAS_2D_SUPPORT"
	newName2 := "INSTANCE"
	newName3 := "GNOME_DESKTOP_SESSION_ID"
	newName4 := "NEWVAR"

	notDeleted := false
	deleted := true

	vars := []Variable{
		{
			Index:   1,
			Name:    "UNITY_HAS_2D_SUPPORT",
			NewName: &newName1,
			Value:   "false",
			Comment: "comm",
			Deleted: &notDeleted,
		},
		{
			Index:   2,
			Name:    "INSTANCE",
			NewName: &newName2,
			Value:   "",
			Comment: "",
			Deleted: &notDeleted,
		},
		{
			Index:   3,
			Name:    "GNOME_DESKTOP_SESSION_ID",
			NewName: &newName3,
			Value:   "this-is-deprecated",
			Comment: "",
			Deleted: nil,
		},
		{
			Index:   5,
			Name:    "NEWVAR",
			NewName: &newName4,
			Value:   "NEWVAL",
			Comment: "NEWCOMM",
			Deleted: &deleted,
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
