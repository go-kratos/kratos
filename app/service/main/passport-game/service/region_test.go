package service

import (
	"context"
	"encoding/json"
	"testing"
)

func TestService_Regions(t *testing.T) {
	once.Do(startService)
	m := s.Regions(context.TODO())
	res, _ := json.Marshal(m)
	t.Logf("res: %s", res)
}

func TestService_Region(t *testing.T) {
	ak := "4de8aecfafc7f91cb650d6371efb1b63"
	r, _ := region(ak)
	if r != _origin {
		t.FailNow()
	}
	t.Logf("ak: %s, region: %s", ak, r)

	ak = "4de8aecfafc7f91cb650d6371efb1b63_t1"
	r, _ = region(ak)
	if r != "t1" {
		t.FailNow()
	}
	t.Logf("ak: %s, region: %s", ak, r)

	ak = "4de8aecfafc7f91cb650d6371efb1b63_"
	_, ok := region(ak)
	if ok {
		t.FailNow()
	}
	t.Logf("ak: %s, ok: %t", ak, ok)
}
