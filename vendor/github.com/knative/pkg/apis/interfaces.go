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

package apis

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// Defaultable defines an interface for setting the defaults for the
// uninitialized fields of this instance.
type Defaultable interface {
	SetDefaults()
}

// Validatable indicates that a particular type may have its fields validated.
type Validatable interface {
	// Validate checks the validity of this types fields.
	Validate() *FieldError
}

// Immutable indicates that a particular type has fields that should
// not change after creation.
type Immutable interface {
	// CheckImmutableFields checks that the current instance's immutable
	// fields haven't changed from the provided original.
	CheckImmutableFields(original Immutable) *FieldError
}

// Listable indicates that a particular type can be returned via the returned
// list type by the API server.
type Listable interface {
	runtime.Object

	GetListType() runtime.Object
}
