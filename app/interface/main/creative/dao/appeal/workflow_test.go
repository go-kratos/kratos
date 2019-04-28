package appeal

import (
	"context"
	"go-common/app/interface/main/creative/model/appeal"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppealAppealList(t *testing.T) {
	var (
		c        = context.TODO()
		mid      = int64(0)
		business = int(0)
		ip       = ""
	)
	convey.Convey("AppealList", t, func(ctx convey.C) {
		as, err := d.AppealList(c, mid, business, ip)
		ctx.Convey("Then err should be nil.as should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(as, convey.ShouldNotBeNil)
		})
	})
}

func TestAppealAppealDetail(t *testing.T) {
	var (
		c        = context.TODO()
		mid      = int64(0)
		id       = int64(0)
		business = int(0)
		ip       = ""
	)
	convey.Convey("AppealDetail", t, func(ctx convey.C) {
		a, err := d.AppealDetail(c, mid, id, business, ip)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(a, convey.ShouldNotBeNil)
		})
	})
}

func TestAppealAddAppeal(t *testing.T) {
	var (
		c           = context.TODO()
		tid         = int64(0)
		aid         = int64(0)
		mid         = int64(0)
		business    = int64(0)
		qq          = ""
		phone       = ""
		email       = ""
		desc        = ""
		attachments = ""
		ip          = ""
		ap          = &appeal.BusinessAppeal{}
	)
	convey.Convey("AddAppeal", t, func(ctx convey.C) {
		apID, err := d.AddAppeal(c, tid, aid, mid, business, qq, phone, email, desc, attachments, ip, ap)
		ctx.Convey("Then err should be nil.apID should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(apID, convey.ShouldNotBeNil)
		})
	})
}

func TestAppealAddReply(t *testing.T) {
	var (
		c           = context.TODO()
		cid         = int64(0)
		event       = int64(0)
		content     = ""
		attachments = ""
		ip          = ""
	)
	convey.Convey("AddReply", t, func(ctx convey.C) {
		err := d.AddReply(c, cid, event, content, attachments, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAppealAppealExtra(t *testing.T) {
	var (
		c        = context.TODO()
		mid      = int64(0)
		cid      = int64(0)
		business = int64(0)
		val      = int64(0)
		key      = ""
		ip       = ""
	)
	convey.Convey("AppealExtra", t, func(ctx convey.C) {
		err := d.AppealExtra(c, mid, cid, business, val, key, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAppealAppealState(t *testing.T) {
	var (
		c        = context.TODO()
		mid      = int64(0)
		id       = int64(0)
		business = int64(0)
		state    = int64(0)
		ip       = ""
	)
	convey.Convey("AppealState", t, func(ctx convey.C) {
		err := d.AppealState(c, mid, id, business, state, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAppealAppealStarInfo(t *testing.T) {
	var (
		c        = context.TODO()
		mid      = int64(0)
		cid      = int64(0)
		business = int(0)
		ip       = ""
	)
	convey.Convey("AppealStarInfo", t, func(ctx convey.C) {
		star, etime, err := d.AppealStarInfo(c, mid, cid, business, ip)
		ctx.Convey("Then err should be nil.star,etime should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(etime, convey.ShouldNotBeNil)
			ctx.So(star, convey.ShouldNotBeNil)
		})
	})
}
