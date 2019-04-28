package service

import (
	"context"
	"testing"

	"encoding/json"
)

func TestService_RenewToken(t *testing.T) {
	once.Do(startService)
	ak := "4de8aecfafc7f91cb650d6371efb1b63"
	if res, err := s.RenewToken(context.TODO(), ak, ""); err != nil {
		t.Errorf("s.RenewToken() error(%v)", err)
		t.FailNow()
	} else if res == nil || res.Expires == 0 {
		t.Errorf("res is not correct, expected res with expires non zero but got %v", res)
		t.FailNow()
	} else {
		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	}
}
