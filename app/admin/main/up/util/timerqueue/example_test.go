package timerqueue_test

import (
	"fmt"
	"time"

	"go-common/app/admin/main/up/util/timerqueue"
)

type event int

func (e event) OnTimer(t time.Time) {
	fmt.Printf("  Event %d executed at %v\n", int(e), t)
}

// Schedule several events with a timerqueue, and dispatch
// them by calling Advance.
func ExampleQueue() {
	queue := timerqueue.New()

	// Schedule an event each day from Jan 1 to Jan 7, 2015.
	tm := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 7; i++ {
		queue.Schedule(event(i), tm)
		tm = tm.Add(24 * time.Hour)
	}

	fmt.Println("Advancing to Jan 4...")
	queue.Advance(time.Date(2015, 1, 4, 0, 0, 0, 0, time.UTC))

	fmt.Println("Advancing to Jan 10...")
	queue.Advance(time.Date(2015, 1, 10, 0, 0, 0, 0, time.UTC))

	// Output:
	// Advancing to Jan 4...
	//   Event 1 executed at 2015-01-01 00:00:00 +0000 UTC
	//   Event 2 executed at 2015-01-02 00:00:00 +0000 UTC
	//   Event 3 executed at 2015-01-03 00:00:00 +0000 UTC
	//   Event 4 executed at 2015-01-04 00:00:00 +0000 UTC
	// Advancing to Jan 10...
	//   Event 5 executed at 2015-01-05 00:00:00 +0000 UTC
	//   Event 6 executed at 2015-01-06 00:00:00 +0000 UTC
	//   Event 7 executed at 2015-01-07 00:00:00 +0000 UTC
}
