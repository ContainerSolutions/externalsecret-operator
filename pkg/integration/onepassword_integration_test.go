package integration

import (
	"fmt"
	"testing"

	"github.com/ContainerSolutions/externalsecret-operator/pkg/secrets"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOnePasswordBackend(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	Convey("Given an initialized OnePasswordBackend", t, func() {
		vault := "Personal"
		key := "testkey"
		expectedValue := "testvalue"

		client := secrets.OnePasswordCliClient{}
		backend := secrets.NewOnePasswordBackend(vault, client)

		err := backend.Init(vault)
		if err != nil {
			fmt.Println("Init: " + err.Error())
		}

		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(key)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, expectedValue)
			})
		})
	})
}
