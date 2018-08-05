package stub

import (
	"context"
	"testing"

	"github.com/ContainerSolutions/externalconfig-operator/pkg/apis/externalconfig-operator/v1alpha1"
	"github.com/ContainerSolutions/externalconfig-operator/pkg/secrets"
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

	Convey("Given an ExternalConfig resource", t, func() {
		externalConfig := v1alpha1.ExternalConfig{
			Spec: v1alpha1.ExternalConfigSpec{
				Backend: "dummy",
				Key:     key,
			},
		}
		externalConfig.Name = "anExternalConfig"
		Convey("When creating a Secret", func() {
			theSecret, err := handler.makeSecret(ctx, &externalConfig)
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
