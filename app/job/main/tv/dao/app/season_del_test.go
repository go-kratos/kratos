package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppDelSeason(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("DelSeason", t, func(ctx convey.C) {
		res, err := d.DelSeason(c)
		if err == nil && len(res) == 0 {
			sid, errPick := pickDelSid()
			if errPick != nil {
				return
			}
			fmt.Println(sid)
			d.DB.Exec(c, "update tv_ep_season set `check` = 2,audit_time = 0 where id = ?", sid)
			res, err = d.DelSeason(c)
		}
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestAppRejectSeason(t *testing.T) {
	var (
		c   = context.Background()
		sid = int(0)
	)
	convey.Convey("RejectSeason", t, func(ctx convey.C) {
		nbRows, err := d.RejectSeason(c, sid)
		ctx.Convey("Then err should be nil.nbRows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(nbRows, convey.ShouldNotBeNil)
		})
	})
}
