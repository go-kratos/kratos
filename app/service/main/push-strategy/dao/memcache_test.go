package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	biz  = int64(1)
	uuid = "uuid"
)

func TestDaopingMC(t *testing.T) {
	convey.Convey("pingMC", t, func(ctx convey.C) {
		err := d.pingMC(context.Background())
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func addUUIDCache() error {
	return d.AddUUIDCache(context.Background(), biz, uuid)
}

func TestDaoExistsUUIDCache(t *testing.T) {
	addUUIDCache()
	convey.Convey("ExistsUUIDCache", t, func(ctx convey.C) {
		exist, err := d.ExistsUUIDCache(context.Background(), biz, uuid)
		ctx.Convey("Then err should be nil.exist should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(exist, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddUUIDCache(t *testing.T) {
	convey.Convey("AddUUIDCache", t, func(ctx convey.C) {
		err := addUUIDCache()
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelUUIDCache(t *testing.T) {
	addUUIDCache()
	convey.Convey("DelUUIDCache", t, func(ctx convey.C) {
		err := d.DelUUIDCache(context.Background(), biz, uuid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
