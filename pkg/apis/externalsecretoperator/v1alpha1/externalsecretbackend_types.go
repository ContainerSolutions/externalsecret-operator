package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ExternalSecretBackendSpec defines the desired state of ExternalSecretBackend
// +k8s:openapi-gen=true
type ExternalSecretBackendSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Type       string
	Parameters map[string]string
}

// ExternalSecretBackendStatus defines the observed state of ExternalSecretBackend
// +k8s:openapi-gen=true
type ExternalSecretBackendStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Type        string
	Initialized string
	Parameters  map[string]string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExternalSecretBackend is the Schema for the externalsecretbackends API
// +k8s:openapi-gen=true
type ExternalSecretBackend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalSecretBackendSpec   `json:"spec,omitempty"`
	Status ExternalSecretBackendStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExternalSecretBackendList contains a list of ExternalSecretBackend
type ExternalSecretBackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalSecretBackend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ExternalSecretBackend{}, &ExternalSecretBackendList{})
}
