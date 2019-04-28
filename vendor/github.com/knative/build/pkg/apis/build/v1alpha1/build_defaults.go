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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefaultTimeout is 10min
const DefaultTimeout = 10 * time.Minute

// SetDefaults for build
func (b *Build) SetDefaults() {
	if b == nil {
		return
	}
	if b.Spec.ServiceAccountName == "" {
		b.Spec.ServiceAccountName = "default"
	}
	if b.Spec.Timeout == nil {
		b.Spec.Timeout = &metav1.Duration{Duration: DefaultTimeout}
	}
	if b.Spec.Template != nil && b.Spec.Template.Kind == "" {
		b.Spec.Template.Kind = BuildTemplateKind
	}
}
