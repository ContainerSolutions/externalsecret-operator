package externalsecret

import (
	"testing"

	"github.com/containersolutions/externalsecretoperator/pkg/apis/externalsecretoperator/v1alpha1"
	"github.com/containersolutions/externalsecretoperator/secrets/backend"
	"github.com/containersolutions/externalsecretoperator/secrets/dummy"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSecretForCR(t *testing.T) {
	key := "key"
	suffix := "-value"

	backend.Register("dummy", dummy.NewBackend)
	backend.Instantiate("dummy", "dummy")
	backend.Instances["dummy"].Init(map[string]string{"suffix": "-value"})

	Convey("Given an ExternalSecret resource", t, func() {
		externalSecret := v1alpha1.ExternalSecret{
			Spec: v1alpha1.ExternalSecretSpec{
				Backend: "dummy",
				Key:     key,
			},
		}
		externalSecret.Name = "anExternalSecret"
		Convey("When creating a Secret", func() {
			theSecret, err := newSecretForCR(&externalSecret)
			Convey("The Secret should have the correct key", func() {
				So(err, ShouldBeNil)
				So(theSecret.Data, ShouldContainKey, key)
			})
			Convey("The Secret should have the correct value", func() {
				So(string(theSecret.Data[key]), ShouldEqual, key+suffix)
			})
		})
	})
}
