package dao

import (
	"context"
	"go-common/app/service/video/stream-mng/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCreateUpStreamDispatch(t *testing.T) {
	convey.Convey("CreateUpStreamDispatch", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			info = &model.UpStreamInfo{
				RoomID:   11891462,
				CDN:      1,
				PlatForm: "ios",
				Country:  "中国",
				City:     "上海",
				ISP:      "电信",
				IP:       "12.12.12.12",
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.CreateUpStreamDispatch(c, info)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetSummaryUpStreamRtmp(t *testing.T) {
	convey.Convey("GetSummaryUpStreamRtmp", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = int64(1546593007)
			end   = int64(1546830611)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			infos, err := d.GetSummaryUpStreamRtmp(c, start, end)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetSummaryUpStreamISP(t *testing.T) {
	convey.Convey("GetSummaryUpStreamISP", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = int64(1546593007)
			end   = int64(1546830611)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			infos, err := d.GetSummaryUpStreamISP(c, start, end)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetSummaryUpStreamCountry(t *testing.T) {
	convey.Convey("GetSummaryUpStreamCountry", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = int64(1546593007)
			end   = int64(1546830611)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			infos, err := d.GetSummaryUpStreamCountry(c, start, end)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetSummaryUpStreamPlatform(t *testing.T) {
	convey.Convey("GetSummaryUpStreamPlatform", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = int64(1546593007)
			end   = int64(1546830611)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			infos, err := d.GetSummaryUpStreamPlatform(c, start, end)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetSummaryUpStreamCity(t *testing.T) {
	convey.Convey("GetSummaryUpStreamCity", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			start = int64(1546593007)
			end   = int64(1546830611)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			infos, err := d.GetSummaryUpStreamCity(c, start, end)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}
