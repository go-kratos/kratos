package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/space/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoprivacyHit(t *testing.T) {
	convey.Convey("privacyHit", t, func(ctx convey.C) {
		var (
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := privacyHit(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoprivacyKey(t *testing.T) {
	convey.Convey("privacyKey", t, func(ctx convey.C) {
		var (
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := privacyKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPrivacy(t *testing.T) {
	convey.Convey("Privacy", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.Privacy(c, mid)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPrivacyModify(t *testing.T) {
	convey.Convey("PrivacyModify", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(2222)
			field = "bangumi"
			value = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PrivacyModify(c, mid, field, value)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetPrivacyCache(t *testing.T) {
	convey.Convey("SetPrivacyCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(2222)
			data = model.DefaultPrivacy
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetPrivacyCache(c, mid, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoPrivacyCache(t *testing.T) {
	convey.Convey("PrivacyCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.PrivacyCache(c, mid)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelPrivacyCache(t *testing.T) {
	convey.Convey("DelPrivacyCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelPrivacyCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
