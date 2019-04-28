/*
Copyright 2018 The Knative Authors

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
	"github.com/knative/pkg/apis"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/knative/pkg/kmeta"
)

// Template is an interface for accessing the BuildTemplateSpec
// from various forms of template (namespace-/cluster-scoped).
type Template interface {
	TemplateSpec() BuildTemplateSpec
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BuildTemplate is a template that can used to easily create Builds.
type BuildTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BuildTemplateSpec `json:"spec"`
}

// Check that our resource implements several interfaces.
var _ kmeta.OwnerRefable = (*BuildTemplate)(nil)
var _ Template = (*BuildTemplate)(nil)
var _ BuildTemplateInterface = (*BuildTemplate)(nil)

// Check that BuildTemplate may be validated and defaulted.
var _ apis.Validatable = (*BuildTemplate)(nil)
var _ apis.Defaultable = (*BuildTemplate)(nil)

// BuildTemplateSpec is the spec for a BuildTemplate.
type BuildTemplateSpec struct {
	// TODO: Generation does not work correctly with CRD. They are scrubbed
	// by the APIserver (https://github.com/kubernetes/kubernetes/issues/58778)
	// So, we add Generation here. Once that gets fixed, remove this and use
	// ObjectMeta.Generation instead.
	// +optional
	Generation int64 `json:"generation,omitempty"`

	// Parameters defines the parameters that can be populated in a template.
	Parameters []ParameterSpec `json:"parameters,omitempty"`

	// Steps are the steps of the build; each step is run sequentially with the
	// source mounted into /workspace.
	Steps []corev1.Container `json:"steps"`

	// Volumes is a collection of volumes that are available to mount into the
	// steps of the build.
	Volumes []corev1.Volume `json:"volumes"`
}

// ParameterSpec defines the possible parameters that can be populated in a
// template.
type ParameterSpec struct {
	// Name is the unique name of this template parameter.
	Name string `json:"name"`

	// Description is a human-readable explanation of this template parameter.
	Description string `json:"description,omitempty"`

	// Default, if specified, defines the default value that should be applied if
	// the build does not specify the value for this parameter.
	Default *string `json:"default,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BuildTemplateList is a list of BuildTemplate resources.
type BuildTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []BuildTemplate `json:"items"`
}

// TemplateSpec returnes the Spec used by the template
func (bt *BuildTemplate) TemplateSpec() BuildTemplateSpec {
	return bt.Spec
}

// Copy performes a deep copy
func (bt *BuildTemplate) Copy() BuildTemplateInterface {
	return bt.DeepCopy()
}

// GetGroupVersionKind gives kind
func (bt *BuildTemplate) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("BuildTemplate")
}

// SetDefaults for build template
func (bt *BuildTemplate) SetDefaults() {}
