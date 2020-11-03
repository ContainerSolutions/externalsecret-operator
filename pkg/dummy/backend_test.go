package dummy

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewBackend(t *testing.T) {
	Convey("When creating a new dummy backend", t, func() {
		backend := NewBackend()
		So(backend, ShouldNotBeNil)
		So(backend, ShouldHaveSameTypeAs, &Backend{})
	})
}

func TestGet(t *testing.T) {
	var (
		secretKey     = "secret"
		keyVersion    = "latest"
		testSuffix    = "test-suffix"
		expectedValue = secretKey + keyVersion + testSuffix
	)

	Convey("Given an uninitialized dummy backend", t, func() {
		backend := Backend{}
		Convey("When retrieving a secret", func() {
			_, err := backend.Get(secretKey, keyVersion)
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "backend is not initialized")
			})
		})
	})

	Convey("Given an initialized dummy backend", t, func() {
		backend := Backend{}
		backend.suffix = testSuffix
		Convey("When retrieving a secret", func() {
			actualValue, err := backend.Get(secretKey, keyVersion)
			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
				So(actualValue, ShouldEqual, expectedValue)
			})
		})

		Convey("When retrieving secret details", func() {
			_, err := backend.Get("", "")
			Convey("An  error is returned when key is empty", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "empty key provided")

			})
		})

		Convey("When mock error key is provided", func() {
			_, err := backend.Get("ErroredKey", "")
			Convey("An  error is returned when key is a mock error key", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Mocked error")

			})
		})
	})
}

func TestInit(t *testing.T) {
	var (
		params      = make(map[string]interface{})
		credentials = []byte{}
	)

	params["Suffix"] = "dummy init"
	Convey("Should initialize backend", t, func() {
		backend := Backend{}
		credentials = []byte{}
		err := backend.Init(params, credentials)
		So(err, ShouldBeNil)
	})

	Convey("Should fail initialize backend with invalid config", t, func() {
		backend := Backend{}
		err := backend.Init(make(map[string]interface{}), credentials)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "empty or invalid parameters")
	})

	Convey("Should fail initialize backend with invalid parameter suffix", t, func() {
		backend := Backend{}
		err := backend.Init(map[string]interface{}{"unknown": "fail"}, credentials)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "missing parameters")
	})
}
