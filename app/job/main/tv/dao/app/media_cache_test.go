package app

import (
	"context"
	"fmt"
	model "go-common/app/job/main/tv/model/pgc"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAppSnCMSCacheKey(t *testing.T) {
	var (
		sid = int(0)
	)
	convey.Convey("SnCMSCacheKey", t, func(ctx convey.C) {
		p1 := d.SnCMSCacheKey(sid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAppEpCMSCacheKey(t *testing.T) {
	var (
		epid = int(0)
	)
	convey.Convey("EpCMSCacheKey", t, func(ctx convey.C) {
		p1 := d.EpCMSCacheKey(epid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestAppSetSnCMSCache(t *testing.T) {
	var (
		c = context.Background()
		s = &model.SeasonCMS{}
	)
	convey.Convey("SetSnCMSCache", t, func(ctx convey.C) {
		err := d.SetSnCMSCache(c, s)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAppSetEpCMSCache(t *testing.T) {
	var (
		c = context.Background()
		s = &model.EpCMS{}
	)
	convey.Convey("SetEpCMSCache", t, func(ctx convey.C) {
		err := d.SetEpCMSCache(c, s)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func pickEpSid() (sid int64, err error) {
	if err = d.DB.QueryRow(context.Background(), "select season_id from tv_content where is_deleted = 0 "+
		"and valid = 1 and state = 3 limit 1").Scan(&sid); err != nil {
		fmt.Println(err)
	}
	return
}

func TestAppNewestOrder(t *testing.T) {
	var c = context.Background()
	convey.Convey("NewestOrder", t, func(ctx convey.C) {
		sid, errPick := pickEpSid()
		if errPick != nil {
			return
		}
		epid, newestOrder, err := d.NewestOrder(c, sid)
		ctx.Convey("Then err should be nil.epid,newestOrder should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(newestOrder, convey.ShouldNotBeNil)
			ctx.So(epid, convey.ShouldNotBeNil)
		})
	})
}

func TestAppAllEP(t *testing.T) {
	var (
		c        = context.Background()
		strategy = int(0)
	)
	convey.Convey("AllEP", t, func(ctx convey.C) {
		sid, errPick := pickEpSid()
		if errPick != nil {
			return
		}
		eps, err := d.AllEP(c, int(sid), strategy)
		ctx.Convey("Then err should be nil.eps should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(eps, convey.ShouldNotBeNil)
		})
	})
}
