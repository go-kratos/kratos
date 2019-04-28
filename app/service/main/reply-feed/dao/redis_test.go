package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyReplyZSet(t *testing.T) {
	convey.Convey("keyReplyZSet", t, func(ctx convey.C) {
		var (
			name = ""
			oid  = int64(0)
			tp   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyReplyZSet(name, oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPingRedis(t *testing.T) {
	convey.Convey("PingRedis", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PingRedis(context.Background())
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoExpireReplyZSetRds(t *testing.T) {
	convey.Convey("ExpireReplyZSetRds", t, func(ctx convey.C) {
		var (
			name = ""
			oid  = int64(0)
			tp   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.ExpireReplyZSetRds(context.Background(), name, oid, tp)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetReplyZSetRds(t *testing.T) {
	convey.Convey("ReplyZSetRds", t, func(ctx convey.C) {
		var (
			name  = ""
			oid   = int64(0)
			tp    = int(0)
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIDs, err := d.ReplyZSetRds(context.Background(), name, oid, tp, start, end)
			ctx.Convey("Then err should be nil.rpIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIDs, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCountReplyZSetRds(t *testing.T) {
	convey.Convey("CountReplyZSetRds", t, func(ctx convey.C) {
		var (
			name = ""
			oid  = int64(-1)
			tp   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.CountReplyZSetRds(context.Background(), name, oid, tp)
			ctx.Convey("Then err should be nil.rpIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldEqual, 0)
			})
		})
	})
}
