package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ModelDeploymentSpec defines the desired state of ModelDeployment
type ModelDeploymentSpec struct {
	ModelName  string `json:"modelName,omitempty"`
	Image      string `json:"image,omitempty"`
	Replicas   *int32 `json:"replicas,omitempty"`
	WeightsPVC string `json:"weightsPVC,omitempty"`
}

// ModelDeploymentStatus defines the observed state of ModelDeployment
type ModelDeploymentStatus struct {
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ModelDeployment is the Schema for the modeldeployments API
type ModelDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelDeploymentSpec   `json:"spec,omitempty"`
	Status ModelDeploymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ModelDeploymentList contains a list of ModelDeployment
type ModelDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelDeployment{}, &ModelDeploymentList{})
}

// AddToScheme registers this API group with the given scheme
func AddToScheme(s *runtime.Scheme) error {
	return SchemeBuilder.AddToScheme(s)
}
