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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ExternalSecretStoreRef is a reference to the external secret SecretStore
type ExternalSecretStoreRef struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Type=string
	Name string `json:"name"`
}

// ExternalSecretTarget ...
type ExternalSecretTarget struct {
	//  Name of the target Secret Resource
	//  defaults to .metadata.name of the ExternalSecret. immutable.
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	CreationPolicy string `json:"creationPolicy,omitempty"`
	// +kubebuilder:validation:Optional
	Template runtime.RawExtension `json:"template,omitempty"`
}

// ExternalSecretData contains Key/Name and Version of keys to be retrieved
type ExternalSecretData struct {
	// The Key/Name of the secret held in the ExternalBackend
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key"`
	// Version of the secret to be retrieved
	Version string `json:"version,omitempty"`
}

// ExternalSecretSpec defines the desired state of ExternalSecret
type ExternalSecretSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Secrets
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxItems=20
	// +kubebuilder:validation:MinItems=1
	Data []ExternalSecretData `json:"data"`
	// SecretStore reference
	// +kubebuilder:validation:Required
	StoreRef ExternalSecretStoreRef `json:"storeRef"`

	// +kubebuilder:validation:Optional
	// Secret Rotation Period;
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	RefreshInterval string `json:"refreshInterval,omitempty"`
	// +kubebuilder:validation:Optional
	Target ExternalSecretTarget `json:"target,omitempty"`
}

// ExternalSecretStatus defines the observed state of ExternalSecret
type ExternalSecretStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Defines where the ExternalSecret is in its lifecycle
	Phase string `json:"phase,omitempty"`
	// Conditions represent the latest available observations of an object's state
	Conditions []metav1.Condition `json:"conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ExternalSecret is the Schema for the externalsecrets API
type ExternalSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalSecretSpec   `json:"spec"`
	Status ExternalSecretStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExternalSecretList contains a list of ExternalSecret
type ExternalSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalSecret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ExternalSecret{}, &ExternalSecretList{})
}
