package lrucache

import (
	"container/list"
	"testing"
)

type Elem struct {
	key   int
	value string
}

func Test_New(t *testing.T) {
	lc := New(5)
	if lc.Len() != 0 {
		t.Error("case 1 failed")
	}
}

func Test_Put(t *testing.T) {
	lc := New(0)
	lc.Put(1, "1")
	if lc.Len() != 0 {
		t.Error("case 1.1 failed")
	}

	lc = New(5)
	lc.Put(1, "1")
	lc.Put(2, "2")
	lc.Put(1, "3")
	if lc.Len() != 2 {
		t.Error("case 2.1 failed")
	}

	l := list.New()
	l.PushBack(&Elem{1, "3"})
	l.PushBack(&Elem{2, "2"})

	e := l.Front()
	for c := lc.Front(); c != nil; c = c.Next() {
		v := e.Value.(*Elem)
		if c.Key.(int) != v.key {
			t.Error("case 2.2 failed: ", c.Key.(int), v.key)
		}
		if c.Value.(string) != v.value {
			t.Error("case 2.3 failed: ", c.Value.(string), v.value)
		}
		e = e.Next()
	}

	lc.Put(3, "4")
	lc.Put(4, "5")
	lc.Put(5, "6")
	lc.Put(2, "7")
	if lc.Len() != 5 {
		t.Error("case 3.1 failed")
	}

	l = list.New()
	l.PushBack(&Elem{2, "7"})
	l.PushBack(&Elem{5, "6"})
	l.PushBack(&Elem{4, "5"})
	l.PushBack(&Elem{3, "4"})
	l.PushBack(&Elem{1, "3"})

	rl := list.New()
	rl.PushBack(&Elem{1, "3"})
	rl.PushBack(&Elem{3, "4"})
	rl.PushBack(&Elem{4, "5"})
	rl.PushBack(&Elem{5, "6"})
	rl.PushBack(&Elem{2, "7"})

	e = l.Front()
	for c := lc.Front(); c != nil; c = c.Next() {
		v := e.Value.(*Elem)
		if c.Key.(int) != v.key {
			t.Error("case 3.2 failed: ", c.Key.(int), v.key)
		}
		if c.Value.(string) != v.value {
			t.Error("case 3.3 failed: ", c.Value.(string), v.value)
		}
		e = e.Next()
	}

	e = rl.Front()
	for c := lc.Back(); c != nil; c = c.Prev() {
		v := e.Value.(*Elem)
		if c.Key.(int) != v.key {
			t.Error("case 3.4 failed: ", c.Key.(int), v.key)
		}
		if c.Value.(string) != v.value {
			t.Error("case 3.5 failed: ", c.Value.(string), v.value)
		}
		e = e.Next()
	}

	lc.Put(6, "8")
	if lc.Len() != 5 {
		t.Error("case 4.1 failed")
	}

	l = list.New()
	l.PushBack(&Elem{6, "8"})
	l.PushBack(&Elem{2, "7"})
	l.PushBack(&Elem{5, "6"})
	l.PushBack(&Elem{4, "5"})
	l.PushBack(&Elem{3, "4"})

	e = l.Front()
	for c := lc.Front(); c != nil; c = c.Next() {
		v := e.Value.(*Elem)
		if c.Key.(int) != v.key {
			t.Error("case 4.2 failed: ", c.Key.(int), v.key)
		}
		if c.Value.(string) != v.value {
			t.Error("case 4.3 failed: ", c.Value.(string), v.value)
		}
		e = e.Next()
	}
}

func Test_Get(t *testing.T) {
	lc := New(2)
	lc.Put(1, "1")
	lc.Put(2, "2")
	if v, _ := lc.Get(1); v != "1" {
		t.Error("case 1.1 failed")
	}
	lc.Put(3, "3")
	if lc.Len() != 2 {
		t.Error("case 1.2 failed")
	}

	l := list.New()
	l.PushBack(&Elem{3, "3"})
	l.PushBack(&Elem{1, "1"})

	e := l.Front()
	for c := lc.Front(); c != nil; c = c.Next() {
		v := e.Value.(*Elem)
		if c.Key.(int) != v.key {
			t.Error("case 1.3 failed: ", c.Key.(int), v.key)
		}
		if c.Value.(string) != v.value {
			t.Error("case 1.4 failed: ", c.Value.(string), v.value)
		}
		e = e.Next()
	}
}

func Test_Delete(t *testing.T) {
	lc := New(5)
	lc.Put(3, "4")
	lc.Put(4, "5")
	lc.Put(5, "6")
	lc.Put(2, "7")
	lc.Put(6, "8")
	lc.Delete(5)

	l := list.New()
	l.PushBack(&Elem{6, "8"})
	l.PushBack(&Elem{2, "7"})
	l.PushBack(&Elem{4, "5"})
	l.PushBack(&Elem{3, "4"})
	if lc.Len() != 4 {
		t.Error("case 1.1 failed")
	}

	e := l.Front()
	for c := lc.Front(); c != nil; c = c.Next() {
		v := e.Value.(*Elem)
		if c.Key.(int) != v.key {
			t.Error("case 1.2 failed: ", c.Key.(int), v.key)
		}
		if c.Value.(string) != v.value {
			t.Error("case 1.3 failed: ", c.Value.(string), v.value)
		}
		e = e.Next()
	}

	lc.Delete(6)

	l = list.New()
	l.PushBack(&Elem{2, "7"})
	l.PushBack(&Elem{4, "5"})
	l.PushBack(&Elem{3, "4"})
	if lc.Len() != 3 {
		t.Error("case 2.1 failed")
	}

	e = l.Front()
	for c := lc.Front(); c != nil; c = c.Next() {
		v := e.Value.(*Elem)
		if c.Key.(int) != v.key {
			t.Error("case 2.2 failed: ", c.Key.(int), v.key)
		}
		if c.Value.(string) != v.value {
			t.Error("case 2.3 failed: ", c.Value.(string), v.value)
		}
		e = e.Next()
	}

	lc.Delete(3)

	l = list.New()
	l.PushBack(&Elem{2, "7"})
	l.PushBack(&Elem{4, "5"})
	if lc.Len() != 2 {
		t.Error("case 3.1 failed")
	}

	e = l.Front()
	for c := lc.Front(); c != nil; c = c.Next() {
		v := e.Value.(*Elem)
		if c.Key.(int) != v.key {
			t.Error("case 3.2 failed: ", c.Key.(int), v.key)
		}
		if c.Value.(string) != v.value {
			t.Error("case 3.3 failed: ", c.Value.(string), v.value)
		}
		e = e.Next()
	}
}

func Test_Range(t *testing.T) {
	lc := New(5)
	lc.Put(3, "4")
	lc.Put(4, "5")
	lc.Put(5, "6")
	lc.Put(2, "7")
	lc.Put(6, "8")

	l := list.New()
	l.PushBack(&Elem{6, "8"})
	l.PushBack(&Elem{2, "7"})
	l.PushBack(&Elem{5, "6"})
	l.PushBack(&Elem{4, "5"})
	l.PushBack(&Elem{3, "4"})

	e := l.Front()
	lc.Range(
		func(key, value interface{}) bool {
			v := e.Value.(*Elem)
			if key.(int) != v.key {
				t.Error("case 1.1 failed: ", key.(int), v.key)
			}
			if value.(string) != v.value {
				t.Error("case 1.2 failed: ", value.(string), v.value)
			}
			e = e.Next()
			return true
		})

	if e != nil {
		t.Error("case 1.3 failed: ", e.Value)
	}
}
