package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/andreiavrammsd/dotenv-editor/env"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

// Env package mock
type EnvMock struct {
	mock.Mock
}

func (m *EnvMock) Current() []env.Variable {
	args := m.Called()
	return args.Get(0).([]env.Variable)
}

func (m *EnvMock) FromInput(input string) []env.Variable {
	return []env.Variable{}
}

func (m *EnvMock) Sync(src string, vars []env.Variable) string {
	return ""
}

func (m *EnvMock) ToString(vars []env.Variable) string {
	return ""
}

// HTTP writer mock
type WriterMock struct {
	mock.Mock
	Body       []byte
	StatusCode int
}

func (w *WriterMock) Header() http.Header {
	args := w.Called()
	return args.Get(0).(http.Header)
}

func (w *WriterMock) Write(data []byte) (int, error) {
	w.Body = data
	args := w.Called()
	return args.Int(0), args.Error(1)
}

func (w *WriterMock) WriteHeader(statusCode int) {
}

func TestHandlers_GetCurrent(t *testing.T) {
	vars := []env.Variable{
		{
			Index:   1,
			Name:    "USER",
			NewName: nil,
			Value:   "MSD",
			Comment: " # auth",
			Deleted: nil,
		},
		{
			Index:   2,
			Name:    "ENV",
			NewName: nil,
			Value:   "uat",
			Comment: "",
			Deleted: nil,
		},
	}
	expectedBody, _ := json.Marshal(vars)
	expectedContentType := "application/json"

	Convey("Given an HTTP handler for /env/current", t, func() {
		writer := &WriterMock{}
		request := &http.Request{}
		environment := &EnvMock{}

		handler := New(environment)

		Convey("The environment variables list should be generated", func() {
			environment.On("Current").Once().Return(vars)

			Convey("Content-Type header should be set", func() {
				header := http.Header{}
				writer.On("Header").Once().Return(header)

				Convey("The variables list should be written to the response", func() {
					writer.On("Write").Once().Return(113, nil)

					handler.GetCurrent(writer, request)
					environment.AssertExpectations(t)
					writer.AssertExpectations(t)

					Convey("The body should be the variables list as JSON", func() {
						So(string(writer.Body), ShouldEqual, string(expectedBody))
					})
					Convey("With Content-Type application/json", func() {
						So(header.Get("Content-Type"), ShouldEqual, expectedContentType)
					})
				})
			})

		})
	})
}
