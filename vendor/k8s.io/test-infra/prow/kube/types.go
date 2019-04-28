/*
Copyright 2016 The Kubernetes Authors.

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
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: Drop all of these, please!

// ObjectMeta is a kubernetes v1 ObjectMeta
type ObjectMeta = metav1.ObjectMeta

// Pod is a kubernetes v1 Pod
type Pod = v1.Pod

// PodTemplateSpec is a kubernetes v1 PodTemplateSpec
type PodTemplateSpec = v1.PodTemplateSpec

// PodSpec is a kubernetes v1 PodSpec
type PodSpec = v1.PodSpec

// PodStatus is a kubernetes v1 PodStatus
type PodStatus = v1.PodStatus

// Phase constants
const (
	PodPending   = v1.PodPending
	PodRunning   = v1.PodRunning
	PodSucceeded = v1.PodSucceeded
	PodFailed    = v1.PodFailed
	PodUnknown   = v1.PodUnknown
)

// PodStatus constants
const (
	Evicted = "Evicted"
)

// Container is a kubernetes v1 Container
type Container = v1.Container

// Port is a kubernetes v1 ContainerPort
type Port = v1.ContainerPort

// EnvVar is a kubernetes v1 EnvVar
type EnvVar = v1.EnvVar

// Volume is a kubernetes v1 Volume
type Volume = v1.Volume

// VolumeMount is a kubernetes v1 VolumeMount
type VolumeMount = v1.VolumeMount

// VolumeSource is a kubernetes v1 VolumeSource
type VolumeSource = v1.VolumeSource

// EmptyDirVolumeSource is a kubernetes v1 EmptyDirVolumeSource
type EmptyDirVolumeSource = v1.EmptyDirVolumeSource

// SecretSource is a kubernetes v1 SecretVolumeSource
type SecretSource = v1.SecretVolumeSource

// ConfigMapSource is a kubernetes v1 ConfigMapVolumeSource
type ConfigMapSource = v1.ConfigMapVolumeSource

// ConfigMap is a kubernetes v1 ConfigMap
type ConfigMap = v1.ConfigMap

// Secret is a kubernetes v1 secret
type Secret = v1.Secret
