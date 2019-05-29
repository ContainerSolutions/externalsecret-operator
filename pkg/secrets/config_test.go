package secrets

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBackendConfigFromJSON(t *testing.T) {
	Convey("Given a JSON backend config string", t, func() {
		configData := `{
			 "Name": "dummy-example",
			 "Type": "dummy",
			 "Parameters": {
         "Suffix": "-ohlord"
			  }
		}`

		Convey("When creating a BackendConfig object", func() {
			backendConfig, err := BackendConfigFromJSON(configData)
			So(err, ShouldBeNil)
			Convey("The data in BackendConfig is as expected", func() {
				So(backendConfig.Name, ShouldEqual, "dummy-example")
				So(backendConfig.Type, ShouldEqual, "dummy")
				So(backendConfig.Parameters, ShouldResemble, map[string]string{"Suffix": "-ohlord"})
			})
		})
	})
}
