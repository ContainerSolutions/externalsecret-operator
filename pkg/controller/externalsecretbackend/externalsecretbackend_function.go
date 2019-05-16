package externalsecretbackend

import (
	"fmt"

	externalsecretoperatorv1alpha1 "github.com/ContainerSolutions/externalsecret-operator/pkg/apis/externalsecretoperator/v1alpha1"
	"github.com/ContainerSolutions/externalsecret-operator/pkg/secrets"
)

func newExternalSecretBackendForCR(cr *externalsecretoperatorv1alpha1.ExternalSecretBackend) error {
	var backend secrets.BackendIface

	switch cr.Spec.Type {
	case "asm":
		backend = secrets.NewAWSSecretsManagerBackend()
	case "dummy":
		backend = secrets.NewDummySecretsBackend()
	default:
		return fmt.Errorf("Backend registration failed: unkown backend type `%v`", cr.Spec.Type)
	}

	secrets.BackendRegister(cr.Name, backend)

	return nil
}
