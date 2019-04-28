/*
Copyright 2014 Workiva, LLC

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

package queue

import "errors"

var (
	// ErrDisposed is returned when an operation is performed on a disposed
	// queue.
	ErrDisposed = errors.New(`queue: disposed`)

	// ErrTimeout is returned when an applicable queue operation times out.
	ErrTimeout = errors.New(`queue: poll timed out`)

	// ErrEmptyQueue is returned when an non-applicable queue operation was called
	// due to the queue's empty item state
	ErrEmptyQueue = errors.New(`queue: empty queue`)
)
