package fanout

import (
	"context"
	"testing"
	"time"
)

func TestFanout_Do(t *testing.T) {
	ca := New("cache", Worker(1), Buffer(1024))
	var run bool
	ca.Do(context.Background(), func(c context.Context) {
		run = true
		panic("error")
	})
	time.Sleep(time.Millisecond * 50)
	t.Log("not panic")
	if !run {
		t.Fatal("expect run be true")
	}
}

func TestFanout_Close(t *testing.T) {
	ca := New("cache", Worker(1), Buffer(1024))
	ca.Close()
	err := ca.Do(context.Background(), func(c context.Context) {})
	if err == nil {
		t.Fatal("expect get err")
	}
}
