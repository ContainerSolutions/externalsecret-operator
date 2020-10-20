package backend

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBackendConfigFromJSON(t *testing.T) {
	Convey("Given a JSON backend config string", t, func() {
		configData := `{
			 "Type": "dummy",
			 "Parameters": {
         "Suffix": "-ohlord"
			  }
		}`

		Convey("When creating a Config object", func() {
			backendConfig, err := ConfigFromJSON(configData)
			So(err, ShouldBeNil)
			Convey("The data in Config is as expected", func() {
				So(backendConfig.Type, ShouldEqual, "dummy")
				So(backendConfig.Parameters, ShouldResemble, map[string]interface{}{"Suffix": "-ohlord"})
			})
		})

		Convey("When creating a Config object from invalid JSON", func() {

			_, err := ConfigFromJSON("")
			So(err, ShouldNotBeNil)

		})
	})
}

func TestBackendConfigFromCtrl(t *testing.T) {
	Convey("Given a JSON RawMessage backend config string", t, func() {
		configData := `{
			"type": "dummy",
			"auth": {
				"secretRef": {
					"name": "credential-secret",
					"namespace": "default"
				}
			},
			"parameters": {
				"Suffix": "I am definitely a param"
			}
		}`

		Convey("When creating a Config object", func() {
			backendConfig, err := ConfigFromCtrl([]byte(configData))
			So(err, ShouldBeNil)
			Convey("The data in Config is as expected", func() {
				So(backendConfig.Type, ShouldEqual, "dummy")
				So(backendConfig.Parameters, ShouldResemble, map[string]interface{}{"Suffix": "I am definitely a param"})
			})
		})

		Convey("When creating a Config object from invalid JSON RawMessage", func() {

			_, err := ConfigFromCtrl([]byte{})
			So(err, ShouldNotBeNil)

		})
	})
}

func TestConfigFromEnv(t *testing.T) {
	Convey("When backend config from env", t, func() {
		Convey("When creating a Config object from env", func() {
			value := `{
				"Type": "dummy",
				"Parameters": {
					"Suffix": "-ohlord"
				}
			}`
			key := "OPERATOR_CONFIG"

			os.Setenv(key, value)

			So(os.Getenv(key), ShouldNotBeBlank)

			backendConfig, err := ConfigFromEnv()
			So(err, ShouldBeNil)
			Convey("The data in Config is as expected", func() {
				So(backendConfig.Type, ShouldEqual, "dummy")
				So(backendConfig.Parameters, ShouldResemble, map[string]interface{}{"Suffix": "-ohlord"})
			})
		})

		Convey("When creating a Config object from a blank env val", func() {
			value := ""
			key := "OPERATOR_CONFIG"

			os.Setenv(key, value)

			So(os.Getenv(key), ShouldBeBlank)

			_, err := ConfigFromEnv()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "cannot find config: `OPERATOR_CONFIG` not set")
		})

		Convey("When OPERATOR_CONFIG is not set", func() {
			key := "OPERATOR_CONFIG"

			os.Unsetenv(key)

			So(os.Getenv(key), ShouldBeBlank)

			_, err := ConfigFromEnv()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "cannot find config: `OPERATOR_CONFIG` not set")
		})

	})
}
