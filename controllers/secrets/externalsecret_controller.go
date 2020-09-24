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
	"github.com/containersolutions/externalsecret-operator/pkg/backend"

	// trigger secrets backend registration
	_ "github.com/containersolutions/externalsecret-operator/pkg/backend"
)

// ExternalSecretReconciler reconciles a ExternalSecret object
type ExternalSecretReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=secrets.externalsecret-operator.container-solutions.com,resources=externalsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=secrets.externalsecret-operator.container-solutions.com,resources=externalsecrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *ExternalSecretReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("externalsecret", req.NamespacedName)

	log.Info("Reconciling ExternalSecret")

	// Fetch the ExternalSecret instance
	externalSecret := &secretsv1alpha1.ExternalSecret{}
	err := r.Get(ctx, req.NamespacedName, externalSecret)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ExternalSecret")
		return ctrl.Result{}, err
	}

	// Check if this Secret already exists
	found := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: externalSecret.Name, Namespace: externalSecret.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Secret object
		secret, err := r.newSecretForCR(externalSecret)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
		log.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		err = r.Create(ctx, secret)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Secret created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Secret already exists - don't requeue
	log.Info("Skip Reconcile: Secret already exists", "Secret.Namespace", found.Namespace, "Secret.Name", found.Name)

	return ctrl.Result{}, nil
}

func (r *ExternalSecretReconciler) newSecretForCR(s *secretsv1alpha1.ExternalSecret) (*corev1.Secret, error) {
	backend, ok := backend.Instances[s.Spec.Backend]
	if !ok {
		return nil, fmt.Errorf("Cannot find backend: %v", s.Spec.Backend)
	}

	value, err := backend.Get(s.Spec.Key)
	if err != nil {
		log.Error(err, "could not create secret due to error from backend")
	}

	secret := map[string][]byte{s.Spec.Key: []byte(value)}

	secretObject := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Data: secret,
	}

	ctrl.SetControllerReference(s, secretObject, r.Scheme)

	return secretObject, err
}

func (r *ExternalSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretsv1alpha1.ExternalSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
