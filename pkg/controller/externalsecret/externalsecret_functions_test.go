package externalsecret

import (
	"testing"

	"github.com/ContainerSolutions/externalsecret-operator/pkg/apis/externalsecretoperator/v1alpha1"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSecretForCR(t *testing.T) {

	suffix := "-value"
	key := "key"

	// TODO: rely on single dummy secret backend until code migration is complete
	secretsBackend.Init(suffix)

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
