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
	secretKey := "secret"
	keyVersion := "latest"
	testSuffix := "test-suffix"
	expectedValue := secretKey + keyVersion + testSuffix

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
	})
}

func TestInit(t *testing.T) {
	params := make(map[string]string)

	params["Suffix"] = "dummy init"

	Convey("Should initialize backend", t, func() {
		backend := Backend{}
		err := backend.Init(params)
		So(err, ShouldBeNil)
	})

	Convey("Should fail initialize backend with invalid config", t, func() {
		backend := Backend{}
		err := backend.Init(make(map[string]string))
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "empty or invalid parameters")
	})

	Convey("Should fail initialize backend with invalid parameter suffix", t, func() {
		backend := Backend{}
		err := backend.Init(map[string]string{"uknown": "fail"})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "missing parameters")
	})
}
