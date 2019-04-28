package http

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestHttpGetRole(t *testing.T) {
	var (
		c   = context.TODO()
		bid = int64(0)
		uid = int64(0)
	)
	convey.Convey("GetRole", t, func(ctx convey.C) {
		httpMock("GET", d.c.Host.Manager+_getRole).Reply(200).JSON(`{"code":0,"data":[{"id":0}]}`)
		roles, err := d.GetRole(c, bid, uid)
		ctx.Convey("Then err should be nil.roles should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(roles, convey.ShouldNotBeNil)
		})
	})
}

func TestHttpGetUserRoles(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("GetUserRoles", t, func(ctx convey.C) {
		httpMock("GET", d.c.Host.Manager+_getRoles).Reply(200).JSON(`{"code":0,"data":[{"id":0}]}`)
		roles, err := d.GetUserRoles(c, 421)
		ctx.Convey("Then err should be nil.roles should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(roles, convey.ShouldNotBeNil)
		})
	})
}

func TestHttpGetUnames(t *testing.T) {
	var (
		c    = context.TODO()
		uids = []int64{421}
	)
	convey.Convey("GetUnames", t, func(ctx convey.C) {
		httpMock("GET", d.c.Host.Manager+_getUname).Reply(200).JSON(`{"code":0,"data":{"421":"丝瓜"}}`)
		unames, err := d.GetUnames(c, uids)
		ctx.Convey("Then err should be nil.unames should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(unames, convey.ShouldNotBeNil)
		})
	})
}

func TestHttpGetUIDs(t *testing.T) {
	convey.Convey("GetUIDs", t, func(ctx convey.C) {
		httpMock("GET", d.c.Host.Manager+_getUIDs).Reply(200).JSON(`{"code":0,"data":{"cxf":481}}`)
		uids, err := d.GetUIDs(cntx, "cxf")
		ctx.Convey("Then err should be nil.unames should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(uids, convey.ShouldNotBeNil)
		})
	})
}

func TestHttpGetUdepartment(t *testing.T) {
	convey.Convey("GetUdepartment", t, func(ctx convey.C) {
		httpMock("GET", d.c.Host.Manager+_getUdepartment).Reply(200).JSON(`{"code":0,"data":{"481":"CTO"}}`)
		depart, err := d.GetUdepartment(cntx, []int64{481})
		ctx.Convey("Then err should be nil.unames should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(depart, convey.ShouldNotBeNil)
		})
	})
}
