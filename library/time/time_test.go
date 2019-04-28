package time

import (
	"context"
	"testing"
	"time"
)

func TestShrink(t *testing.T) {
	var d Duration
	err := d.UnmarshalText([]byte("1s"))
	if err != nil {
		t.Fatalf("TestShrink:  d.UnmarshalText failed!err:=%v", err)
	}
	c := context.Background()
	to, ctx, cancel := d.Shrink(c)
	defer cancel()
	if time.Duration(to) != time.Second {
		t.Fatalf("new timeout must be equal 1 second")
	}
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > time.Second || time.Until(deadline) < time.Millisecond*500 {
		t.Fatalf("ctx deadline must be less than 1s and greater than 500ms")
	}
}

func TestShrinkWithTimeout(t *testing.T) {
	var d Duration
	err := d.UnmarshalText([]byte("1s"))
	if err != nil {
		t.Fatalf("TestShrink:  d.UnmarshalText failed!err:=%v", err)
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	to, ctx, cancel := d.Shrink(c)
	defer cancel()
	if time.Duration(to) != time.Second {
		t.Fatalf("new timeout must be equal 1 second")
	}
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > time.Second || time.Until(deadline) < time.Millisecond*500 {
		t.Fatalf("ctx deadline must be less than 1s and greater than 500ms")
	}
}

func TestShrinkWithDeadline(t *testing.T) {
	var d Duration
	err := d.UnmarshalText([]byte("1s"))
	if err != nil {
		t.Fatalf("TestShrink:  d.UnmarshalText failed!err:=%v", err)
	}
	c, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	to, ctx, cancel := d.Shrink(c)
	defer cancel()
	if time.Duration(to) >= time.Millisecond*500 {
		t.Fatalf("new timeout must be less than 500 ms")
	}
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > time.Millisecond*500 || time.Until(deadline) < time.Millisecond*200 {
		t.Fatalf("ctx deadline must be less than 500ms and greater than 200ms")
	}
}
