package dao

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetUserByUserID(t *testing.T) {
	var (
		c            = context.TODO()
		userID int64 = 123
	)
	convey.Convey("GetUserByMid", t, func(ctx convey.C) {
		res, err := d.GetUserByMid(c, userID)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}

func TestDaoGetUserBaseByUserID(t *testing.T) {
	var (
		c      = context.TODO()
		userID = "bishi"
	)
	convey.Convey("GetUserBaseByUserID", t, func(ctx convey.C) {
		res, err := d.GetUserBaseByUserID(c, userID)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}
func TestDaoGetUserTelByTel(t *testing.T) {
	var (
		c   = context.TODO()
		tel = []byte{123}
	)
	convey.Convey("GetUserTelByTel", t, func(ctx convey.C) {
		res, err := d.GetUserTelByTel(c, tel)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}

func TestDaoGetUserMailByMail(t *testing.T) {
	var (
		c    = context.TODO()
		mail = []byte{123}
	)
	convey.Convey("GetUserMailByMail", t, func(ctx convey.C) {
		res, err := d.GetUserMailByMail(c, mail)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})

		str, _ := json.Marshal(res)
		t.Logf("res: %s", str)
	})
}
