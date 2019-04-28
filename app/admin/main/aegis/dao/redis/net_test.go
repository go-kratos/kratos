package redis

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	kvs = map[string]string{
		"test:1": "test111-11",
		"test:2": "test222-22",
	}
	ks = []string{"test:1", "test:2"}
)

func TestRedisSetMulti(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("SetMulti", t, func(ctx convey.C) {
		err := d.SetMulti(c, kvs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestRedisMGet(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("MGet", t, func(ctx convey.C) {
		dest, err := d.MGet(c, ks...)
		ctx.Convey("Then err should be nil.dest should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(dest, convey.ShouldNotBeNil)
		})
	})
}

func TestRedisDelMulti(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("DelMulti", t, func(ctx convey.C) {
		err := d.DelMulti(c, ks...)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
