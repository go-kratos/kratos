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
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/knative/pkg/apis"
	"github.com/knative/pkg/apis/duck"
)

// Generation is the schema for the generational portion of the payload
type Generation int64

// Generation is an Implementable "duck type".
var _ duck.Implementable = (*Generation)(nil)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Generational is a skeleton type wrapping Generation in the manner we expect
// resource writers defining compatible resources to embed it.  We will
// typically use this type to deserialize Generation ObjectReferences and
// access the Generation data.  This is not a real resource.
type Generational struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GenerationalSpec `json:"spec"`
}

// GenerationalSpec shows how we expect folks to embed Generation in
// their Spec field.
type GenerationalSpec struct {
	Generation Generation `json:"generation,omitempty"`
}

// In order for Generation to be Implementable, Generational must be Populatable.
var _ duck.Populatable = (*Generational)(nil)

// Ensure Generational satisfies apis.Listable
var _ apis.Listable = (*Generational)(nil)

// GetFullType implements duck.Implementable
func (_ *Generation) GetFullType() duck.Populatable {
	return &Generational{}
}

// Populate implements duck.Populatable
func (t *Generational) Populate() {
	t.Spec.Generation = 1234
}

// GetListType implements apis.Listable
func (r *Generational) GetListType() runtime.Object {
	return &GenerationalList{}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GenerationalList is a list of Generational resources
type GenerationalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Generational `json:"items"`
}
