package backend

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type MockBackend struct {
	Param1 string
}

func NewBackend() Backend {
	return &MockBackend{}
}

func (m *MockBackend) Init(params map[string]string) error {
	m.Param1 = params["Param1"]
	return nil
}

func (m *MockBackend) Get(key string) (string, error) {
	return m.Param1, nil
}

func TestRegister(t *testing.T) {
	Convey("Given a mocked backend", t, func() {
		Convey("When registering it as a backend type", func() {
			Register("mock", NewBackend)
			Convey("Then the instantiation function is registered with the correct label", func() {
				function, found := Functions["mock"]
				So(found, ShouldBeTrue)
				So(function, ShouldEqual, NewBackend)
			})
		})
	})
}

func TestInstantiate(t *testing.T) {
	Convey("Given a registered backend type", t, func() {
		Register("mock", NewBackend)
		Convey("When Instantiating it using the right label", func() {
			err := Instantiate("mock-backend", "mock")
			So(err, ShouldBeNil)
			Convey("Then a backend is instantiated with the right label", func() {
				backend, found := Instances["mock-backend"]
				So(found, ShouldBeTrue)
				So(reflect.TypeOf(backend), ShouldEqual, reflect.TypeOf(&MockBackend{}))
			})
		})
		Convey("When Instantiating it using the wrong label", func() {
			err := Instantiate("mock-backend", "mock-wrong-label")
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown backend type: 'mock-wrong-label'")
			})
		})
	})
}

func TestInitFromEnv(t *testing.T) {

	configStruct := Config{
		Type: "mock",
		Parameters: map[string]string{
			"Param1": "Value1",
		},
	}

	Convey("Given a registered backend type", t, func() {
		Register("mock", NewBackend)
		Convey("Given a valid config", func() {
			configData, _ := json.Marshal(configStruct)
			os.Setenv("OPERATOR_CONFIG", string(configData))
			os.Setenv("OPERATOR_NAME", "mock-backend")
			Convey("When initializing backend from env", func() {
				err := InitFromEnv()
				So(err, ShouldBeNil)
				Convey("Then a backend is instantiated and initialized correctly", func() {
					backend, found := Instances["mock-backend"]
					So(found, ShouldBeTrue)
					So(reflect.TypeOf(backend), ShouldEqual, reflect.TypeOf(&MockBackend{}))
					value, _ := backend.Get("")
					So(value, ShouldEqual, "Value1")
				})
			})
		})

		Convey("Given a valid config but no OPERATOR_NAME", func() {
			configData, _ := json.Marshal(configStruct)
			os.Setenv("OPERATOR_CONFIG", string(configData))
			os.Unsetenv("OPERATOR_NAME")
			Convey("When initializing backend from env", func() {
				err := InitFromEnv()
				So(err, ShouldNotBeNil)
				Convey("Then an error message is returned", func() {
					So(err.Error(), ShouldStartWith, "OPERATOR_NAME must be set")
				})
			})
		})

		Convey("Given a valid config with unknown backend type ", func() {
			configStruct.Type = "unknown"
			configData, _ := json.Marshal(configStruct)
			os.Setenv("OPERATOR_CONFIG", string(configData))
			os.Setenv("OPERATOR_NAME", "mock-backend")
			Convey("When initializing backend from env", func() {
				err := InitFromEnv()
				So(err, ShouldNotBeNil)
				Convey("Then an error message is returned", func() {
					So(err.Error(), ShouldEqual, "unknown backend type: 'unknown'")
				})
			})
		})

		Convey("Given an invalid config", func() {
			os.Setenv("OPERATOR_CONFIG", "garbage")
			Convey("When initializing backend from env", func() {
				err := InitFromEnv()
				So(err, ShouldNotBeNil)
				Convey("Then an error is returned", func() {
					So(err.Error(), ShouldStartWith, "invalid")
				})
			})
		})

		Convey("Given a missing config", func() {
			os.Unsetenv("OPERATOR_CONFIG")
			Convey("When initializing backend from env", func() {
				err := InitFromEnv()
				So(err, ShouldNotBeNil)
				Convey("Then an error is returned", func() {
					So(err.Error(), ShouldStartWith, "cannot find config")
				})
			})
		})
	})
}
