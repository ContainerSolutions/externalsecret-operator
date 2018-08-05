package stub

import (
	"context"

	"github.com/ContainerSolutions/externalconfig-operator/pkg/apis/externalconfig-operator/v1alpha1"
	"github.com/ContainerSolutions/externalconfig-operator/pkg/secrets"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	SecretsBackend secrets.SecretsBackend
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.ExternalConfig:
		//FIXME: Status doesn't work
		if o.Status.Injected {
			return nil
		}
		secret, err := h.makeSecret(ctx, o)
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to retrieve secret from backend: %v", err)
			return err
		}
		err = sdk.Create(secret)
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create new secret: %v", err)
			return err
		}
		o.Status.Injected = true
		logrus.Info("Created secret %v", secret)
	}
	return nil
}

func (h *Handler) makeSecret(ctx context.Context, cr *v1alpha1.ExternalConfig) (*corev1.Secret, error) {
	var backendKey secrets.ContextKey = "backend"
	var backend secrets.SecretsBackend
	key := cr.Spec.Key
	backend = ctx.Value(backendKey).(secrets.SecretsBackend)
	value, err := backend.Get(cr.Spec.Key)
	if err != nil {
		return nil, err
	}
	secret := map[string][]byte{key: []byte(value)}

	secretObject := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "ExternalConfig",
				}),
			},
		},
		Data: secret,
	}

	return secretObject, err
}
