package cms

import (
	"go-common/app/interface/main/tv/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCmsauthMsg(t *testing.T) {
	var (
		cond = ""
	)
	convey.Convey("authMsg", t, func(ctx convey.C) {
		msg := d.authMsg(cond)
		ctx.Convey("Then msg should not be nil.", func(ctx convey.C) {
			ctx.So(msg, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsSnErrMsg(t *testing.T) {
	var (
		season = &model.SnAuth{}
	)
	convey.Convey("SnErrMsg", t, func(ctx convey.C) {
		p1, p2 := d.SnErrMsg(season)
		ctx.Convey("Then p1,p2 should not be nil.", func(ctx convey.C) {
			ctx.So(p2, convey.ShouldNotBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsUgcErrMsg(t *testing.T) {
	var (
		deleted = int(0)
		result  = int(0)
		valid   = int(0)
	)
	convey.Convey("UgcErrMsg", t, func(ctx convey.C) {
		p1, p2 := d.UgcErrMsg(deleted, result, valid)
		ctx.Convey("Then p1,p2 should not be nil.", func(ctx convey.C) {
			ctx.So(p2, convey.ShouldNotBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsAuditingMsg(t *testing.T) {
	convey.Convey("AuditingMsg", t, func(ctx convey.C) {
		p1, p2 := d.AuditingMsg()
		ctx.Convey("Then p1,p2 should not be nil.", func(ctx convey.C) {
			ctx.So(p2, convey.ShouldNotBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsEpErrMsg(t *testing.T) {
	var (
		ep = &model.EpAuth{}
	)
	convey.Convey("EpErrMsg", t, func(ctx convey.C) {
		p1, p2 := d.EpErrMsg(ep)
		ctx.Convey("Then p1,p2 should not be nil.", func(ctx convey.C) {
			ctx.So(p2, convey.ShouldNotBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
