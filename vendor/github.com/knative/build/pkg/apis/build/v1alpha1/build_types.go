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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/knative/pkg/apis"
	duckv1alpha1 "github.com/knative/pkg/apis/duck/v1alpha1"
	"github.com/knative/pkg/kmeta"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Build represents a build of a container image. A Build is made up of a
// source, and a set of steps. Steps can mount volumes to share data between
// themselves. A build may be created by instantiating a BuildTemplate.
type Build struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BuildSpec   `json:"spec"`
	Status BuildStatus `json:"status"`
}

// Check that our resource implements several interfaces.
var _ kmeta.OwnerRefable = (*Build)(nil)

// Check that Build may be validated and defaulted.
var _ apis.Validatable = (*Build)(nil)
var _ apis.Defaultable = (*Build)(nil)

// BuildSpec is the spec for a Build resource.
type BuildSpec struct {
	// TODO: Generation does not work correctly with CRD. They are scrubbed
	// by the APIserver (https://github.com/kubernetes/kubernetes/issues/58778)
	// So, we add Generation here. Once that gets fixed, remove this and use
	// ObjectMeta.Generation instead.
	// +optional
	Generation int64 `json:"generation,omitempty"`

	// Source specifies the input to the build.
	// +optional
	Source *SourceSpec `json:"source,omitempty"`

	// Sources specifies the inputs to the build.
	// +optional
	Sources []SourceSpec `json:"sources,omitempty"`

	// Steps are the steps of the build; each step is run sequentially with the
	// source mounted into /workspace.
	// +optional
	Steps []corev1.Container `json:"steps,omitempty"`

	// Volumes is a collection of volumes that are available to mount into the
	// steps of the build.
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// The name of the service account as which to run this build.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// Template, if specified, references a BuildTemplate resource to use to
	// populate fields in the build, and optional Arguments to pass to the
	// template. The default Kind of template is BuildTemplate
	// +optional
	Template *TemplateInstantiationSpec `json:"template,omitempty"`

	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Time after which the build times out. Defaults to 10 minutes.
	// Specified build timeout should be less than 24h.
	// Refer Go's ParseDuration documentation for expected format: https://golang.org/pkg/time/#ParseDuration
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// If specified, the pod's scheduling constraints
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
}

// TemplateKind defines the type of BuildTemplate used by the build.
type TemplateKind string

const (
	// BuildTemplateKind indicates that the template type has a namepace scope.
	BuildTemplateKind TemplateKind = "BuildTemplate"
	// ClusterBuildTemplateKind indicates that template type has a cluster scope.
	ClusterBuildTemplateKind TemplateKind = "ClusterBuildTemplate"
)

// TemplateInstantiationSpec specifies how a BuildTemplate is instantiated into
// a Build.
type TemplateInstantiationSpec struct {
	// Name references the BuildTemplate resource to use.
	// The template is assumed to exist in the Build's namespace.
	Name string `json:"name"`

	// The Kind of the template to be used, possible values are BuildTemplate
	// or ClusterBuildTemplate. If nothing is specified, the default if is BuildTemplate
	// +optional
	Kind TemplateKind `json:"kind,omitempty"`

	// Arguments, if specified, lists values that should be applied to the
	// parameters specified by the template.
	// +optional
	Arguments []ArgumentSpec `json:"arguments,omitempty"`

	// Env, if specified will provide variables to all build template steps.
	// This will override any of the template's steps environment variables.
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
}

// ArgumentSpec defines the actual values to use to populate a template's
// parameters.
type ArgumentSpec struct {
	// Name is the name of the argument.
	Name string `json:"name"`
	// Value is the value of the argument.
	Value string `json:"value"`
	// TODO(jasonhall): ValueFrom?
}

// SourceSpec defines the input to the Build
type SourceSpec struct {
	// Git represents source in a Git repository.
	// +optional
	Git *GitSourceSpec `json:"git,omitempty"`

	// GCS represents source in Google Cloud Storage.
	// +optional
	GCS *GCSSourceSpec `json:"gcs,omitempty"`

	// Custom indicates that source should be retrieved using a custom
	// process defined in a container invocation.
	// +optional
	Custom *corev1.Container `json:"custom,omitempty"`

	// SubPath specifies a path within the fetched source which should be
	// built. This option makes parent directories *inaccessible* to the
	// build steps. (The specific source type may, in fact, not even fetch
	// files not in the SubPath.)
	// +optional
	SubPath string `json:"subPath,omitempty"`

	// Name is the name of source. This field is used to uniquely identify the
	// source init containers
	// Restrictions on the allowed charatcers
	// Must be a basename (no /)
	// Must be a valid DNS name (only alphanumeric characters, no _)
	// https://tools.ietf.org/html/rfc1123#section-2
	// +optional
	Name string `json:"name,omitempty"`

	// TargetPath is the path in workspace directory where the source will be copied.
	// TargetPath is optional and if its not set source will be copied under workspace.
	// TargetPath should not be set for custom source.
	TargetPath string `json:"targetPath,omitempty"`
}

// GitSourceSpec describes a Git repo source input to the Build.
type GitSourceSpec struct {
	// URL of the Git repository to clone from.
	Url string `json:"url"`

	// Git revision (branch, tag, commit SHA or ref) to clone.  See
	// https://git-scm.com/docs/gitrevisions#_specifying_revisions for more
	// information.
	Revision string `json:"revision"`
}

// GCSSourceSpec describes source input to the Build in the form of an archive,
// or a source manifest describing files to fetch.
type GCSSourceSpec struct {
	// Type declares the style of source to fetch.
	Type GCSSourceType `json:"type,omitempty"`

	// Location specifies the location of the source archive or manifest file.
	Location string `json:"location,omitempty"`
}

// GCSSourceType defines a type of GCS source fetch.
type GCSSourceType string

const (
	// GCSArchive indicates that source should be fetched from a typical archive file.
	GCSArchive GCSSourceType = "Archive"

	// GCSManifest indicates that source should be fetched using a
	// manifest-based protocol which enables incremental source upload.
	GCSManifest GCSSourceType = "Manifest"
)

// BuildProvider defines a build execution implementation.
type BuildProvider string

const (
	// GoogleBuildProvider indicates that this build was performed with Google Cloud Build.
	GoogleBuildProvider BuildProvider = "Google"
	// ClusterBuildProvider indicates that this build was performed on-cluster.
	ClusterBuildProvider BuildProvider = "Cluster"
)

// BuildStatus is the status for a Build resource
type BuildStatus struct {
	// +optional
	Builder BuildProvider `json:"builder,omitempty"`

	// Cluster provides additional information if the builder is Cluster.
	// +optional
	Cluster *ClusterSpec `json:"cluster,omitempty"`

	// Google provides additional information if the builder is Google.
	// +optional
	Google *GoogleSpec `json:"google,omitempty"`

	// StartTime is the time the build is actually started.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is the time the build completed.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// StepStates describes the state of each build step container.
	// +optional
	StepStates []corev1.ContainerState `json:"stepStates,omitempty"`

	// StepsCompleted lists the name of build steps completed.
	// +optional
	StepsCompleted []string `json:"stepsCompleted",omitempty`

	// Conditions describes the set of conditions of this build.
	// +optional
	Conditions duckv1alpha1.Conditions `json:"conditions,omitempty"`
}

// Check that BuildStatus may have its conditions managed.
var _ duckv1alpha1.ConditionsAccessor = (*BuildStatus)(nil)

// ClusterSpec provides information about the on-cluster build, if applicable.
type ClusterSpec struct {
	// Namespace is the namespace in which the pod is running.
	Namespace string `json:"namespace"`
	// PodName is the name of the pod responsible for executing this build's steps.
	PodName string `json:"podName"`
}

// GoogleSpec provides information about the GCB build, if applicable.
type GoogleSpec struct {
	// Operation is the unique name of the GCB API Operation for the build.
	Operation string `json:"operation"`
}

// BuildSucceeded is set when the build is running, and becomes True when the
// build finishes successfully.
//
// If the build is ongoing, its status will be Unknown. If it fails, its status
// will be False.
const BuildSucceeded = duckv1alpha1.ConditionSucceeded

var buildCondSet = duckv1alpha1.NewBatchConditionSet()

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BuildList is a list of Build resources
type BuildList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items is the list of Build items in this list.
	Items []Build `json:"items"`
}

// GetCondition returns the Condition matching the given type.
func (bs *BuildStatus) GetCondition(t duckv1alpha1.ConditionType) *duckv1alpha1.Condition {
	return buildCondSet.Manage(bs).GetCondition(t)
}

// SetCondition sets the condition, unsetting previous conditions with the same
// type as necessary.
func (bs *BuildStatus) SetCondition(newCond *duckv1alpha1.Condition) {
	if newCond != nil {
		buildCondSet.Manage(bs).SetCondition(*newCond)
	}
}

// GetConditions returns the Conditions array. This enables generic handling of
// conditions by implementing the duckv1alpha1.Conditions interface.
func (bs *BuildStatus) GetConditions() duckv1alpha1.Conditions {
	return bs.Conditions
}

// SetConditions sets the Conditions array. This enables generic handling of
// conditions by implementing the duckv1alpha1.Conditions interface.
func (bs *BuildStatus) SetConditions(conditions duckv1alpha1.Conditions) {
	bs.Conditions = conditions
}

func (b *Build) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("Build")
}
