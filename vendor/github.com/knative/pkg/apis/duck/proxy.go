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

	"k8s.io/apimachinery/pkg/watch"
)

// NewProxyWatcher is based on the same concept from Kubernetes apimachinery in 1.12 here:
//    https://github.com/kubernetes/apimachinery/blob/c6dd271be/pkg/watch/watch.go#L272
// Replace this copy once we've update our client libraries.

// proxyWatcher lets you wrap your channel in watch.Interface. Threadsafe.
type proxyWatcher struct {
	result chan watch.Event
	stopCh chan struct{}

	mutex   sync.Mutex
	stopped bool
}

var _ watch.Interface = (*proxyWatcher)(nil)

// NewProxyWatcher creates new proxyWatcher by wrapping a channel
func NewProxyWatcher(ch chan watch.Event) watch.Interface {
	return &proxyWatcher{
		result:  ch,
		stopCh:  make(chan struct{}),
		stopped: false,
	}
}

// Stop implements Interface
func (pw *proxyWatcher) Stop() {
	pw.mutex.Lock()
	defer pw.mutex.Unlock()
	if !pw.stopped {
		pw.stopped = true
		close(pw.stopCh)
	}
}

// Stopping returns true if Stop() has been called
func (pw *proxyWatcher) Stopping() bool {
	pw.mutex.Lock()
	defer pw.mutex.Unlock()
	return pw.stopped
}

// ResultChan implements watch.Interface
func (pw *proxyWatcher) ResultChan() <-chan watch.Event {
	return pw.result
}

// StopChan returns stop channel
func (pw *proxyWatcher) StopChan() <-chan struct{} {
	return pw.stopCh
}
