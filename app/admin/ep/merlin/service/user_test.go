package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_QueryUserInfo(t *testing.T) {
	Convey("query user when user exists", t, func() {
		userInfo, err := s.QueryUserInfo(context.Background(), "fengyifeng")
		So(err, ShouldBeNil)
		So(userInfo.EMail, ShouldEqual, "fengyifeng@bilibili.com")
	})
}
