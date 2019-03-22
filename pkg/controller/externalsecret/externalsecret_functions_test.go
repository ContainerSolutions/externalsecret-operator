package externalsecret

import (
	"testing"

	"github.com/ContainerSolutions/externalsecret-operator/pkg/apis/externalsecretoperator/v1alpha1"
	"github.com/ContainerSolutions/externalsecret-operator/pkg/secrets"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSecretForCR(t *testing.T) {
	key := "key"
	suffix := "-value"

	dummy := secrets.NewDummySecretsBackend()
	dummy.Init("-value")
	secrets.BackendRegister("dummy", dummy)

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
