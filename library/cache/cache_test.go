package cache

import (
	"testing"
	"time"
)

func TestCache_Save(t *testing.T) {
	ca := New(1, 1024)
	var run bool
	ca.Save(func() {
		run = true
		panic("error")
	})
	time.Sleep(time.Millisecond * 50)
	t.Log("don't panic")
	if !run {
		t.Fatal("expect run be true")
	}
}
