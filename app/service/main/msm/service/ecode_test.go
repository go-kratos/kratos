package service

import (
	"context"
	"testing"

	"go-common/app/service/main/msm/model"
)

func TestService(t *testing.T) {
	testIncUpdate(t, svr)
	testCodes(t, svr)
	testUpdate(t, svr)
	testAllCodes(t, svr)
}
func testIncUpdate(t *testing.T, svr *Service) {
	codes := svr.codes.Load().(*model.Codes)
	codes.Ver = 1499742647
	if err := svr.diff(); err != nil {
		t.Logf("update(%v)", err)
		t.FailNow()
	}
}

func testCodes(t *testing.T, svr *Service) {
	if code, err := svr.Codes(context.TODO(), 1499742647); err != nil {
		t.Logf("codes() error(%v)", err)
		t.FailNow()
	} else {
		t.Logf("update() data(%v) ", code)
	}
}
func testUpdate(t *testing.T, svr *Service) {
	if err := svr.all(); err != nil {
		t.Logf("update(%v)", err)
		t.FailNow()
	}
}

func testAllCodes(t *testing.T, svr *Service) {
	if code, err := svr.Codes(context.TODO(), 0); err != nil {
		t.Logf("codes() error(%v)", err)
		t.FailNow()
	} else {
		t.Logf("update() data(%v)", code)
	}

}
func BenchmarkAllCodes(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := svr.Codes(context.TODO(), 0); err != nil {
				b.Logf("codes() error(%v)", err)
				b.FailNow()
			}
		}
	})
}
func BenchmarkCodes(b *testing.B) {
	//svr.all()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			//svr.diff()
			if _, err := svr.Codes(context.TODO(), svr.codes.Load().(*model.Codes).Ver); err != nil {
				b.Logf("codes() error(%v)", err)
				b.FailNow()
			}
		}
	})
}
