package notice

import (
	"context"
	"go-common/app/service/live/xuserex/conf"
	"testing"
	"time"

	"flag"
	"github.com/smartystreets/goconvey/convey"
	"os"
)

var (
	d        *Dao
	UID      = int64(10000)
	targetID = int64(10000)
	date     = "20190101"
	term     = time.Now()
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "live.live.xuserex")
		flag.Set("conf_token", "4e6ace268a9ee6d04fad131ad551f61e")
		flag.Set("tree_id", "82470")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/test.toml")
		flag.Set("deploy.env", "uat")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestRoomNoticeNew(t *testing.T) {
	convey.Convey("New", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ctx.Convey("Then dao should not be nil.", func(ctx convey.C) {
				ctx.So(d, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomNoticekeyShouldNotice(t *testing.T) {
	convey.Convey("keyShouldNotice", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyShouldNotice(UID, targetID, date)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomNoticeIsNotice(t *testing.T) {
	convey.Convey("IsNotice", t, func(ctx convey.C) {
		convey.Convey("When everything goes positive", func(ctx convey.C) {
			c := context.TODO()
			p1, err := d.IsNotice(c, UID, targetID)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomNoticegetThreshold(t *testing.T) {
	convey.Convey("getThreshold", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			threshold, err := d.getThreshold()
			ctx.Convey("Then err should be nil.threshold should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(threshold, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomNoticegetShouldNotice(t *testing.T) {
	convey.Convey("getShouldNotice", t, func() {
		var (
			ctx = context.Background()
		)
		convey.Convey("When everything goes positive", func() {
			shouldNotice, err := d.getShouldNotice(ctx, UID, targetID, term)
			convey.Convey("Then err should be nil.shouldNotice should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(shouldNotice, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomNoticehbaseRowKey(t *testing.T) {
	convey.Convey("hbaseRowKey", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := hbaseRowKey(UID, targetID, date)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomNoticeRawMonthConsume(t *testing.T) {
	convey.Convey("RawMonthConsume", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawMonthConsume(c, UID, targetID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomNoticeGetTermBegin(t *testing.T) {
	convey.Convey("GetTermBegin", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.GetTermBegin()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomNoticeGetTermEnd(t *testing.T) {
	convey.Convey("GetTermEnd", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.GetTermEnd()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomNoticeIsValidTerm(t *testing.T) {
	convey.Convey("IsValidTerm", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.IsValidTerm(term)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestRoomNoticeisGuard(t *testing.T) {
	convey.Convey("isGuard", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			isGuard, err := d.isGuard(c, UID, targetID)
			ctx.Convey("Then err should be nil.isGuard should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(isGuard, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGetTaskFinish(t *testing.T) {
	convey.Convey("GetTaskFinish", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			term := d.GetTermBegin()
			isFinish, err := d.GetTaskFinish(c, term)
			ctx.Convey("Then err should be nil.isGuard should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(isFinish, convey.ShouldNotBeNil)
			})
		})
	})
}
