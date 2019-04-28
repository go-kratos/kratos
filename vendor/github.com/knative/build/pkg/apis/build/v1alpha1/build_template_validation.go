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
	"github.com/knative/pkg/apis"
	corev1 "k8s.io/api/core/v1"
)

// Validate build template
func (b *BuildTemplate) Validate() *apis.FieldError {
	return validateObjectMetadata(b.GetObjectMeta()).ViaField("metadata").Also(b.Spec.Validate().ViaField("spec"))
}

// Validate Build Template
func (b *BuildTemplateSpec) Validate() *apis.FieldError {
	if err := validateSteps(b.Steps); err != nil {
		return err
	}
	if err := ValidateVolumes(b.Volumes); err != nil {
		return err
	}
	if err := validateParameters(b.Parameters); err != nil {
		return err
	}
	return nil
}

//ValidateVolumes validates collection of volumes that are available to mount into the
// steps of the build ot build template.
func ValidateVolumes(volumes []corev1.Volume) *apis.FieldError {
	// Build must not duplicate volume names.
	vols := map[string]struct{}{}
	for _, v := range volumes {
		if _, ok := vols[v.Name]; ok {
			return apis.ErrMultipleOneOf("volumeName")
		}
		vols[v.Name] = struct{}{}
	}
	return nil
}

func validateSteps(steps []corev1.Container) *apis.FieldError {
	// Build must not duplicate step names.
	names := map[string]struct{}{}
	for _, s := range steps {
		if s.Image == "" {
			return apis.ErrMissingField("Image")
		}

		if s.Name == "" {
			continue
		}
		if _, ok := names[s.Name]; ok {
			return apis.ErrMultipleOneOf("stepName")
		}
		names[s.Name] = struct{}{}
	}
	return nil
}

func validateParameters(params []ParameterSpec) *apis.FieldError {
	// Template must not duplicate parameter names.
	seen := map[string]struct{}{}
	for _, p := range params {
		if _, ok := seen[p.Name]; ok {
			return apis.ErrInvalidKeyName("ParamName", "b.spec.params")
		}
		seen[p.Name] = struct{}{}
	}
	return nil
}
