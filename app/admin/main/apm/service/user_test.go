package service

import (
	"context"
	"reflect"
	"testing"

	"go-common/app/admin/main/apm/dao"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestUser(t *testing.T) {
	convey.Convey("user", t, func() {
		t.Log("user test")
	})
}
func TestServiceGetUser(t *testing.T) {
	convey.Convey("GetUser", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			username = "chenjianrong"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(svr.dao), "GitLabFace", func(_ *dao.Dao, _ context.Context, _ string) (string, error) {
				return "http://git.bilibili.co/some/pic.jpg", nil
			})
			defer guard.Unpatch()
			usr, err := svr.GetUser(c, username)
			ctx.Convey("Then err should be nil.usr should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(usr, convey.ShouldNotBeNil)
			})
		})
	})
}
