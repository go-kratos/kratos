package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/library/ecode"
	"testing"
)

func queryDelCache(t *testing.T, uid int64) (success bool, err error) {

	bp := getTestDefaultBasicParam("")
	_, err = s.DelCache(ctx, bp, uid, "")
	if err == nil {
		success = true
	}
	return
}

func TestService_DelCache(t *testing.T) {
	Convey("not found", t, testWithRandUid(func(uid int64) {
		queryDelCache(t, uid)

		s, e := queryDelCache(t, uid)
		So(s, ShouldEqual, false)
		So(e, ShouldEqual, ecode.NothingFound)

	}))

	Convey("normal", t, testWithTestUser(func(u *TestUser) {
		getTestWallet(t, u.uid, "pc")
		s, e := queryDelCache(t, u.uid)
		So(s, ShouldEqual, true)
		So(e, ShouldBeNil)
	}))
}
