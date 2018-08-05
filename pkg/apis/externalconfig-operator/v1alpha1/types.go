package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ExternalConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ExternalConfig `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ExternalConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ExternalConfigSpec   `json:"spec"`
	Status            ExternalConfigStatus `json:"status,omitempty"`
}

type ExternalConfigSpec struct {
	Backend string
	Key     string
}
type ExternalConfigStatus struct {
	Injected bool
}
