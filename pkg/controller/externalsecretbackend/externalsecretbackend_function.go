package externalsecretbackend

import (
	externalsecretoperatorv1alpha1 "github.com/ContainerSolutions/externalsecret-operator/pkg/apis/externalsecretoperator/v1alpha1"
	"github.com/ContainerSolutions/externalsecret-operator/pkg/secrets"
)

func newBackendInstanceForCR(cr *externalsecretoperatorv1alpha1.ExternalSecretBackend) error {
	return secrets.BackendInstantiate(cr.Name, cr.Spec.Type)
}

func initBackendInstanceForCR(cr *externalsecretoperatorv1alpha1.ExternalSecretBackend) error {
	return secrets.BackendInstances[cr.Name].Init(cr.Spec.Parameters)
}
