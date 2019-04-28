package client

import (
	"context"

	"testing"
	"time"

	"go-common/app/service/main/seq-server/model"
)

func TestDynamic(t *testing.T) {
	s := New2(nil)
	time.Sleep(5 * time.Second)
	testID(t, s)
	testID32(t, s)
}

func testID(t *testing.T, s *Service2) {
	res := make(map[int64]struct{})
	for i := 0; i < 10000; i++ {
		id, err := s.ID(context.TODO(), &model.ArgBusiness{BusinessID: 7, Token: "RA8yy0RjDCBTGgFUha4hPOnhxfXvM8hR"})
		if err != nil {
			t.Errorf("s.ID error(%v)", err)
			continue
		}
		if _, ok := res[id]; ok {
			t.Errorf("s.ID repeat id:%d", id)
			t.FailNow()
		}
		res[id] = struct{}{}
		t.Logf("got ID(%d)", id)
	}
}

func testID32(t *testing.T, s *Service2) {
	res := make(map[int32]struct{})
	for i := 0; i < 10000; i++ {
		id, err := s.ID32(context.TODO(), &model.ArgBusiness{BusinessID: 7, Token: "RA8yy0RjDCBTGgFUha4hPOnhxfXvM8hR"})
		if err != nil {
			t.Errorf("s.ID error(%v)", err)
			continue
		}
		if _, ok := res[id]; ok {
			t.Errorf("s.ID repeat id:%d", id)
			t.FailNow()
		}
		res[id] = struct{}{}
		t.Logf("got ID(%d)", id)
	}
}
