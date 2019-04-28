[![Build Status](https://travis-ci.org/beevik/timerqueue.svg?branch=master)](https://travis-ci.org/beevik/timerqueue)
[![GoDoc](https://godoc.org/github.com/beevik/timerqueue?status.svg)](https://godoc.org/github.com/beevik/timerqueue)

timerqueue
==========

The timerqueue package implements a priority queue for objects scheduled to
perform actions at clock times.

See http://godoc.org/github.com/beevik/timerqueue for godoc-formatted API
documentation.

###Example: Scheduling timers

The following code declares an object implementing the Timer interface,
creates a timerqueue, and adds three events to the timerqueue.

```go
type event int

func (e event) OnTimer(t time.Time) {
    fmt.Printf("event.OnTimer %d fired at %v\n", int(e), t)
}

queue := timerqueue.New()
queue.Schedule(event(1), time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC))
queue.Schedule(event(2), time.Date(2015, 1, 3, 0, 0, 0, 0, time.UTC))
queue.Schedule(event(3), time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC))

```

###Example: Peeking at the next timer to be scheduled

Using the queue initialized in the first example, the following code
examines the head of the timerqueue and outputs the id and time of
the event found there.

```go
e, t := queue.PeekFirst()
if e != nil {
    fmt.Printf("Event %d will be first to fire at %v.\n", int(e.(event)), t)
    fmt.Printf("%d events remain in the timerqueue.", queue.Len())
}
```

Output:
```
Event 1 will be first to fire at 2015-01-01 00:00:00 +0000 UTC.
3 events remain in the timerqueue.
```

###Example: Popping the next timer to be scheduled

Using the queue initialized in the first example, this code
removes the next timer to be executed until the queue is empty.

```go
for queue.Len() > 0 {
    e, t := queue.PopFirst()
    fmt.Printf("Event %d fires at %v.\n", int(e.(event)), t)
}
```

Output:
```
Event 1 fires at 2015-01-01 00:00:00 +0000 UTC.
Event 3 fires at 2015-01-02 00:00:00 +0000 UTC.
Event 2 fires at 2015-01-03 00:00:00 +0000 UTC.
```

###Example: Issuing OnTimer callbacks with Advance

The final example shows how to dispatch OnTimer callbacks to
timers using the timerqueue's Advance method.

Advance calls the OnTimer method for each timer scheduled
before the requested time. Timers are removed from the timerqueue
in order of their scheduling.

```go
// Call the OnTimer method for each event scheduled before
// January 10, 2015. Pop the called timer from the queue.
queue.Advance(time.Date(2015, 1, 10, 0, 0, 0, 0, time.UTC))
```

Output:
```
event.OnTimer 1 fired at 2015-01-01 00:00:00 +0000 UTC.
event.OnTimer 3 fired at 2015-01-02 00:00:00 +0000 UTC.
event.OnTimer 2 fired at 2015-01-03 00:00:00 +0000 UTC.
```
