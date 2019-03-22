package stub

import (
	"context"
	"testing"

	"github.com/ContainerSolutions/externalsecret-operator/pkg/apis/externalsecret-operator/v1alpha1"
	"github.com/ContainerSolutions/externalsecret-operator/pkg/secrets"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMakeSecret(t *testing.T) {

	suffix := "-value"
	key := "key"
	handler := Handler{}
	backend := secrets.NewDummySecretsBackend()
	backend.Init(suffix)
	var backendKey secrets.ContextKey = "backend"
	ctx := context.Background()
	ctx = context.WithValue(ctx, backendKey, backend)

	Convey("Given an ExternalSecret resource", t, func() {
		externalSecret := v1alpha1.ExternalSecret{
			Spec: v1alpha1.ExternalSecretSpec{
				Backend: "dummy",
				Key:     key,
			},
		}
		externalSecret.Name = "anExternalSecret"
		Convey("When creating a Secret", func() {
			theSecret, err := handler.makeSecret(ctx, &externalSecret)
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
