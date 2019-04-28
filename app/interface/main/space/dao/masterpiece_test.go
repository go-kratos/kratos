package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaomasterpieceHit(t *testing.T) {
	convey.Convey("masterpieceHit", t, func(ctx convey.C) {
		var (
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := masterpieceHit(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaomasterpieceKey(t *testing.T) {
	convey.Convey("masterpieceKey", t, func(ctx convey.C) {
		var (
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := masterpieceKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRawMasterpiece(t *testing.T) {
	convey.Convey("RawMasterpiece", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawMasterpiece(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddMasterpiece(t *testing.T) {
	convey.Convey("AddMasterpiece", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(2222)
			aid    = int64(2222)
			reason = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddMasterpiece(c, mid, aid, reason)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoEditMasterpiece(t *testing.T) {
	convey.Convey("EditMasterpiece", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(2222)
			aid    = int64(2222)
			preAid = int64(3333)
			reason = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.EditMasterpiece(c, mid, aid, preAid, reason)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelMasterpiece(t *testing.T) {
	convey.Convey("DelMasterpiece", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2222)
			aid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelMasterpiece(c, mid, aid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
