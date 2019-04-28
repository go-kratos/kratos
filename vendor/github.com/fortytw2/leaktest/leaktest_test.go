package leaktest

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

type testReporter struct {
	failed bool
	msg    string
}

func (tr *testReporter) Errorf(format string, args ...interface{}) {
	tr.failed = true
	tr.msg = fmt.Sprintf(format, args...)
}

var leakyFuncs = []func(){
	// Infinite for loop
	func() {
		for {
			time.Sleep(time.Second)
		}
	},
	// Select on a channel not referenced by other goroutines.
	func() {
		c := make(chan struct{})
		<-c
	},
	// Blocked select on channels not referenced by other goroutines.
	func() {
		c := make(chan struct{})
		c2 := make(chan struct{})
		select {
		case <-c:
		case c2 <- struct{}{}:
		}
	},
	// Blocking wait on sync.Mutex that isn't referenced by other goroutines.
	func() {
		var mu sync.Mutex
		mu.Lock()
		mu.Lock()
	},
	// Blocking wait on sync.RWMutex that isn't referenced by other goroutines.
	func() {
		var mu sync.RWMutex
		mu.RLock()
		mu.Lock()
	},
	func() {
		var mu sync.Mutex
		mu.Lock()
		c := sync.NewCond(&mu)
		c.Wait()
	},
}

func TestCheck(t *testing.T) {
	// this works because the running goroutine is left running at the
	// start of the next test case - so the previous leaks don't affect the
	// check for the next one
	for i, fn := range leakyFuncs {
		checker := &testReporter{}
		snapshot := CheckTimeout(checker, time.Second)
		go fn()

		snapshot()
		if !checker.failed {
			t.Errorf("didn't catch sleeping goroutine, test #%d", i)
		}
	}
}

// TestSlowTest verifies that the timeout works on slow tests: it should
// be based on time after the test finishes rather than time after the test's
// start.
func TestSlowTest(t *testing.T) {
	defer CheckTimeout(t, 1000 * time.Millisecond)()

	go time.Sleep(1500 * time.Millisecond)
	time.Sleep(750 * time.Millisecond)
}

func TestEmptyLeak(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer CheckContext(ctx, t)()
	time.Sleep(time.Second)
}

// TestChangingStackTrace validates that a change in a preexisting goroutine's
// stack is not detected as a leaked goroutine.
func TestChangingStackTrace(t *testing.T) {
	started := make(chan struct{})
	c1 := make(chan struct{})
	c2 := make(chan struct{})
	defer close(c2)
	go func() {
		close(started)
		<-c1
		<-c2
	}()
	<-started
	func() {
		defer CheckTimeout(t, time.Second)()
		close(c1)
	}()
}

func TestInterestingGoroutine(t *testing.T) {
	s := "goroutine 123 [running]:\nmain.main()"
	gr, err := interestingGoroutine(s)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if gr.id != 123 {
		t.Errorf("goroutine id = %d; want %d", gr.id, 123)
	}
	if gr.stack != s {
		t.Errorf("goroutine stack = %q; want %q", gr.stack, s)
	}

	stacks := []struct {
		stack string
		err   error
	}{
		{
			stack: "goroutine 123 [running]:",
			err:   errors.New(`error parsing stack: "goroutine 123 [running]:"`),
		},
		{
			stack: "goroutine 123 [running]:\ntesting.RunTests",
			err:   nil,
		},
		{
			stack: "goroutine 123 [running]:\nfoo\nbar\nruntime.goexit\nbaz",
			err:   nil,
		},
		{
			stack: "goroutine 123:\nmain.main()",
			err:   errors.New(`error parsing stack header: "goroutine 123:"`),
		},
		{
			stack: "goroutine NaN [running]:\nmain.main()",
			err:   errors.New(`error parsing goroutine id: strconv.ParseUint: parsing "NaN": invalid syntax`),
		},
	}
	for i, s := range stacks {
		gr, err := interestingGoroutine(s.stack)
		if s.err == nil && err != nil {
			t.Errorf("%d: error = %v; want nil", i, err)
		} else if s.err != nil && (err == nil || err.Error() != s.err.Error()) {
			t.Errorf("%d: error = %v; want %s", i, err, s.err)
		}
		if gr != nil {
			t.Errorf("%d: gr = %v; want nil", i, gr)
		}

	}
}
