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

package kmeta

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

// The methods in this file are used for managing subresources in cases where
// a controller instantiates different resources for each version of itself.
// There are two sets of methods available here:
// * `*VersionLabel*`: these methods act on `metadata.resourceVersion` and
//  create new labels for EVERY change to the resource (incl. `/status`).
// * `*GenerationLabel*`: these methods act on `metadata.generation` and
//  create new labels for changes to the resource's "spec" (typically, but
//  some K8s resources change `metadata.generation` for annotations as well
//  e.g. Deployment).
//
// For example, if an A might instantiate N B's at version 1 and M B's at
// version 2 then it can use MakeVersionLabels to decorate each subresource
// with the appropriate labels for the version at which it was instantiated.
//
// During reconciliation, MakeVersionLabelSelector can be used with the
// informer listers to access the appropriate subresources for the current
// version of the parent resource.
//
// Likewise during reconciliation, MakeOldVersionLabelSelector can be used
// with the API client's DeleteCollection method to clean up subresources
// for older versions of the resource.

// MakeVersionLabels constructs a set of labels to apply to subresources
// instantiated at this version of the parent resource, so that we can
// efficiently select them.
func MakeVersionLabels(om metav1.ObjectMetaAccessor) labels.Set {
	return map[string]string{
		"controller": string(om.GetObjectMeta().GetUID()),
		"version":    om.GetObjectMeta().GetResourceVersion(),
	}
}

// MakeVersionLabelSelector constructs a selector for subresources
// instantiated at this version of the parent resource.  This keys
// off of the labels populated by MakeVersionLabels.
func MakeVersionLabelSelector(om metav1.ObjectMetaAccessor) labels.Selector {
	return labels.SelectorFromSet(MakeVersionLabels(om))
}

// MakeOldVersionLabelSelector constructs a selector for subresources
// instantiated at an older version of the parent resource.  This keys
// off of the labels populated by MakeVersionLabels.
func MakeOldVersionLabelSelector(om metav1.ObjectMetaAccessor) labels.Selector {
	return labels.NewSelector().Add(
		mustNewRequirement("controller", selection.Equals, []string{string(om.GetObjectMeta().GetUID())}),
		mustNewRequirement("version", selection.NotEquals, []string{om.GetObjectMeta().GetResourceVersion()}),
	)
}

// MakeGenerationLabels constructs a set of labels to apply to subresources
// instantiated at this version of the parent resource, so that we can
// efficiently select them.
func MakeGenerationLabels(om metav1.ObjectMetaAccessor) labels.Set {
	return map[string]string{
		"controller": string(om.GetObjectMeta().GetUID()),
		"generation": genStr(om),
	}
}

// MakeGenerationLabelSelector constructs a selector for subresources
// instantiated at this version of the parent resource.  This keys
// off of the labels populated by MakeGenerationLabels.
func MakeGenerationLabelSelector(om metav1.ObjectMetaAccessor) labels.Selector {
	return labels.SelectorFromSet(MakeGenerationLabels(om))
}

// MakeOldGenerationLabelSelector constructs a selector for subresources
// instantiated at an older version of the parent resource.  This keys
// off of the labels populated by MakeGenerationLabels.
func MakeOldGenerationLabelSelector(om metav1.ObjectMetaAccessor) labels.Selector {
	return labels.NewSelector().Add(
		mustNewRequirement("controller", selection.Equals, []string{string(om.GetObjectMeta().GetUID())}),
		mustNewRequirement("generation", selection.NotEquals, []string{genStr(om)}),
	)
}

func genStr(om metav1.ObjectMetaAccessor) string {
	return fmt.Sprintf("%05d", om.GetObjectMeta().GetGeneration())
}

// mustNewRequirement panics if there are any errors constructing our selectors.
func mustNewRequirement(key string, op selection.Operator, vals []string) labels.Requirement {
	r, err := labels.NewRequirement(key, op, vals)
	if err != nil {
		panic(fmt.Sprintf("mustNewRequirement(%v, %v, %v) = %v", key, op, vals, err))
	}
	return *r
}
