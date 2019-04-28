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

package duck

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

// Implementable is implemented by the Fooable duck type that consumers
// are expected to embed as a `.status.fooable` field.
type Implementable interface {
	// GetFullType returns an instance of a full resource wrapping
	// an instance of this Implementable that can populate its fields
	// to verify json roundtripping.
	GetFullType() Populatable
}

// Populatable is implemented by a skeleton resource wrapping an Implementable
// duck type.  It will generally have TypeMeta, ObjectMeta, and a Status field
// wrapping a Fooable field.
type Populatable interface {
	// Populate fills in all possible fields, so that we can verify that
	// they roundtrip properly through JSON.
	Populate()
}

// VerifyType verifies that a particular concrete resource properly implements
// the provided Implementable duck type.  It is expected that under the resource
// definition implementing a particular "Fooable" that one would write:
//
//   type ConcreteResource struct { ... }
//
//   // Check that ConcreteResource properly implement Fooable.
//   err := duck.VerifyType(&ConcreteResource{}, &something.Fooable{})
//
// This will return an error if the duck typing is not satisfied.
func VerifyType(instance interface{}, iface Implementable) error {
	// Create instances of the full resource for our input and ultimate result
	// that we will compare at the end.
	input, output := iface.GetFullType(), iface.GetFullType()

	// Populate our input resource with values we will roundtrip.
	input.Populate()

	// Serialize the input to JSON and deserialize that into the provided instance
	// of the type that we are checking.
	if before, err := json.Marshal(input); err != nil {
		return fmt.Errorf("error serializing duck type %T", input)
	} else if err := json.Unmarshal(before, instance); err != nil {
		return fmt.Errorf("error deserializing duck type %T into %T", input, instance)
	}

	// Serialize the instance we are checking to JSON and deserialize that into the
	// output resource.
	if after, err := json.Marshal(instance); err != nil {
		return fmt.Errorf("error serializing %T", instance)
	} else if err := json.Unmarshal(after, output); err != nil {
		return fmt.Errorf("error deserializing %T into dock type %T", instance, output)
	}

	// Now verify that we were able to roundtrip all of our fields through the type
	// we are checking.
	if diff := cmp.Diff(input, output); diff != "" {
		return fmt.Errorf("%T does not implement the duck type %T, the following fields were lost: %s",
			instance, iface, diff)
	}
	return nil
}
