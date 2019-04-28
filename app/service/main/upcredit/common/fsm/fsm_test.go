// Copyright (c) 2013 - Max Persson <max@looplab.se>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fsm

import (
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"
)

type fakeTransitionerObj struct {
}

func (t fakeTransitionerObj) transition(f *FSM) error {
	return &InternalError{}
}

func TestSameState(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "start"},
		},
		Callbacks{},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}

func TestSetState(t *testing.T) {
	fsm := NewFSM(
		"walking",
		Events{
			{Name: "walk", Src: []string{"start"}, Dst: "walking"},
		},
		Callbacks{},
	)
	fsm.SetState("start")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'walking'")
	}
	err := fsm.Event("walk")
	if err != nil {
		t.Error("transition is expected no error")
	}
}

func TestBadTransition(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "running"},
		},
		Callbacks{},
	)
	fsm.transitionerObj = new(fakeTransitionerObj)
	err := fsm.Event("run")
	if err == nil {
		t.Error("bad transition should give an error")
	}
}

func TestInappropriateEvent(t *testing.T) {
	fsm := NewFSM(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Callbacks{},
	)
	err := fsm.Event("close")
	if e, ok := err.(InvalidEventError); !ok && e.Event != "close" && e.State != "closed" {
		t.Error("expected 'InvalidEventError' with correct state and event")
	}
}

func TestInvalidEvent(t *testing.T) {
	fsm := NewFSM(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Callbacks{},
	)
	err := fsm.Event("lock")
	if e, ok := err.(UnknownEventError); !ok && e.Event != "close" {
		t.Error("expected 'UnknownEventError' with correct event")
	}
}

func TestMultipleSources(t *testing.T) {
	fsm := NewFSM(
		"one",
		Events{
			{Name: "first", Src: []string{"one"}, Dst: "two"},
			{Name: "second", Src: []string{"two"}, Dst: "three"},
			{Name: "reset", Src: []string{"one", "two", "three"}, Dst: "one"},
		},
		Callbacks{},
	)

	fsm.Event("first")
	if fsm.Current() != "two" {
		t.Error("expected state to be 'two'")
	}
	fsm.Event("reset")
	if fsm.Current() != "one" {
		t.Error("expected state to be 'one'")
	}
	fsm.Event("first")
	fsm.Event("second")
	if fsm.Current() != "three" {
		t.Error("expected state to be 'three'")
	}
	fsm.Event("reset")
	if fsm.Current() != "one" {
		t.Error("expected state to be 'one'")
	}
}

func TestMultipleEvents(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "first", Src: []string{"start"}, Dst: "one"},
			{Name: "second", Src: []string{"start"}, Dst: "two"},
			{Name: "reset", Src: []string{"one"}, Dst: "reset_one"},
			{Name: "reset", Src: []string{"two"}, Dst: "reset_two"},
			{Name: "reset", Src: []string{"reset_one", "reset_two"}, Dst: "start"},
		},
		Callbacks{},
	)

	fsm.Event("first")
	fsm.Event("reset")
	if fsm.Current() != "reset_one" {
		t.Error("expected state to be 'reset_one'")
	}
	fsm.Event("reset")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}

	fsm.Event("second")
	fsm.Event("reset")
	if fsm.Current() != "reset_two" {
		t.Error("expected state to be 'reset_two'")
	}
	fsm.Event("reset")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}

func TestGenericCallbacks(t *testing.T) {
	beforeEvent := false
	leaveState := false
	enterState := false
	afterEvent := false

	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"before_event": func(e *Event) {
				beforeEvent = true
			},
			"leave_state": func(e *Event) {
				leaveState = true
			},
			"enter_state": func(e *Event) {
				enterState = true
			},
			"after_event": func(e *Event) {
				afterEvent = true
			},
		},
	)

	fsm.Event("run")
	if !(beforeEvent && leaveState && enterState && afterEvent) {
		t.Error("expected all callbacks to be called")
	}
}

func TestSpecificCallbacks(t *testing.T) {
	beforeEvent := false
	leaveState := false
	enterState := false
	afterEvent := false

	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"before_run": func(e *Event) {
				beforeEvent = true
			},
			"leave_start": func(e *Event) {
				leaveState = true
			},
			"enter_end": func(e *Event) {
				enterState = true
			},
			"after_run": func(e *Event) {
				afterEvent = true
			},
		},
	)

	fsm.Event("run")
	if !(beforeEvent && leaveState && enterState && afterEvent) {
		t.Error("expected all callbacks to be called")
	}
}

func TestSpecificCallbacksShortform(t *testing.T) {
	enterState := false
	afterEvent := false

	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"end": func(e *Event) {
				enterState = true
			},
			"run": func(e *Event) {
				afterEvent = true
			},
		},
	)

	fsm.Event("run")
	if !(enterState && afterEvent) {
		t.Error("expected all callbacks to be called")
	}
}

func TestBeforeEventWithoutTransition(t *testing.T) {
	beforeEvent := true

	fsm := NewFSM(
		"start",
		Events{
			{Name: "dontrun", Src: []string{"start"}, Dst: "start"},
		},
		Callbacks{
			"before_event": func(e *Event) {
				beforeEvent = true
			},
		},
	)

	err := fsm.Event("dontrun")
	if e, ok := err.(NoTransitionError); !ok && e.Err != nil {
		t.Error("expected 'NoTransitionError' without custom error")
	}

	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
	if !beforeEvent {
		t.Error("expected callback to be called")
	}
}

func TestCancelBeforeGenericEvent(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"before_event": func(e *Event) {
				e.Cancel()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}

func TestCancelBeforeSpecificEvent(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"before_run": func(e *Event) {
				e.Cancel()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}

func TestCancelLeaveGenericState(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"leave_state": func(e *Event) {
				e.Cancel()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}

func TestCancelLeaveSpecificState(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"leave_start": func(e *Event) {
				e.Cancel()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}

func TestCancelWithError(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"before_event": func(e *Event) {
				e.Cancel(fmt.Errorf("error"))
			},
		},
	)
	err := fsm.Event("run")
	if _, ok := err.(CanceledError); !ok {
		t.Error("expected only 'CanceledError'")
	}

	if e, ok := err.(CanceledError); ok && e.Err.Error() != "error" {
		t.Error("expected 'CanceledError' with correct custom error")
	}

	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}

func TestAsyncTransitionGenericState(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"leave_state": func(e *Event) {
				e.Async()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
	fsm.Transition()
	if fsm.Current() != "end" {
		t.Error("expected state to be 'end'")
	}
}

func TestAsyncTransitionSpecificState(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"leave_start": func(e *Event) {
				e.Async()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
	fsm.Transition()
	if fsm.Current() != "end" {
		t.Error("expected state to be 'end'")
	}
}

func TestAsyncTransitionInProgress(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
			{Name: "reset", Src: []string{"end"}, Dst: "start"},
		},
		Callbacks{
			"leave_start": func(e *Event) {
				e.Async()
			},
		},
	)
	fsm.Event("run")
	err := fsm.Event("reset")
	if e, ok := err.(InTransitionError); !ok && e.Event != "reset" {
		t.Error("expected 'InTransitionError' with correct state")
	}
	fsm.Transition()
	fsm.Event("reset")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}

func TestAsyncTransitionNotInProgress(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
			{Name: "reset", Src: []string{"end"}, Dst: "start"},
		},
		Callbacks{},
	)
	err := fsm.Transition()
	if _, ok := err.(NotInTransitionError); !ok {
		t.Error("expected 'NotInTransitionError'")
	}
}

func TestCallbackNoError(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"run": func(e *Event) {
			},
		},
	)
	e := fsm.Event("run")
	if e != nil {
		t.Error("expected no error")
	}
}

func TestCallbackError(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"run": func(e *Event) {
				e.Err = fmt.Errorf("error")
			},
		},
	)
	e := fsm.Event("run")
	if e.Error() != "error" {
		t.Error("expected error to be 'error'")
	}
}

func TestCallbackArgs(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"run": func(e *Event) {
				if len(e.Args) != 1 {
					t.Error("too few arguments")
				}
				arg, ok := e.Args[0].(string)
				if !ok {
					t.Error("not a string argument")
				}
				if arg != "test" {
					t.Error("incorrect argument")
				}
			},
		},
	)
	fsm.Event("run", "test")
}

func TestNoDeadLock(t *testing.T) {
	var fsm *FSM
	fsm = NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"run": func(e *Event) {
				fsm.Current() // Should not result in a panic / deadlock
			},
		},
	)
	fsm.Event("run")
}

func TestThreadSafetyRaceCondition(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"run": func(e *Event) {
			},
		},
	)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = fsm.Current()
	}()
	fsm.Event("run")
	wg.Wait()
}

func TestDoubleTransition(t *testing.T) {
	var fsm *FSM
	var wg sync.WaitGroup
	wg.Add(2)
	fsm = NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Callbacks{
			"before_run": func(e *Event) {
				wg.Done()
				// Imagine a concurrent event coming in of the same type while
				// the data access mutex is unlocked because the current transition
				// is running its event callbacks, getting around the "active"
				// transition checks
				if len(e.Args) == 0 {
					// Must be concurrent so the test may pass when we add a mutex that synchronizes
					// calls to Event(...). It will then fail as an inappropriate transition as we
					// have changed state.
					go func() {
						if err := fsm.Event("run", "second run"); err != nil {
							fmt.Println(err)
							wg.Done() // It should fail, and then we unfreeze the test.
						}
					}()
					time.Sleep(20 * time.Millisecond)
				} else {
					panic("Was able to reissue an event mid-transition")
				}
			},
		},
	)
	if err := fsm.Event("run"); err != nil {
		fmt.Println(err)
	}
	wg.Wait()
}

func TestNoTransition(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "start"},
		},
		Callbacks{},
	)
	err := fsm.Event("run")
	if _, ok := err.(NoTransitionError); !ok {
		t.Error("expected 'NoTransitionError'")
	}
}

func ExampleNewFSM() {
	fsm := NewFSM(
		"green",
		Events{
			{Name: "warn", Src: []string{"green"}, Dst: "yellow"},
			{Name: "panic", Src: []string{"yellow"}, Dst: "red"},
			{Name: "panic", Src: []string{"green"}, Dst: "red"},
			{Name: "calm", Src: []string{"red"}, Dst: "yellow"},
			{Name: "clear", Src: []string{"yellow"}, Dst: "green"},
		},
		Callbacks{
			"before_warn": func(e *Event) {
				fmt.Println("before_warn")
			},
			"before_event": func(e *Event) {
				fmt.Println("before_event")
			},
			"leave_green": func(e *Event) {
				fmt.Println("leave_green")
			},
			"leave_state": func(e *Event) {
				fmt.Println("leave_state")
			},
			"enter_yellow": func(e *Event) {
				fmt.Println("enter_yellow")
			},
			"enter_state": func(e *Event) {
				fmt.Println("enter_state")
			},
			"after_warn": func(e *Event) {
				fmt.Println("after_warn")
			},
			"after_event": func(e *Event) {
				fmt.Println("after_event")
			},
		},
	)
	fmt.Println(fsm.Current())
	err := fsm.Event("warn")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
	// Output:
	// green
	// before_warn
	// before_event
	// leave_green
	// leave_state
	// enter_yellow
	// enter_state
	// after_warn
	// after_event
	// yellow
}

func ExampleFSM_Current() {
	fsm := NewFSM(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Callbacks{},
	)
	fmt.Println(fsm.Current())
	// Output: closed
}

func ExampleFSM_Is() {
	fsm := NewFSM(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Callbacks{},
	)
	fmt.Println(fsm.Is("closed"))
	fmt.Println(fsm.Is("open"))
	// Output:
	// true
	// false
}

func ExampleFSM_Can() {
	fsm := NewFSM(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Callbacks{},
	)
	fmt.Println(fsm.Can("open"))
	fmt.Println(fsm.Can("close"))
	// Output:
	// true
	// false
}

func ExampleFSM_AvailableTransitions() {
	fsm := NewFSM(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
			{Name: "kick", Src: []string{"closed"}, Dst: "broken"},
		},
		Callbacks{},
	)
	// sort the results ordering is consistent for the output checker
	transitions := fsm.AvailableTransitions()
	sort.Strings(transitions)
	fmt.Println(transitions)
	// Output:
	// [kick open]
}

func ExampleFSM_Cannot() {
	fsm := NewFSM(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Callbacks{},
	)
	fmt.Println(fsm.Cannot("open"))
	fmt.Println(fsm.Cannot("close"))
	// Output:
	// false
	// true
}

func ExampleFSM_Event() {
	fsm := NewFSM(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Callbacks{},
	)
	fmt.Println(fsm.Current())
	err := fsm.Event("open")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
	err = fsm.Event("close")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
	// Output:
	// closed
	// open
	// closed
}

func ExampleFSM_Transition() {
	fsm := NewFSM(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Callbacks{
			"leave_closed": func(e *Event) {
				e.Async()
			},
		},
	)
	err := fsm.Event("open")
	if e, ok := err.(AsyncError); !ok && e.Err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
	err = fsm.Transition()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
	// Output:
	// closed
	// open
}
