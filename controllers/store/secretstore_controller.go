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
	"time"

	storev1alpha1 "github.com/containersolutions/externalsecret-operator/apis/store/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/containersolutions/externalsecret-operator/pkg/backend"
	config "github.com/containersolutions/externalsecret-operator/pkg/config"
)

const (
	defaulRetryPeriod = time.Second * 30
)

// SecretStoreReconciler reconciles a SecretStore object
type SecretStoreReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=store.externalsecret-operator.container-solutions.com,resources=secretstores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=store.externalsecret-operator.container-solutions.com,resources=secretstores/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete

func (r *SecretStoreReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("secretstore", req.NamespacedName)

	log.Info("Reconciling SecretStore")
	defer log.Info("Reconcile SecretStore Complete")

	// Fetch the SecretStore instance
	secretStore := &storev1alpha1.SecretStore{}
	err := r.Get(ctx, req.NamespacedName, secretStore)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("SecretStore not found")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get SecretStore")
		return ctrl.Result{}, err
	}

	contrl := secretStore.Spec.Controller

	storeConfig := secretStore.Spec.Store.Raw

	config, err := config.ConfigFromCtrl(storeConfig)
	if err != nil {
		return ctrl.Result{}, err
	}

	secretRef := config.Auth["secretRef"].(map[string]interface{})

	// Fetch credential Secret
	credentialsSecret := &corev1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Name: secretRef["name"].(string), Namespace: secretRef["namespace"].(string)}, credentialsSecret)
	if err != nil {
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get credentials Secret")
		return ctrl.Result{RequeueAfter: defaulRetryPeriod}, err
	}

	credentials := credentialsSecret.Data["operator-credentials.json"]

	err = backend.InitFromCtrl(contrl, config, credentials)
	if err != nil {
		log.Error(err, "Failed to intialize backend")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SecretStoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&storev1alpha1.SecretStore{}).
		Complete(r)
}
