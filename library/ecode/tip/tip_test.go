package tip

import (
	"sync"
	"testing"

	"go-common/library/ecode"
)

var (
	once sync.Once
)

func initEcodes() {
	once.Do(func() {
		Init(nil)
	})
}

func TestInit(t *testing.T) {
	initEcodes()
	testCodes(t)
}

func BenchmarkLookup(b *testing.B) {
	initEcodes()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.Logf("Ecodes: lookup ServerErr: %s", ecode.ServerErr.Message())
			b.Logf("Ecodes: lookup ServerErr: %s", ecode.NotModified.Message())
		}
	})
}

func testCodes(t *testing.T) {
	if ver, err := defualtEcodes.update(1499843315); err != nil {
		t.Logf("codes(%v)", err)
		t.FailNow()
	} else {
		t.Logf("ver(%d)", ver)
	}
	if codes, ok := defualtEcodes.codes.Load().(map[int]string); !ok {
		t.Errorf("codes load not ok")
		t.FailNow()
	} else {
		t.Logf("%v", codes)
	}
}
