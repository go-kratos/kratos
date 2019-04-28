package service

import (
	"context"
	"encoding/json"
	"testing"
)

const (
	_gameAppKey = "d4761645d1632e8c"
)

// TestService_MyInfo oauth via origin.
func TestService_MyInfo(t *testing.T) {
	once.Do(startService)
	ak := "4de8aecfafc7f91cb650d6371efb1b63"
	expectMid := int64(110000139)
	if res, err := s.MyInfo(context.TODO(), s.appMap[_gameAppKey], ak); err != nil {
		t.Errorf("s.MyInfo() error(%v)", err)
		t.FailNow()
	} else if res == nil || res.Mid != expectMid {
		t.Errorf("res is not correct, expected res with mid %d but got %v", expectMid, res)
		t.FailNow()
	} else {
		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	}
}
