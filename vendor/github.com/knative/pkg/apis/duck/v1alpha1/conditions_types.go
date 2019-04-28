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
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/knative/pkg/apis"
	"github.com/knative/pkg/apis/duck"
)

// Conditions is the schema for the conditions portion of the payload
type Conditions []Condition

// ConditionType is a camel-cased condition type.
type ConditionType string

const (
	// ConditionReady specifies that the resource is ready.
	// For long-running resources.
	ConditionReady ConditionType = "Ready"
	// ConditionSucceeded specifies that the resource has finished.
	// For resource which run to completion.
	ConditionSucceeded ConditionType = "Succeeded"
)

// ConditionSeverity expresses the severity of a Condition Type failing.
type ConditionSeverity string

const (
	// ConditionSeverityError specifies that a failure of a condition type
	// should be viewed as an error.
	ConditionSeverityError ConditionSeverity = "Error"
	// ConditionSeverityWarning specifies that a failure of a condition type
	// should be viewed as a warning, but that things could still work.
	ConditionSeverityWarning ConditionSeverity = "Warning"
	// ConditionSeverityInfo specifies that a failure of a condition type
	// should be viewed as purely informational, and that things could still work.
	ConditionSeverityInfo ConditionSeverity = "Info"
)

// Conditions defines a readiness condition for a Knative resource.
// See: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#typical-status-properties
// +k8s:deepcopy-gen=true
type Condition struct {
	// Type of condition.
	// +required
	Type ConditionType `json:"type" description:"type of status condition"`

	// Status of the condition, one of True, False, Unknown.
	// +required
	Status corev1.ConditionStatus `json:"status" description:"status of the condition, one of True, False, Unknown"`

	// Severity with which to treat failures of this type of condition.
	// When this is not specified, it defaults to Error.
	// +optional
	Severity ConditionSeverity `json:"severity,omitempty" description:"how to interpret failures of this condition, one of Error, Warning, Info"`

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// We use VolatileTime in place of metav1.Time to exclude this from creating equality.Semantic
	// differences (all other things held constant).
	// +optional
	LastTransitionTime apis.VolatileTime `json:"lastTransitionTime,omitempty" description:"last time the condition transit from one status to another"`

	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" description:"one-word CamelCase reason for the condition's last transition"`

	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" description:"human-readable message indicating details about last transition"`
}

// IsTrue is true if the condition is True
func (c *Condition) IsTrue() bool {
	if c == nil {
		return false
	}
	return c.Status == corev1.ConditionTrue
}

// IsFalse is true if the condition is False
func (c *Condition) IsFalse() bool {
	if c == nil {
		return false
	}
	return c.Status == corev1.ConditionFalse
}

// IsUnknown is true if the condition is Unknown
func (c *Condition) IsUnknown() bool {
	if c == nil {
		return true
	}
	return c.Status == corev1.ConditionUnknown
}

// Conditions is an Implementable "duck type".
var _ duck.Implementable = (*Conditions)(nil)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KResource is a skeleton type wrapping Conditions in the manner we expect
// resource writers defining compatible resources to embed it.  We will
// typically use this type to deserialize Conditions ObjectReferences and
// access the Conditions data.  This is not a real resource.
type KResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status KResourceStatus `json:"status"`
}

// KResourceStatus shows how we expect folks to embed Conditions in
// their Status field.
type KResourceStatus struct {
	Conditions Conditions `json:"conditions,omitempty"`
}

func (krs *KResourceStatus) GetConditions() Conditions {
	return krs.Conditions
}

func (krs *KResourceStatus) SetConditions(conditions Conditions) {
	krs.Conditions = conditions
}

// Ensure KResourceStatus satisfies ConditionsAccessor
var _ ConditionsAccessor = (*KResourceStatus)(nil)

// In order for Conditions to be Implementable, KResource must be Populatable.
var _ duck.Populatable = (*KResource)(nil)

// Ensure KResource satisfies apis.Listable
var _ apis.Listable = (*KResource)(nil)

// GetFullType implements duck.Implementable
func (_ *Conditions) GetFullType() duck.Populatable {
	return &KResource{}
}

// Populate implements duck.Populatable
func (t *KResource) Populate() {
	t.Status.Conditions = Conditions{{
		// Populate ALL fields
		Type:               "Birthday",
		Status:             corev1.ConditionTrue,
		LastTransitionTime: apis.VolatileTime{Inner: metav1.NewTime(time.Date(1984, 02, 28, 18, 52, 00, 00, time.UTC))},
		Reason:             "Celebrate",
		Message:            "n3wScott, find your party hat :tada:",
	}}
}

// GetListType implements apis.Listable
func (r *KResource) GetListType() runtime.Object {
	return &KResourceList{}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KResourceList is a list of KResource resources
type KResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []KResource `json:"items"`
}
