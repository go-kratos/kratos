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

// LegacyTargetable left around until we migrate to Addressable in the
// dependent resources. Addressable has more structure in the way it
// defines the fields. LegacyTargetable only assumed a single string
// in the Status field and we're moving towards defining proper structs
// under Status rather than strings.
// This is to support existing resources until they migrate.
//
// Do not use this for anything new, use Addressable

// LegacyTargetable is the old schema for the addressable portion
// of the payload
//
// For new resources use Addressable.
type LegacyTargetable struct {
	DomainInternal string `json:"domainInternal,omitempty"`
}

// LegacyTargetable is an Implementable "duck type".
var _ duck.Implementable = (*LegacyTargetable)(nil)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LegacyTarget is a skeleton type wrapping LegacyTargetable in the manner we
// want to support unless they get migrated into supporting Legacy.
// We will typically use this type to deserialize LegacyTargetable
// ObjectReferences and access the LegacyTargetable data.  This is not a
// real resource.
// ** Do not use this for any new resources **
type LegacyTarget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status LegacyTargetable `json:"status"`
}

// In order for LegacyTargetable to be Implementable, LegacyTarget must be Populatable.
var _ duck.Populatable = (*LegacyTarget)(nil)

// Ensure LegacyTarget satisfies apis.Listable
var _ apis.Listable = (*LegacyTarget)(nil)

// GetFullType implements duck.Implementable
func (_ *LegacyTargetable) GetFullType() duck.Populatable {
	return &LegacyTarget{}
}

// Populate implements duck.Populatable
func (t *LegacyTarget) Populate() {
	t.Status = LegacyTargetable{
		// Populate ALL fields
		DomainInternal: "this is not empty",
	}
}

// GetListType implements apis.Listable
func (r *LegacyTarget) GetListType() runtime.Object {
	return &LegacyTargetList{}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LegacyTargetList is a list of LegacyTarget resources
type LegacyTargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []LegacyTarget `json:"items"`
}
