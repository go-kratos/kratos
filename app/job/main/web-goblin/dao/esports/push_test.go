package esports

import (
	"testing"

	mdlesp "go-common/app/job/main/web-goblin/model/esports"

	"github.com/smartystreets/goconvey/convey"
)

func TestEsportsNoticeUser(t *testing.T) {
	var (
		mids    = []int64{11, 22, 33}
		body    = ""
		contest = &mdlesp.Contest{}
	)
	convey.Convey("NoticeUser", t, func(ctx convey.C) {
		err := d.NoticeUser(mids, body, contest)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestEsportssignature(t *testing.T) {
	var (
		params map[string]string
		secret = ""
	)
	convey.Convey("signature", t, func(ctx convey.C) {
		p1 := d.signature(params, secret)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestEsportsgetUUID(t *testing.T) {
	var (
		mids    = "11,22,33"
		contest = &mdlesp.Contest{}
	)
	convey.Convey("getUUID", t, func(ctx convey.C) {
		p1 := d.getUUID(mids, contest)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
