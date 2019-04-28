package lrulist

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	l := New()
	testAppend(l, t)
	testPreAppend(l, t)
	testFind(l, t)
	testRev(l, t)
	testMoveToHead(l, t)
	testEmpty(l, t)
}

// 测试尾插
func testAppend(l *List, t *testing.T) {
	t.Log("testAppend")
	l.Append("chenxuewen1", "1")
	l.Append("chenxuewen2", "2")
	l.Append("chenxuewen3", "3")
	l.Append("chenxuewen4", "4")
	l.Append("chenxuewen5", "5")
	l.Each(func(node Node) {
		t.Log(node)
	})
}

// 测试头插
func testPreAppend(l *List, t *testing.T) {
	t.Log("testPreAppend")
	l.Prepend("chenxuewen0", "0")
	l.Each(func(node Node) {
		t.Log(node)
	})
}

// 测试查找
func testFind(l *List, t *testing.T) {
	t.Log("testFind")
	t.Log(l.Find("chenxuewen1"))
}

// 测试删除
func testRev(l *List, t *testing.T) {
	t.Log("testRev")
	l.Remove(l.Find("chenxuewen2"))
	l.Each(func(node Node) {
		t.Log(node)
	})
}

// 测试移动到头部
func testMoveToHead(l *List, t *testing.T) {
	t.Log("testMoveToHead")
	l.MoveToHead(l.Remove(l.Find("chenxuewen4")))
	l.Each(func(node Node) {
		t.Log(node)
	})
}

// 测试清空
func testEmpty(l *List, t *testing.T) {
	t.Log("testEmpty")
	l.Empty()
	l.Each(func(node Node) {
		t.Log(node)
	})
}

func BenchmarkListLRU(b *testing.B) {
	l := New()
	lock := sync.RWMutex{}
	l.Append("chenxuewen1", "1")
	l.Append("chenxuewen2", "2")
	l.Append("chenxuewen3", "3")
	l.Append("chenxuewen4", "4")
	l.Append("chenxuewen5", "5")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lock.RLock()
			n1 := l.Find("chenxuewen3")
			lock.RUnlock()
			lock.Lock()
			if n1 != nil {
				l.MoveToHead(l.Remove(n1))
			}
			lock.Unlock()

			lock.RLock()
			n2 := l.Find("chenxuewen4")
			lock.RUnlock()
			lock.Lock()
			if n2 != nil {
				l.MoveToHead(l.Remove(n2))
			}
			lock.Unlock()
		}
	})
}
