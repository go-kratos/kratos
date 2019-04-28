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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/knative/pkg/apis"
	"github.com/knative/pkg/kmeta"
)

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterBuildTemplate is a template that can used to easily create Builds.
type ClusterBuildTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BuildTemplateSpec `json:"spec"`
}

// Check that our resource implements several interfaces.
var _ kmeta.OwnerRefable = (*ClusterBuildTemplate)(nil)
var _ Template = (*ClusterBuildTemplate)(nil)
var _ BuildTemplateInterface = (*ClusterBuildTemplate)(nil)

// Check that ClusterBuildTemplate may be validated and defaulted.
var _ apis.Validatable = (*ClusterBuildTemplate)(nil)
var _ apis.Defaultable = (*ClusterBuildTemplate)(nil)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterBuildTemplateList is a list of BuildTemplate resources.
type ClusterBuildTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ClusterBuildTemplate `json:"items"`
}

// TemplateSpec returnes the Spec used by the template
func (bt *ClusterBuildTemplate) TemplateSpec() BuildTemplateSpec {
	return bt.Spec
}

// Copy performes a deep copy
func (bt *ClusterBuildTemplate) Copy() BuildTemplateInterface {
	return bt.DeepCopy()
}

func (bt *ClusterBuildTemplate) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("ClusterBuildTemplate")
}

// SetDefaults
func (b *ClusterBuildTemplate) SetDefaults() {}
