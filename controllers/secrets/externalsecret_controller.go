/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	secretsv1alpha1 "github.com/containersolutions/externalsecret-operator/apis/secrets/v1alpha1"
	storev1alpha1 "github.com/containersolutions/externalsecret-operator/apis/store/v1alpha1"

	"github.com/containersolutions/externalsecret-operator/pkg/backend"
)

const (
	// seconds
	defaulRetryPeriod = 30
)

// ExternalSecretReconciler reconciles a ExternalSecret object
type ExternalSecretReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=secrets.externalsecret-operator.container-solutions.com,resources=externalsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=secrets.externalsecret-operator.container-solutions.com,resources=externalsecrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=store.externalsecret-operator.container-solutions.com,resources=secretstores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *ExternalSecretReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("externalsecret", req.NamespacedName)

	log.Info("Reconciling ExternalSecret")
	defer log.Info("Reconcile ExternalSecret Complete")

	// Fetch the ExternalSecret instance
	externalSecret := &secretsv1alpha1.ExternalSecret{}
	err := r.Get(ctx, req.NamespacedName, externalSecret)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("External Secret not found.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ExternalSecret")
		return ctrl.Result{}, err
	}

	// Fetch referenced store
	secretStore := &storev1alpha1.SecretStore{}
	err = r.Get(ctx, types.NamespacedName{Name: externalSecret.Spec.StoreRef.Name, Namespace: externalSecret.Spec.StoreRef.Namespace}, secretStore)
	if err != nil {
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get SecretStore")
		return ctrl.Result{RequeueAfter: time.Second * defaulRetryPeriod}, err
	}

	// Check if this Secret already exists
	foundSecret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: externalSecret.Name, Namespace: externalSecret.Namespace}, foundSecret)
	if err != nil {
		if errors.IsNotFound(err) {
			// Define a new Secret object
			secret, err := r.newSecretForCR(externalSecret, secretStore)
			if err != nil {
				log.Error(err, "Failed to create Secret")
				return ctrl.Result{RequeueAfter: time.Second * defaulRetryPeriod}, err
			}

			log.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
			err = r.Create(ctx, secret)
			if err != nil {
				log.Error(err, "Failed to create Secret", "secret", secret)
				return ctrl.Result{}, err
			}

			// Secret created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Secret")
		return ctrl.Result{}, err
	}

	// update Secret if it already exists
	secretMap, err := r.backendGet(externalSecret, secretStore)
	if err != nil {
		log.Error(err, "backendGet")
		return ctrl.Result{}, err
	}

	updateLabels := makeLabels(secretStore.Spec.Controller, externalSecret.Spec.StoreRef.Name)

	foundSecret.ObjectMeta.Labels = updateLabels
	foundSecret.Data = secretMap
	err = r.Update(ctx, foundSecret)
	if err != nil {
		log.Error(err, "Failed to update secret")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ExternalSecretReconciler) newSecretForCR(s *secretsv1alpha1.ExternalSecret, st *storev1alpha1.SecretStore) (*corev1.Secret, error) {
	secretMap, err := r.backendGet(s, st)
	if err != nil {
		log.Error(err, "backendGet")
		return nil, err
	}

	secretLabels := makeLabels(st.Spec.Controller, s.Spec.StoreRef.Name)

	secretObject := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
			Labels:    secretLabels,
		},
		Data: secretMap,
	}

	// Allows deleted objects to be garbage collected.
	err = ctrl.SetControllerReference(s, secretObject, r.Scheme)
	if err != nil {
		log.Error(err, "Error setting owner references", secretObject)
		return nil, err
	}

	return secretObject, nil
}

func (r *ExternalSecretReconciler) backendGet(s *secretsv1alpha1.ExternalSecret, st *storev1alpha1.SecretStore) (map[string][]byte, error) {
	secrets := s.Spec.Secrets
	secretMap := make(map[string][]byte)

	stCtrl := st.Spec.Controller
	backend, ok := backend.Instances[stCtrl]
	if !ok {
		log.Error("Cannot find controller:", stCtrl)
		return secretMap, fmt.Errorf("Cannot find backend: %v", stCtrl)
	}

	for _, secret := range secrets {
		retrievedValue, err := backend.Get(secret.Key, secret.Version)
		if err != nil {
			log.Error(err, "could not create secret due to error from backend")
			return secretMap, fmt.Errorf("could not create secret due to error from backend: %v", err)
		}

		secretMap[secret.Key] = []byte(retrievedValue)
	}

	return secretMap, nil
}

func makeLabels(contrl string, storeRef string) map[string]string {
	return map[string]string{
		"secret-controller": contrl,
		"secret-storeRef":   storeRef,
	}
}

func (r *ExternalSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretsv1alpha1.ExternalSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
