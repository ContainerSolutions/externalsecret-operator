package onepassword

import (
	"fmt"
	"testing"

	"github.com/ContainerSolutions/externalsecret-operator/secrets"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOnePasswordBackend(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	Convey("Given an initialized OnePasswordBackend", t, func() {
		key := "testkey"
		expectedValue := "testvalue"

		err := secrets.InitFromEnv()
		if err != nil {
			fmt.Println("Init: Error parsing the OPERATOR_CONFIG env var. " + err.Error())
			t.Fail()
		}

		backend := secrets.Instances["onepassword"]

		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(key)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, expectedValue)
			})
		})
	})
}
