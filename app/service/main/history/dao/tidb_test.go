package dao

import (
	"context"
	"testing"
	"time"

	pb "go-common/app/service/main/history/api/grpc"
	"go-common/app/service/main/history/model"

	"github.com/smartystreets/goconvey/convey"
)

var _history = &model.History{
	Mid:        1,
	BusinessID: 4,
	Business:   "pgc",
	Kid:        2,
	Aid:        3,
	Sid:        4,
	Epid:       5,
	Cid:        6,
	SubType:    7,
	Device:     8,
	Progress:   9,
	ViewAt:     10,
}
var hs = []*model.History{_history}

func TestDaoAddHistories(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("AddHistories", t, func(ctx convey.C) {
		err := d.AddHistories(c, hs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoQueryBusinesses(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("QueryBusinesses", t, func(ctx convey.C) {
		res, err := d.QueryBusinesses(c)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
		})
	})
}

func TestDaoUserHistories(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(1)
		viewAt = time.Now().Unix()
		ps     = int64(1)
	)
	convey.Convey("UserHistories all business", t, func(ctx convey.C) {
		var businesses []string
		res, err := d.UserHistories(c, businesses, mid, viewAt, ps)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
			res["pgc"][0].Ctime = 0
			res["pgc"][0].Mtime = 0
			ctx.So(res, convey.ShouldResemble, map[string][]*model.History{
				"pgc": {_history},
			})
		})
	})

	convey.Convey("UserHistories one business", t, func(ctx convey.C) {
		var businesses = []string{"pgc"}
		res, err := d.UserHistories(c, businesses, mid, viewAt, ps)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
			res["pgc"][0].Ctime = 0
			res["pgc"][0].Mtime = 0
			ctx.So(res, convey.ShouldResemble, map[string][]*model.History{
				"pgc": {_history},
			})
		})
	})
}

func TestDaoHistories(t *testing.T) {
	var (
		c        = context.Background()
		business = "pgc"
		mid      = int64(1)
		ids      = []int64{2}
	)
	convey.Convey("Histories", t, func(ctx convey.C) {
		res, err := d.Histories(c, business, mid, ids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeEmpty)
		})
	})
}

func TestDaoDeleteHistories(t *testing.T) {
	var (
		c = context.Background()
		h = &pb.DelHistoriesReq{Mid: 1, Records: []*pb.DelHistoriesReq_Record{
			{Business: "pgc", ID: 2},
		},
		}
	)
	convey.Convey("DeleteHistories", t, func(ctx convey.C) {
		d.AddHistories(c, hs)
		err := d.DeleteHistories(c, h)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoClearHistory(t *testing.T) {
	var (
		c          = context.Background()
		mid        = int64(1)
		businesses = []string{"pgc"}
	)
	convey.Convey("ClearHistory", t, func(ctx convey.C) {
		d.AddHistories(c, hs)
		err := d.ClearHistory(c, mid, businesses)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoClearAllHistory(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("ClearAllHistory", t, func(ctx convey.C) {
		err := d.ClearAllHistory(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpdateUserHide(t *testing.T) {
	var (
		c    = context.Background()
		mid  = int64(1)
		hide = true
	)
	convey.Convey("UpdateUserHide", t, func(ctx convey.C) {
		err := d.UpdateUserHide(c, mid, hide)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUserHide(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
	)
	convey.Convey("UserHide", t, func(ctx convey.C) {
		hide, err := d.UserHide(c, mid)
		ctx.Convey("Then err should be nil.hide should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(hide, convey.ShouldEqual, true)
		})
		hide, err = d.UserHide(c, 200)
		ctx.Convey("not found .Then err should be nil.hide should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(hide, convey.ShouldBeFalse)
		})
	})
}
