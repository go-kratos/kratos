package app

import (
	"context"
	"testing"

	model "go-common/app/job/main/tv/model/pgc"
	"go-common/library/database/sql"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppRemoveCont(t *testing.T) {
	var (
		c    = context.Background()
		epid = int(12373)
	)
	convey.Convey("RemoveCont", t, func(ctx convey.C) {
		nbRows, err := d.RemoveCont(c, epid)
		ctx.Convey("Then err should be nil.nbRows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(nbRows, convey.ShouldNotBeNil)
		})
	})
}

func TestAppPickData(t *testing.T) {
	var (
		c   = context.Background()
		sid = int64(296)
	)
	convey.Convey("PickData", t, func(ctx convey.C) {
		_, err := d.PickData(c, sid)
		if err == sql.ErrNoRows {
			println(err)
		} else {
			ctx.So(err, convey.ShouldBeNil)
		}
	})
}

func TestAppSeason(t *testing.T) {
	var (
		c   = context.Background()
		sid = int(296)
	)
	convey.Convey("Season", t, func(ctx convey.C) {
		r, err := d.Season(c, sid)
		if err == sql.ErrNoRows {
			println(err)
		} else {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		}
	})
}

func TestAppEP(t *testing.T) {
	var (
		c    = context.Background()
		epid = int(12373)
	)
	convey.Convey("EP", t, func(ctx convey.C) {
		r, err := d.EP(c, epid)
		if err == sql.ErrNoRows {
			println(err)
		} else {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(r, convey.ShouldNotBeNil)
		}
	})
}

func TestAppCont(t *testing.T) {
	var (
		c    = context.Background()
		epid = int(12373)
	)
	convey.Convey("Cont", t, func(ctx convey.C) {
		res, err := d.Cont(c, epid)
		if err == sql.ErrNoRows {
			println(err)
		} else {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		}
	})
}

func TestAppWaitCall(t *testing.T) {
	var (
		c    = context.Background()
		epid = int(12373)
	)
	convey.Convey("WaitCall", t, func(ctx convey.C) {
		nbRows, err := d.WaitCall(c, epid)
		ctx.Convey("Then err should be nil.nbRows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(nbRows, convey.ShouldNotBeNil)
		})
	})
}

func TestAppDeleteEP(t *testing.T) {
	var (
		c   = context.Background()
		sid = int(296)
	)
	convey.Convey("DeleteEP", t, func(ctx convey.C) {
		nbRows, err := d.DeleteEP(c, sid)
		ctx.Convey("Then err should be nil.nbRows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(nbRows, convey.ShouldNotBeNil)
		})
	})
}

func TestAppDeleteCont(t *testing.T) {
	var (
		c   = context.Background()
		sid = int(296)
	)
	convey.Convey("DeleteCont", t, func(ctx convey.C) {
		nbRows, err := d.DeleteCont(c, sid)
		ctx.Convey("Then err should be nil.nbRows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(nbRows, convey.ShouldNotBeNil)
		})
	})
}

func TestAppRejectCont(t *testing.T) {
	var (
		c   = context.Background()
		sid = int(296)
	)
	convey.Convey("RejectCont", t, func(ctx convey.C) {
		nbRows, err := d.RejectCont(c, sid)
		ctx.Convey("Then err should be nil.nbRows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(nbRows, convey.ShouldNotBeNil)
		})
	})
}

func TestAppAuditingCont(t *testing.T) {
	var (
		c     = context.Background()
		conts []*model.Content
	)
	a := &model.Content{
		EPID: 12373,
	}
	conts = append(conts, a)
	convey.Convey("AuditingCont", t, func(ctx convey.C) {
		_, err := d.AuditingCont(c, conts)
		if err == sql.ErrNoRows {
			println(err)
		} else {
			ctx.So(err, convey.ShouldBeNil)
		}
	})
}

func TestAppReadySns(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("ReadySns", t, func(ctx convey.C) {
		_, err := d.ReadySns(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
