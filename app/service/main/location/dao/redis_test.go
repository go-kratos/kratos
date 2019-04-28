package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyZlimit(t *testing.T) {
	var (
		aid = int64(123456)
	)
	convey.Convey("keyZlimit", t, func(ctx convey.C) {
		key := keyZlimit(aid)
		ctx.Convey("key should not be equal to 123456", func(ctx convey.C) {
			ctx.So(key, convey.ShouldEqual, "zl_123456")
		})
	})
}

func TestDaoExistsAuth(t *testing.T) {
	var (
		c       = context.TODO()
		aid     = int64(123456)
		zoneids = map[int64]map[int64]int64{int64(123456): {int64(234567): int64(345678)}}
	)
	convey.Convey("ExistsAuth", t, WithDao(func(d *Dao) {
		d.AddAuth(c, zoneids)
		ok, err := d.ExistsAuth(c, aid)
		convey.Convey("Error should be nil, ok should be true", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldBeTrue)
		})
	}))
}

func TestDaoAuth(t *testing.T) {
	var (
		c       = context.TODO()
		aid     = int64(0)
		zoneid  = []int64{int64(0)}
		zoneids = map[int64]map[int64]int64{int64(123456): {int64(234567): int64(345678)}}
	)
	convey.Convey("Auth", t, WithDao(func(d *Dao) {
		d.AddAuth(c, zoneids)
		res, err := d.Auth(c, aid, zoneid)
		convey.Convey("Error should be nil, res should not be empty", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
		})
	}))
}

func TestDaoAddAuth(t *testing.T) {
	var (
		c       = context.TODO()
		zoneids = map[int64]map[int64]int64{int64(123456): {int64(234567): int64(345678)}}
	)
	convey.Convey("AddAuth", t, WithDao(func(d *Dao) {
		err := d.AddAuth(c, zoneids)
		convey.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	}))
}
