package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUids(t *testing.T) {
	convey.Convey("Uids", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			httpMock("GET", d.c.Host.Manager+_uidsURL).Reply(200).JSON(`{"code":0}`)
			_, err := d.Uids(c, []string{})
			ctx.Convey("Then nil should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUnames(t *testing.T) {
	convey.Convey("Unames", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			httpMock("GET", d.c.Host.Manager+_unamesURL).Reply(200).JSON(`{"code":0,"message":"0","data":{"10086":"cxf"}}`)
			_, err := d.Unames(c, []int64{})
			ctx.Convey("Then nil should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetUIDByName(t *testing.T) {
	convey.Convey("GetUIDByName", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			httpMock("GET", d.c.Host.Manager+_uidsURL).Reply(200).JSON(`{"code":0,"message":"0","data":{"cxf":10086}}`)
			uid, err := d.GetUIDByName(c, "cxf")
			ctx.Convey("Then nil should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(uid, convey.ShouldEqual, 10086)
			})
		})
	})
}

func TestDaoGetNameByUID(t *testing.T) {
	convey.Convey("GetNameByUID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			httpMock("GET", d.c.Host.Manager+_unamesURL).Reply(200).JSON(`{"code":0,"message":"0","data":{"10086":"cxf"}}`)
			_, err := d.GetNameByUID(c, []int64{10086})
			ctx.Convey("Then nil should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
