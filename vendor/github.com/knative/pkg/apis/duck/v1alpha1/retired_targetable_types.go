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

// Targetable is an earlier version of the Callable interface.
// Callable is a higher-level interface which implements Addressable
// but further promises that the destination may synchronously return
// response messages in reply to a message.
//
// Targetable implementations should instead implement Addressable and
// include an `eventing.knative.dev/returns=any` annotation.

// Targetable is retired; implement Addressable for now.
type Targetable struct {
	DomainInternal string `json:"domainInternal,omitempty"`
}

// Targetable is an Implementable "duck type".
var _ duck.Implementable = (*Targetable)(nil)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Target is a skeleton type wrapping Targetable in the manner we expect
// resource writers defining compatible resources to embed it.  We will
// typically use this type to deserialize Targetable ObjectReferences and
// access the Targetable data.  This is not a real resource.
type Target struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status TargetStatus `json:"status"`
}

// TargetStatus shows how we expect folks to embed Targetable in
// their Status field.
type TargetStatus struct {
	Targetable *Targetable `json:"targetable,omitempty"`
}

// In order for Targetable to be Implementable, Target must be Populatable.
var _ duck.Populatable = (*Target)(nil)

// Ensure Target satisfies apis.Listable
var _ apis.Listable = (*Target)(nil)

// GetFullType implements duck.Implementable
func (_ *Targetable) GetFullType() duck.Populatable {
	return &Target{}
}

// Populate implements duck.Populatable
func (t *Target) Populate() {
	t.Status = TargetStatus{
		&Targetable{
			// Populate ALL fields
			DomainInternal: "this is not empty",
		},
	}
}

// GetListType implements apis.Listable
func (r *Target) GetListType() runtime.Object {
	return &TargetList{}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TargetList is a list of Target resources
type TargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Target `json:"items"`
}
