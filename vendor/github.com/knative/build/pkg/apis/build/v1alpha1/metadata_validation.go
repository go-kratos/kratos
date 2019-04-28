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
	"strings"

	"github.com/knative/pkg/apis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	maxLength = 63
)

func validateObjectMetadata(meta metav1.Object) *apis.FieldError {
	name := meta.GetName()

	if strings.Contains(name, ".") {
		return &apis.FieldError{
			Message: "Invalid resource name: special character . must not be present",
			Paths:   []string{"name"},
		}
	}

	if len(name) > maxLength {
		return &apis.FieldError{
			Message: "Invalid resource name: length must be no more than 63 characters",
			Paths:   []string{"name"},
		}
	}
	return nil
}
