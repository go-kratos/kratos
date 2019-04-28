package dao

import (
	"context"
	"go-common/app/interface/main/space/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoindexOrderHit(t *testing.T) {
	convey.Convey("indexOrderHit", t, func(ctx convey.C) {
		var (
			mid = int64(708)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := indexOrderHit(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoindexOrderKey(t *testing.T) {
	convey.Convey("indexOrderKey", t, func(ctx convey.C) {
		var (
			mid = int64(708)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := indexOrderKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIndexOrder(t *testing.T) {
	convey.Convey("IndexOrder", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(708)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			indexOrder, err := d.IndexOrder(c, mid)
			ctx.Convey("Then err should be nil.indexOrder should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(indexOrder, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIndexOrderModify(t *testing.T) {
	convey.Convey("IndexOrderModify", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(708)
			orderStr = `["1","2","8","7","3","4","5","6","21","22","23","24","25"]`
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.IndexOrderModify(c, mid, orderStr)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetIndexOrderCache(t *testing.T) {
	convey.Convey("SetIndexOrderCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(708)
			data = []*model.IndexOrder{{ID: 1, Name: "我的稿件"}}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetIndexOrderCache(c, mid, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoIndexOrderCache(t *testing.T) {
	convey.Convey("IndexOrderCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(708)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.IndexOrderCache(c, mid)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelIndexOrderCache(t *testing.T) {
	convey.Convey("DelIndexOrderCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(708)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelIndexOrderCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
