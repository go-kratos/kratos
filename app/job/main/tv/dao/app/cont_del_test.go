package app

import (
	"context"
	"fmt"
	"testing"

	model "go-common/app/job/main/tv/model/pgc"

	"github.com/smartystreets/goconvey/convey"
)

func pickDelEpid() (epid int64, err error) {
	if err := d.DB.QueryRow(context.Background(), "select epid from tv_content where is_deleted = 1 limit 1").Scan(&epid); err != nil {
		fmt.Println("Pick EPid Err ", err)
	}
	return
}

func pickDelSid() (sid int64, err error) {
	if err := d.DB.QueryRow(context.Background(), "select id from tv_ep_season where is_deleted = 1 limit 1").Scan(&sid); err != nil {
		fmt.Println("Pick sid Err ", err)
	}
	return
}

func TestAppDelCont(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("DelCont", t, func(ctx convey.C) {
		res, err := d.DelCont(c)
		if err == nil && len(res) == 0 {
			fmt.Println("No deleted data, let me create one")
			epid, errPick := pickDelEpid()
			if errPick != nil {
				fmt.Println("pick err ", errPick)
				return
			}
			d.DB.Exec(c, "update tv_content set state = 1,audit_time = 0 where epid = ?", epid)
			res, err = d.DelCont(c)
		}
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("DelCont", t, func(ctx convey.C) {
		d.conf.Sync.LConf.SizeMsg = -1
		_, err := d.DelCont(c)
		ctx.So(err, convey.ShouldNotBeNil)
		fmt.Println(err)
	})
}

func TestAppSyncCont(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("SyncCont", t, func(ctx convey.C) {
		epid, errPick := pickDelEpid()
		if errPick != nil {
			fmt.Println("pick err ", errPick)
			return
		}
		nbRows, err := d.SyncCont(c, int(epid))
		ctx.Convey("Then err should be nil.nbRows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(nbRows, convey.ShouldNotBeNil)
		})
	})
}

func TestAppDelaySync(t *testing.T) {
	var c = context.Background()
	convey.Convey("DelaySync", t, func(ctx convey.C) {
		epid, errPick := pickDelEpid()
		if errPick != nil {
			fmt.Println("pick err ", errPick)
			return
		}
		conts := []*model.Content{{EPID: int(epid)}}
		nbRows, err := d.DelaySync(c, conts)
		ctx.Convey("Then err should be nil.nbRows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(nbRows, convey.ShouldNotBeNil)
		})
	})
}
