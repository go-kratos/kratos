package service

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/library/ecode"
)

func TestService_Oauth(t *testing.T) {
	once.Do(startService)
	ak := "4de8aecfafc7f91cb650d6371efb1b63"
	expectMid := int64(110000139)
	if res, err := s.Oauth(context.TODO(), s.appMap[_gameAppKey], ak, ""); err != nil {
		t.Errorf("s.Oauth() error(%v)", err)
		t.FailNow()
	} else if res == nil || res.Mid != expectMid {
		t.Errorf("res is not correct, expected res with mid %d but got %v", expectMid, res)
		t.FailNow()
	} else {
		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	}
}

func TestService_Oauth_Expires(t *testing.T) {
	once.Do(startService)
	ak := "c1e220e12fd4c89a0c5449b9f8c7b062"
	if _, err := s.Oauth(context.TODO(), s.appMap[_gameAppKey], ak, ""); err != ecode.AccessTokenExpires {
		t.Errorf("res is not correct, expected error %v, but got %v", ecode.AccessTokenExpires, err)
		t.FailNow()
	}
}
