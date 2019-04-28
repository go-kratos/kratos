package up

import (
	"context"
	"go-common/app/service/main/up/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpAddCacheUp(t *testing.T) {
	var (
		c   = context.TODO()
		id  = int64(0)
		val = &model.Up{}
	)
	convey.Convey("AddCacheUp", t, func(ctx convey.C) {
		err := d.AddCacheUp(c, id, val)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUpCacheUp(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("CacheUp", t, func(ctx convey.C) {
		res, err := d.CacheUp(c, id)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestUpDelCacheUp(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("DelCacheUp", t, func(ctx convey.C) {
		err := d.DelCacheUp(c, id)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUpAddCacheUpSwitch(t *testing.T) {
	var (
		c   = context.TODO()
		id  = int64(0)
		val = &model.UpSwitch{}
	)
	convey.Convey("AddCacheUpSwitch", t, func(ctx convey.C) {
		err := d.AddCacheUpSwitch(c, id, val)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUpCacheUpSwitch(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("CacheUpSwitch", t, func(ctx convey.C) {
		res, err := d.CacheUpSwitch(c, id)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestUpDelCacheUpSwitch(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("DelCacheUpSwitch", t, func(ctx convey.C) {
		err := d.DelCacheUpSwitch(c, id)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUpAddCacheUpInfoActive(t *testing.T) {
	var (
		c   = context.TODO()
		id  = int64(0)
		val = &model.UpInfoActiveReply{}
	)
	convey.Convey("AddCacheUpInfoActive", t, func(ctx convey.C) {
		err := d.AddCacheUpInfoActive(c, id, val)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUpCacheUpInfoActive(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("CacheUpInfoActive", t, func(ctx convey.C) {
		res, err := d.CacheUpInfoActive(c, id)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestUpDelCacheUpInfoActive(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(0)
	)
	convey.Convey("DelCacheUpInfoActive", t, func(ctx convey.C) {
		err := d.DelCacheUpInfoActive(c, id)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
