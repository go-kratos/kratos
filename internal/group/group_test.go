package group

import (
	"reflect"
	"testing"
)

func TestGroupGet(t *testing.T) {
	count := 0
	g := NewGroup[int](func() int {
		count++
		return count
	})
	v := g.Get("key_0")
	if !reflect.DeepEqual(v, 1) {
		t.Errorf("expect 1, actual %v", v)
	}

	v = g.Get("key_1")
	if !reflect.DeepEqual(v, 2) {
		t.Errorf("expect 2, actual %v", v)
	}

	v = g.Get("key_0")
	if !reflect.DeepEqual(v, 1) {
		t.Errorf("expect 1, actual %v", v)
	}
	if !reflect.DeepEqual(count, 2) {
		t.Errorf("expect count 2, actual %v", count)
	}
}

func TestGroupReset(t *testing.T) {
	g := NewGroup(func() int {
		return 1
	})
	g.Get("key")
	call := false
	g.Reset(func() int {
		call = true
		return 1
	})

	length := 0
	for range g.vals {
		length++
	}
	if !reflect.DeepEqual(length, 0) {
		t.Errorf("expect length 0, actual %v", length)
	}

	g.Get("key")
	if !reflect.DeepEqual(call, true) {
		t.Errorf("expect call true, actual %v", call)
	}
}

func TestGroupClear(t *testing.T) {
	g := NewGroup(func() int {
		return 1
	})
	g.Get("key")
	length := 0
	for range g.vals {
		length++
	}
	if !reflect.DeepEqual(length, 1) {
		t.Errorf("expect length 1, actual %v", length)
	}

	g.Clear()
	length = 0
	for range g.vals {
		length++
	}
	if !reflect.DeepEqual(length, 0) {
		t.Errorf("expect length 0, actual %v", length)
	}
}
