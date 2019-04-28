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

package errorutil

import (
	"fmt"
	"strings"
)

// Aggregate represents an object that contains multiple errors, but does not
// necessarily have singular semantic meaning.
type Aggregate interface {
	error
	Errors() []error
	Strings() []string
}

// NewAggregate converts a slice of errors into an Aggregate interface, which
// is itself an implementation of the error interface.  If the slice is empty,
// this returns nil.
// It will check if any of the element of input error list is nil, to avoid
// nil pointer panic when call Error().
func NewAggregate(errlist ...error) Aggregate {
	if len(errlist) == 0 {
		return nil
	}
	// In case of input error list contains nil
	var errs []error
	for _, e := range errlist {
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return aggregate(errs)
}

// This helper implements the error and Errors interfaces.  Keeping it private
// prevents people from making an aggregate of 0 errors, which is not
// an error, but does satisfy the error interface.
type aggregate []error

// Error is part of the error interface.
func (agg aggregate) Error() string {
	if len(agg) == 0 {
		// This should never happen, really.
		return ""
	}
	return fmt.Sprintf("[%s]", strings.Join(agg.Strings(), ", "))
}

// Strings flattens the aggregate (and any sub aggregates) into a
// slice of strings.
func (agg aggregate) Strings() []string {
	strs := make([]string, 0, len(agg))
	for _, e := range agg {
		if subAgg, ok := e.(aggregate); ok {
			strs = append(strs, subAgg.Strings()...)
		} else {
			strs = append(strs, e.Error())
		}
	}
	return strs
}

// Errors is part of the Aggregate interface.
func (agg aggregate) Errors() []error {
	return []error(agg)
}
