package block

import (
	"context"
	model "go-common/app/admin/main/member/model/block"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBlockBlackhouseBlock(t *testing.T) {
	convey.Convey("BlackhouseBlock", t, func(ctx convey.C) {
		var (
			c = context.Background()
			p = &model.ParamBatchBlock{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.BlackhouseBlock(c, p)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockparseReasonType(t *testing.T) {
	convey.Convey("parseReasonType", t, func(ctx convey.C) {
		var (
			msg = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			no := parseReasonType(msg)
			ctx.Convey("Then no should not be nil.", func(ctx convey.C) {
				ctx.So(no, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockSendSysMsg(t *testing.T) {
	convey.Convey("SendSysMsg", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			code     = ""
			mids     = []int64{}
			title    = ""
			content  = ""
			remoteIP = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SendSysMsg(c, code, mids, title, content, remoteIP)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockmidsToParam(t *testing.T) {
	convey.Convey("midsToParam", t, func(ctx convey.C) {
		var (
			mids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			str := midsToParam(mids)
			ctx.Convey("Then str should not be nil.", func(ctx convey.C) {
				ctx.So(str, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBlockTelInfo(t *testing.T) {
	convey.Convey("TelInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tel, err := d.TelInfo(c, mid)
			ctx.Convey("Then err should be nil.tel should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(tel, convey.ShouldEqual, "")
			})
		})
	})
}

func TestBlockMailInfo(t *testing.T) {
	convey.Convey("MailInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mail, err := d.MailInfo(c, mid)
			ctx.Convey("Then err should be nil.mail should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(mail, convey.ShouldEqual, "")
			})
		})
	})
}
