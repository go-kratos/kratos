package toview

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/history/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestToviewhashRowKey(t *testing.T) {
	convey.Convey("hashRowKey", t, func(ctx convey.C) {
		var (
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := hashRowKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestToviewAdd(t *testing.T) {
	convey.Convey("Add", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
			aid = int64(141787)
			now = time.Now().Unix()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Add(c, mid, aid, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestToviewAdds(t *testing.T) {
	convey.Convey("Adds", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(14771787)
			aids = []int64{147717, 147787}
			now  = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Adds(c, mid, aids, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestToviewAddMap(t *testing.T) {
	convey.Convey("AddMap", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(14771787)
			views map[int64]*model.ToView
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddMap(c, mid, views)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestToviewListInfo(t *testing.T) {
	convey.Convey("ListInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(14771787)
			aids = []int64{147717, 147787}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ListInfo(c, mid, aids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestToviewMapInfo(t *testing.T) {
	convey.Convey("MapInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(14771787)
			aids = []int64{147717, 147787}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MapInfo(c, mid, aids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestToviewDel(t *testing.T) {
	convey.Convey("Del", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(14771787)
			aids = []int64{147717, 147787}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Del(c, mid, aids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestToviewClear(t *testing.T) {
	convey.Convey("Clear", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Clear(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
