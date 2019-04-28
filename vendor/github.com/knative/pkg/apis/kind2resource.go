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
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// KindToResource converts a GroupVersionKind to a GroupVersionResource
// through the world's simplest (worst) pluralizer.
func KindToResource(gvk schema.GroupVersionKind) schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: pluralizeKind(gvk.Kind),
	}
}

// Takes a kind and pluralizes it. This is super terrible, but I am
// not aware of a generic way to do this.
// I am not alone in thinking this and I haven't found a better solution:
// This seems relevant:
// https://github.com/kubernetes/kubernetes/issues/18622
func pluralizeKind(kind string) string {
	ret := strings.ToLower(kind)
	if strings.HasSuffix(ret, "s") {
		return fmt.Sprintf("%ses", ret)
	}
	return fmt.Sprintf("%ss", ret)
}
