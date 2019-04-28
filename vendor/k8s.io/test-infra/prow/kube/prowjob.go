/*
Copyright 2018 The Kubernetes Authors.

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

package kube

import (
	"k8s.io/test-infra/prow/apis/prowjobs/v1"
)

// The following are aliases to aid in the refactoring while we move
// API definitions under prow/apis/

// ProwJobType specifies how the job is triggered.
type ProwJobType = v1.ProwJobType

// ProwJobState specifies whether the job is running
type ProwJobState = v1.ProwJobState

// ProwJobAgent specifies the controller (such as plank or jenkins-agent) that runs the job.
type ProwJobAgent = v1.ProwJobAgent

// Various job types.
const (
	// PresubmitJob means it runs on unmerged PRs.
	PresubmitJob = v1.PresubmitJob
	// PostsubmitJob means it runs on each new commit.
	PostsubmitJob = v1.PostsubmitJob
	// Periodic job means it runs on a time-basis, unrelated to git changes.
	PeriodicJob = v1.PeriodicJob
	// BatchJob tests multiple unmerged PRs at the same time.
	BatchJob = v1.BatchJob
)

// Various job states.
const (
	// TriggeredState means the job has been created but not yet scheduled.
	TriggeredState = v1.TriggeredState
	// PendingState means the job is scheduled but not yet running.
	PendingState = v1.PendingState
	// SuccessState means the job completed without error (exit 0)
	SuccessState = v1.SuccessState
	// FailureState means the job completed with errors (exit non-zero)
	FailureState = v1.FailureState
	// AbortedState means prow killed the job early (new commit pushed, perhaps).
	AbortedState = v1.AbortedState
	// ErrorState means the job could not schedule (bad config, perhaps).
	ErrorState = v1.ErrorState
)

const (
	// KubernetesAgent means prow will create a pod to run this job.
	KubernetesAgent = v1.KubernetesAgent
	// JenkinsAgent means prow will schedule the job on jenkins.
	JenkinsAgent = v1.JenkinsAgent
)

const (
	// CreatedByProw is added on pods created by prow. We cannot
	// really use owner references because pods may reside on a
	// different namespace from the namespace parent prowjobs
	// live and that would cause the k8s garbage collector to
	// identify those prow pods as orphans and delete them
	// instantly.
	// TODO: Namespace this label.
	CreatedByProw = "created-by-prow"
	// ProwJobTypeLabel is added in pods created by prow and
	// carries the job type (presubmit, postsubmit, periodic, batch)
	// that the pod is running.
	ProwJobTypeLabel = "prow.k8s.io/type"
	// ProwJobIDLabel is added in pods created by prow and
	// carries the ID of the ProwJob that the pod is fulfilling.
	// We also name pods after the ProwJob that spawned them but
	// this allows for multiple resources to be linked to one
	// ProwJob.
	ProwJobIDLabel = "prow.k8s.io/id"
	// ProwJobAnnotation is added in pods created by prow and
	// carries the name of the job that the pod is running. Since
	// job names can be arbitrarily long, this is added as
	// an annotation instead of a label.
	ProwJobAnnotation = "prow.k8s.io/job"
	// OrgLabel is added in resources created by prow and
	// carries the org associated with the job, eg kubernetes-sigs.
	OrgLabel = "prow.k8s.io/refs.org"
	// RepoLabel is added in resources created by prow and
	// carries the repo associated with the job, eg test-infra
	RepoLabel = "prow.k8s.io/refs.repo"
	// PullLabel is added in resources created by prow and
	// carries the PR number associated with the job, eg 321.
	PullLabel = "prow.k8s.io/refs.pull"
)

// ProwJob contains the spec as well as runtime metadata.
type ProwJob = v1.ProwJob

// ProwJobSpec configures the details of the prow job.
//
// Details include the podspec, code to clone, the cluster it runs
// any child jobs, concurrency limitations, etc.
type ProwJobSpec = v1.ProwJobSpec

// DecorationConfig specifies how to augment pods.
//
// This is primarily used to provide automatic integration with gubernator
// and testgrid.
type DecorationConfig = v1.DecorationConfig

// UtilityImages holds pull specs for the utility images
// to be used for a job
type UtilityImages = v1.UtilityImages

// PathStrategy specifies minutia about how to contruct the url.
// Usually consumed by gubernator/testgrid.
const (
	PathStrategyLegacy   = v1.PathStrategyLegacy
	PathStrategySingle   = v1.PathStrategySingle
	PathStrategyExplicit = v1.PathStrategyExplicit
)

// GCSConfiguration holds options for pushing logs and
// artifacts to GCS from a job.
type GCSConfiguration = v1.GCSConfiguration

// ProwJobStatus provides runtime metadata, such as when it finished, whether it is running, etc.
type ProwJobStatus = v1.ProwJobStatus

// Pull describes a pull request at a particular point in time.
type Pull = v1.Pull

// Refs describes how the repo was constructed.
type Refs = v1.Refs
