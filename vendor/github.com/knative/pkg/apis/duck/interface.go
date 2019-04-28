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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

// InformerFactory is used to create Informer/Lister pairs for a schema.GroupVersionResource
type InformerFactory interface {
	// Get returns a synced Informer/Lister pair for the provided schema.GroupVersionResource.
	Get(schema.GroupVersionResource) (cache.SharedIndexInformer, cache.GenericLister, error)
}
