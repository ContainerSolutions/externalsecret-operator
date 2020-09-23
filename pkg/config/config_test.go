package backend

import (
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
				So(backendConfig.Parameters, ShouldResemble, map[string]string{"Suffix": "-ohlord"})
			})
		})
	})
}
