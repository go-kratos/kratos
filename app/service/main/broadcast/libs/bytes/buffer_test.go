package bytes

import (
	"testing"
)

func TestBuffer(t *testing.T) {
	p := NewPool(2, 10)
	b := p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}
	b = p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}
	b = p.Get()
	if b.Bytes() == nil || len(b.Bytes()) == 0 {
		t.FailNow()
	}
}
