package netutil

import "testing"

func TestListen(t *testing.T) {
	l := limitListener{cur: 0, max: 10, sem: make(chan struct{}, 10)}
	ok := l.acquire()
	t.Logf("LimitListener: acquire status: %t", ok)
	if l.cur != 1 {
		t.Errorf("LimitListener: expceted cur=1 but got %d", l.cur)
	}
	l.release()
	if l.cur != 0 {
		t.Errorf("LimitListener: expceted cur=0 but got %d", l.cur)
	}
}
