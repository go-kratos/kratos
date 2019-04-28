package archive

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"go-common/library/cache/redis"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

var (
	errConnClosed = errors.New("redigo: connection closed")
)

func TestArchivekeyUpFavTpsPrefix(t *testing.T) {
	var (
		mid = int64(888952460)
	)
	convey.Convey("keyUpFavTpsPrefix", t, func(ctx convey.C) {
		p1 := keyUpFavTpsPrefix(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveFavTypes(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(888952460)
	)
	convey.Convey("FavTypes", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.redis), "Get", func(_ *redis.Pool, _ context.Context) redis.Conn {
			return redis.MockWith(errConnClosed)
		})
		defer connGuard.Unpatch()
		items, err := d.FavTypes(c, mid)
		ctx.Convey("Then err should be nil.items should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(items, convey.ShouldBeNil)
		})
	})
}

func TestArchivepingRedis(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("pingRedis", t, func(ctx convey.C) {
		err := d.pingRedis(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
