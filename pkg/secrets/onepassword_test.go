package secrets

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetOnePassword(t *testing.T) {
	fmt.Println("Hello")

	secretKey := "secret"
	secretValue := "secretValue"
	expectedValue := secretValue

	Convey("Given an initialized OnePasswordBackend", t, func() {
		backend := NewOnePasswordBackend()

		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(secretKey)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, expectedValue)
			})
		})
	})
}
