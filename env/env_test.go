package env

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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

	Convey("Given the current machine environment", t, func() {
		Convey("The variables list should be extracted", func() {
			vars := New().Current()
			So(vars, ShouldResemble, expected)
		})
	})
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

	Convey("Given a string input", t, func() {
		Convey("The variables list should be extracted", func() {
			vars := New().FromInput(input)
			So(vars, ShouldResemble, expected)
		})
	})
}

func TestEnv_Sync(t *testing.T) {
	src := `
UNITY_HAS_2D_SUPPORT=false # comm
INSTANCE=

TDB=TDB

# Some inner text

GNOME_DESKTOP_SESSION_ID=this-is-deprecated

`

	newName1 := "UNITY_HAS_3D_SUPPORT"
	newName2 := "NEW_INSTANCE"
	newName3 := "GNOME_DESKTOP_SESSION_ID"

	notDeleted := false
	deleted := true

	vars := []Variable{
		{
			Index:   1,
			Name:    "UNITY_HAS_2D_SUPPORT",
			NewName: &newName1,
			Value:   "true",
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
			Value:   "this is deprecated",
			Comment: "updated",
			Deleted: nil,
		},
		{
			Index:   4,
			Name:    "TDB",
			NewName: nil,
			Value:   "TBD",
			Comment: "TBD",
			Deleted: &deleted,
		},
		{
			Index:   5,
			Name:    "ADD_VAR",
			NewName: nil,
			Value:   "-9",
			Comment: "just added",
			Deleted: nil,
		},
	}

	expected := `
UNITY_HAS_3D_SUPPORT=true # comm
NEW_INSTANCE=


# Some inner text

GNOME_DESKTOP_SESSION_ID=this is deprecated # updated


ADD_VAR=-9 # just added`

	Convey("Given a string of variables and a list of variables", t, func() {
		Convey("The string should be synced with the variables", func() {
			result := New().Sync(src, vars)
			So(result, ShouldEqual, expected)
		})
	})
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
	Convey("Given a list of variables", t, func() {
		Convey("They should be converted to a string list", func() {
			out := New().ToString(vars)
			So(out, ShouldEqual, expected)
		})
	})
}

func TestNew(t *testing.T) {
	Convey("Given a new instance of package", t, func() {
		e := New()

		Convey("It should implement the Env interface", func() {
			So(e, ShouldImplement, (*Env)(nil))
		})
	})
}
