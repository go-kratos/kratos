package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDaoAiScore(t *testing.T) {
	convey.Convey("AiScore", t, func(ctx convey.C) {
		var (
			c       = context.TODO()
			content = "ai测试"
			stype   = "live_danmu" //reply
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.aiScoreURL).Reply(200).JSON(`{"code": 0}`)
			res, err := d.AiScore(c, content, stype)
			ctx.So(res, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoReplyDel(t *testing.T) {
	convey.Convey("ReplyDel", t, func(ctx convey.C) {
		var (
			c      = context.TODO()
			adid   = int64(1)
			oid    = int64(2)
			rpid   = int64(3)
			avType = int8(1)
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.mngReplyDelURL).Reply(200).JSON(`{"code": 0}`)
			err := d.ReplyDel(c, adid, oid, rpid, avType)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoReplyLabel(t *testing.T) {
	convey.Convey("ReplyLabel", t, func(ctx convey.C) {
		var (
			c      = context.TODO()
			adid   = int64(1)
			oid    = int64(2)
			rpid   = int64(3)
			avType = int8(4)
		)
		ctx.Convey("When everything right.", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.mngReplyLabelURL).Reply(200).JSON(`{"code": 0}`)
			err := d.ReplyLabel(c, adid, oid, rpid, avType)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
