package apis

import (
	"k8s.io/apimachinery/pkg/runtime"

	externalsecret "github.com/containersolutions/externalsecret-operator/apis/secrets"
	secretstore "github.com/containersolutions/externalsecret-operator/apis/store"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes,
		secretstore.AddToScheme,
		externalsecret.AddToScheme,
	)
}

// AddToSchemes may be used to add all resources defined in the project to a Scheme
var AddToSchemes runtime.SchemeBuilder

// AddToScheme adds all Resources to the Scheme
func AddToScheme(s *runtime.Scheme) error {
	return AddToSchemes.AddToScheme(s)
}
