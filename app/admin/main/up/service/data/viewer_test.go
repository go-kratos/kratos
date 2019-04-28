package data

import (
	"context"
	"go-common/app/admin/main/up/model/datamodel"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDataViewerBase(t *testing.T) {
	convey.Convey("ViewerBase", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.ViewerBase(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataViewerArea(t *testing.T) {
	convey.Convey("ViewerArea", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.ViewerArea(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataCacheTrend(t *testing.T) {
	convey.Convey("CacheTrend", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.CacheTrend(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataGetTags(t *testing.T) {
	convey.Convey("GetTags", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result := s.GetTags(c, ids)
			ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataviewerTrend(t *testing.T) {
	convey.Convey("viewerTrend", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			dt  = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.viewerTrend(c, mid, dt)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataGetUpViewInfo(t *testing.T) {
	convey.Convey("GetUpViewInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &datamodel.GetUpViewInfoArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.GetUpViewInfo(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataGetViewData(t *testing.T) {
	convey.Convey("GetViewData", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.GetViewData(c, mid)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataGetFansSummary(t *testing.T) {
	convey.Convey("GetFansSummary", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &datamodel.GetFansSummaryArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.GetFansSummary(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataGetRelationFansDay(t *testing.T) {
	convey.Convey("GetRelationFansDay", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &datamodel.GetRelationFansHistoryArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.GetRelationFansDay(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
