package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CurrentTemperatureSpec defines the desired state of CurrentTemperature.
type CurrentTemperatureSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// City is the name of the city for which the current temperature is to be fetched.
	City string `json:"city"`
	// FlowUnit is the unit of measurement for the flow rate.
	// +kubebuilder:validation:Enum=m3/s;Beer/s
	// +kubebuilder:default=m3/s
	FlowUnit string `json:"flowUnit,omitempty"`
	// +kubebuilder:validation:Required
	// UpdateInterval defines how often the current temperature should be updated.
	UpdateInterval metav1.Duration `json:"updateInterval"`
}

// CurrentTemperatureStatus defines the observed state of CurrentTemperature.
type CurrentTemperatureStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Location    string      `json:"location"`
	Temperature string      `json:"temperature"`
	Text        string      `json:"text"`
	Flow        string      `json:"flow"`
	Updated     metav1.Time `json:"time"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=cta
// +kubebuilder:printcolumn:name="LOCATION",type=string,JSONPath=`.status.location`
// +kubebuilder:printcolumn:name="TEMPERATURE",type=string,JSONPath=`.status.temperature`
// +kubebuilder:printcolumn:name="Flow",type=string,JSONPath=`.status.flow`
// +kubebuilder:printcolumn:name="TEXT",type=string,JSONPath=`.status.text`
// +kubebuilder:printcolumn:name="UPDATED",type=string,JSONPath=`.status.time`

// CurrentTemperature is the Schema for the currenttemperatures API.
type CurrentTemperature struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CurrentTemperatureSpec   `json:"spec,omitempty"`
	Status CurrentTemperatureStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CurrentTemperatureList contains a list of CurrentTemperature.
type CurrentTemperatureList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CurrentTemperature `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CurrentTemperature{}, &CurrentTemperatureList{})
}
