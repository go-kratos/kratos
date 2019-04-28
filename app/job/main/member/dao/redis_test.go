package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoexpAddedKey(t *testing.T) {
	convey.Convey("expAddedKey", t, func(convCtx convey.C) {
		var (
			tp  = ""
			mid = int64(0)
			day = int(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := expAddedKey(tp, mid, day)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoleader(t *testing.T) {
	convey.Convey("leader", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			key, value := leader()
			convCtx.Convey("Then key,value should not be nil.", func(convCtx convey.C) {
				convCtx.So(value, convey.ShouldNotBeNil)
				convCtx.So(key, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetExpAdded(t *testing.T) {
	convey.Convey("SetExpAdded", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			day = int(0)
			tp  = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			b, err := d.SetExpAdded(c, mid, day, tp)
			convCtx.Convey("Then err should be nil.b should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(b, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpAdded(t *testing.T) {
	convey.Convey("ExpAdded", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			day = int(0)
			tp  = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			b, err := d.ExpAdded(c, mid, day, tp)
			convCtx.Convey("Then err should be nil.b should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(b, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLeaderEleciton(t *testing.T) {
	convey.Convey("LeaderEleciton", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			elected := d.LeaderEleciton(c)
			convCtx.Convey("Then elected should not be nil.", func(convCtx convey.C) {
				convCtx.So(elected, convey.ShouldNotBeNil)
			})
		})
	})
}
