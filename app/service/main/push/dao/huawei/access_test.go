package huawei

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewAccess(t *testing.T) {
	Convey("new access", t, func() {
		ac, err := NewAccess("10125085", "iejq6hn3ds3d4neq1m21v443lmbm31gs")
		if err != nil {
			t.Errorf("new access error(%v)", err)
		} else {
			t.Log(ac.Token, ac.Expire)
		}
	})
}

func Test_AccessExpire(t *testing.T) {
	Convey("access expire", t, func() {
		ac := Access{Expire: time.Now().Add(-8 * time.Hour).Unix()}
		if !ac.IsExpired() {
			t.Errorf("access should be expire")
		}
		ac.Expire -= 10
		if ac.IsExpired() {
			t.Error("access should not be expire")
		}
	})
}
