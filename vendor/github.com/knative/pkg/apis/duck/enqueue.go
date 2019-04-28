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

// EnqueueInformerFactory implements InformerFactory by delegating to another
// InformerFactory, but attaching a ResourceEventHandler to the informer.
type EnqueueInformerFactory struct {
	Delegate InformerFactory

	EventHandler cache.ResourceEventHandler
}

// Check that EnqueueInformerFactory implements InformerFactory.
var _ InformerFactory = (*EnqueueInformerFactory)(nil)

// Get implements InformerFactory.
func (cif *EnqueueInformerFactory) Get(gvr schema.GroupVersionResource) (cache.SharedIndexInformer, cache.GenericLister, error) {
	inf, lister, err := cif.Delegate.Get(gvr)
	if err != nil {
		return nil, nil, err
	}
	// If there is an informer, attach our event handler.
	inf.AddEventHandler(cif.EventHandler)
	return inf, lister, nil
}
