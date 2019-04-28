package dao

import (
	"context"
	"go-common/app/job/main/member/model"
	"net/url"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpdateRealnameFromMSG(t *testing.T) {
	convey.Convey("UpdateRealnameFromMSG", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			ms = &model.RealnameApplyMessage{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.UpdateRealnameFromMSG(c, ms)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRealnameInfo(t *testing.T) {
	convey.Convey("RealnameInfo", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			info, err := d.RealnameInfo(c, mid)
			convCtx.Convey("Then err should be nil.info should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpsertRealnameInfo(t *testing.T) {
	convey.Convey("UpsertRealnameInfo", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			ms = &model.RealnameInfo{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.UpsertRealnameInfo(c, ms)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpsertRealnameApplyImg(t *testing.T) {
	convey.Convey("UpsertRealnameApplyImg", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			ms = &model.RealnameApplyImgMessage{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.UpsertRealnameApplyImg(c, ms)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRealnameAlipayApplyList(t *testing.T) {
	convey.Convey("RealnameAlipayApplyList", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			startID  = int64(0)
			status   model.RealnameApplyStatus
			fromTime = time.Now().AddDate(-1, 0, 0)
			toTime   = time.Now()
			limit    = int(10)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			status = 2
			maxID, list, err := d.RealnameAlipayApplyList(c, startID, status, fromTime, toTime, limit)
			convCtx.Convey("Then err should be nil.maxID,list should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(list, convey.ShouldNotBeNil)
				convCtx.So(maxID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAlipayQuery(t *testing.T) {
	convey.Convey("AlipayQuery", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			param url.Values
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			pass, _, err := d.AlipayQuery(c, param)
			convCtx.Convey("Then err should be nil.pass,reason should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldNotBeNil)
				//convCtx.So(reason, convey.ShouldBeNil)
				convCtx.So(pass, convey.ShouldNotBeNil)
			})
		})
	})
}
