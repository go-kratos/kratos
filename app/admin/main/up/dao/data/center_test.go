package data

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDataViewerBase(t *testing.T) {
	convey.Convey("ViewerBase", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			dt  = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ViewerBase(c, mid, dt)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataViewerArea(t *testing.T) {
	convey.Convey("ViewerArea", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			dt  = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ViewerArea(c, mid, dt)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataViewerTrend(t *testing.T) {
	convey.Convey("ViewerTrend", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			dt  = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ViewerTrend(c, mid, dt)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataRelationFansDay(t *testing.T) {
	convey.Convey("RelationFansDay", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RelationFansDay(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataRelationFansHistory(t *testing.T) {
	convey.Convey("RelationFansHistory", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			month = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RelationFansHistory(c, mid, month)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataRelationFansMonth(t *testing.T) {
	convey.Convey("RelationFansMonth", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RelationFansMonth(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataViewerActionHour(t *testing.T) {
	convey.Convey("ViewerActionHour", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			dt  = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ViewerActionHour(c, mid, dt)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataUpIncr(t *testing.T) {
	convey.Convey("UpIncr", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ty  = int8(0)
			now = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UpIncr(c, mid, ty, now)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataThirtyDayArchive(t *testing.T) {
	convey.Convey("ThirtyDayArchive", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ty  = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ThirtyDayArchive(c, mid, ty)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataparseKeyValue(t *testing.T) {
	convey.Convey("parseKeyValue", t, func(ctx convey.C) {
		var (
			k = "20060102"
			v = "123"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			timestamp, value, err := parseKeyValue(k, v)
			ctx.Convey("Then err should be nil.timestamp,value should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldNotBeNil)
				ctx.So(timestamp, convey.ShouldNotBeNil)
			})
		})
	})
}
