package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/share/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRedisKey(t *testing.T) {
	var (
		oid = int64(0)
		tp  = int(0)
	)
	convey.Convey("redisKey", t, func(ctx convey.C) {
		p1 := redisKey(oid, tp)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoRedisValue(t *testing.T) {
	convey.Convey("redisValue", t, func(ctx convey.C) {
		p := &model.ShareParams{
			OID: 22,
			MID: 33,
			TP:  2,
			IP:  "",
		}
		p1 := redisValue(p)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoShareKey(t *testing.T) {
	var (
		oid = int64(0)
		tp  = int(0)
	)
	convey.Convey("shareKey", t, func(ctx convey.C) {
		p1 := shareKey(oid, tp)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddShareMember(t *testing.T) {
	convey.Convey("AddShareMember", t, func(ctx convey.C) {
		p := &model.ShareParams{
			OID: int64(1),
			MID: int64(1),
			TP:  int(3),
		}
		ok, err := d.AddShareMember(context.Background(), p)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSetShareCache(t *testing.T) {
	var (
		c      = context.TODO()
		oid    = int64(0)
		tp     = int(0)
		shared = int64(0)
	)
	convey.Convey("SetShareCache", t, func(ctx convey.C) {
		err := d.SetShareCache(c, oid, tp, shared)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoShareCache(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(0)
		tp  = int(0)
	)
	convey.Convey("ShareCache", t, func(ctx convey.C) {
		shared, err := d.ShareCache(c, oid, tp)
		ctx.Convey("Then err should be nil.shared should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(shared, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSharesCache(t *testing.T) {
	var (
		c    = context.TODO()
		oids = []int64{}
		tp   = int(0)
	)
	convey.Convey("SharesCache", t, func(ctx convey.C) {
		shares, err := d.SharesCache(c, oids, tp)
		ctx.Convey("Then err should be nil.shares should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(shares, convey.ShouldNotBeNil)
		})
	})
}
