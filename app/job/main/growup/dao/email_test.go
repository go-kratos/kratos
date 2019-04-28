package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTotalIncome(t *testing.T) {
	convey.Convey("TotalIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
			from  = int64(0)
			limit = int64(100)
		)
		fmt.Println("date:", date)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_income(av_id,mid,income,date) VALUES (1,2,3,'2018-06-24') ON DUPLICATE KEY UPDATE av_id=VALUES(av_id), date=VALUES(date)")
			infos, err := d.TotalIncome(c, date, from, limit)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAV(t *testing.T) {
	convey.Convey("GetAV", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_income(av_id, mid, income, date) VALUES (1,2,3,'2018-06-24') ON DUPLICATE KEY UPDATE av_id=VALUES(av_id), date=VALUES(date)")
			infos, err := d.GetAV(c, date)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetTag(t *testing.T) {
	convey.Convey("GetTag", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = time.Date(2018, 6, 24, 0, 0, 0, 0, time.Local)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO tag_info(tag,category_id,is_common,date) VALUES(2,3,1, '2018-06-24')")
			infos, err := d.GetTag(c, date)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetMID(t *testing.T) {
	convey.Convey("GetMID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			TagID = int64(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO tag_up_info(mid) VALUES(10) ON DUPLICATE KEY UPDATE mid=VALUES(mid)")
			infos, err := d.GetMID(c, TagID)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTagToAV(t *testing.T) {
	convey.Convey("TagToAV", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			category = int(2)
			date     = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.TagToAV(c, category, date)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMIDToAV(t *testing.T) {
	convey.Convey("MIDToAV", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(100)
			category = int(10)
			date     = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.MIDToAV(c, mid, category, date)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}
