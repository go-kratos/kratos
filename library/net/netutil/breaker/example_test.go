package breaker_test

import (
	"fmt"
	"time"

	"go-common/library/net/netutil/breaker"
	xtime "go-common/library/time"
)

// ExampleGroup show group usage.
func ExampleGroup() {
	c := &breaker.Config{
		Window:  xtime.Duration(3 * time.Second),
		Sleep:   xtime.Duration(100 * time.Millisecond),
		Bucket:  10,
		Ratio:   0.5,
		Request: 100,
	}
	// init default config
	breaker.Init(c)
	// new group
	g := breaker.NewGroup(c)
	// reload group config
	c.Bucket = 100
	c.Request = 200
	g.Reload(c)
	// get breaker by key
	g.Get("key")
}

// ExampleBreaker show breaker usage.
func ExampleBreaker() {
	// new group,use default breaker config
	g := breaker.NewGroup(nil)
	brk := g.Get("key")
	// mark request success
	brk.MarkSuccess()
	// mark request failed
	brk.MarkFailed()
	// check if breaker allow or not
	if brk.Allow() == nil {
		fmt.Println("breaker allow")
	} else {
		fmt.Println("breaker not allow")
	}
}

// ExampleGo this example create a default group and show function callback
// according to the state of breaker.
func ExampleGo() {
	run := func() error {
		return nil
	}
	fallback := func() error {
		return fmt.Errorf("unknown error")
	}
	if err := breaker.Go("example_go", run, fallback); err != nil {
		fmt.Println(err)
	}
}
