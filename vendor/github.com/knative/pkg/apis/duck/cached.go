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
	"sync"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

// CachedInformerFactory implements InformerFactory by delegating to another
// InformerFactory, but memoizing the results.
type CachedInformerFactory struct {
	Delegate InformerFactory

	m     sync.Mutex
	cache map[schema.GroupVersionResource]*result
}

// Check that CachedInformerFactory implements InformerFactory.
var _ InformerFactory = (*CachedInformerFactory)(nil)

// Get implements InformerFactory.
func (cif *CachedInformerFactory) Get(gvr schema.GroupVersionResource) (cache.SharedIndexInformer, cache.GenericLister, error) {
	cif.m.Lock()
	if cif.cache == nil {
		cif.cache = make(map[schema.GroupVersionResource]*result)
	}
	elt, ok := cif.cache[gvr]
	if !ok {
		elt = &result{}
		elt.init = func() {
			elt.inf, elt.lister, elt.err = cif.Delegate.Get(gvr)
		}
		cif.cache[gvr] = elt
	}
	// If this were done via "defer", then TestDifferentGVRs will fail.
	cif.m.Unlock()

	// The call to the delegate could be slow because it syncs informers, so do
	// this outside of the main lock.
	return elt.Get()
}

type result struct {
	sync.Once
	init func()

	inf    cache.SharedIndexInformer
	lister cache.GenericLister
	err    error
}

func (t *result) Get() (cache.SharedIndexInformer, cache.GenericLister, error) {
	t.Do(t.init)
	return t.inf, t.lister, t.err
}
