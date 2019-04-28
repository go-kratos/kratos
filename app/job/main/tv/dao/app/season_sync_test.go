package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppSnEmpty(t *testing.T) {
	var (
		c   = context.Background()
		sid = int64(12373)
	)
	convey.Convey("SnEmpty", t, func(ctx convey.C) {
		res, err := d.SnEmpty(c, sid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestAppModSeason(t *testing.T) {
	var c = context.Background()
	convey.Convey("ModSeason", t, func(ctx convey.C) {
		res, err := d.ModSeason(c)
		if len(res) == 0 && err == nil {
			sid, errPick := pickEpSid()
			if errPick != nil {
				return
			}
			fmt.Println(sid)
			d.DB.Exec(c, "update tv_ep_season set `check` = 2,audit_time=0 where id = ?", sid)
			res, err = d.ModSeason(c)
		}
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestAppAuditSeason(t *testing.T) {
	var (
		c   = context.Background()
		sid = int(12373)
	)
	convey.Convey("AuditSeason", t, func(ctx convey.C) {
		nbRows, err := d.AuditSeason(c, sid)
		ctx.Convey("Then err should be nil.nbRows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(nbRows, convey.ShouldNotBeNil)
		})
	})
}

func TestAppDelaySeason(t *testing.T) {
	var (
		c   = context.Background()
		sid = int64(12373)
	)
	convey.Convey("DelaySeason", t, func(ctx convey.C) {
		nbRows, err := d.DelaySeason(c, sid)
		ctx.Convey("Then err should be nil.nbRows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(nbRows, convey.ShouldNotBeNil)
		})
	})
}
