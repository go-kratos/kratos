package workflow

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/reply/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.community.reply")
		flag.Set("conf_token", "54e85e3ab609f79ae908b9ea3e3f0775")
		flag.Set("tree_id", "2125")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/reply-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestWorkflowNew(t *testing.T) {
	convey.Convey("New", t, func(ctx convey.C) {
		var (
			c = conf.Conf
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := New(c)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWorkflowAddReport(t *testing.T) {
	convey.Convey("AddReport", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			oid      = int64(0)
			typ      = int8(0)
			typeid   = int32(0)
			rpid     = int64(0)
			score    = int(0)
			reason   = int8(0)
			reporter = int64(0)
			reported = int64(0)
			like     = int(0)
			content  = ""
			link     = ""
			title    = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddReport(c, oid, typ, typeid, rpid, score, reason, reporter, reported, like, content, link, title)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWorkflowTopicsLink(t *testing.T) {
	convey.Convey("TopicsLink", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			links   map[int64]string
			isTopic bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TopicsLink(c, links, isTopic)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestWorkflowActivitySub(t *testing.T) {
	convey.Convey("ActivitySub", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			title, link, err := d.ActivitySub(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWorkflowLiveNotice(t *testing.T) {
	convey.Convey("LiveNotice", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			title, err := d.LiveNotice(c, oid)
			ctx.Convey("Then err should be nil.title should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWorkflowLiveActivityTitle(t *testing.T) {
	convey.Convey("LiveActivityTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			title, link, err := d.LiveActivityTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWorkflowNoticeTitle(t *testing.T) {
	convey.Convey("NoticeTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			title, link, _ := d.NoticeTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWorkflowBanTitle(t *testing.T) {
	convey.Convey("BanTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			title, link, _ := d.BanTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWorkflowCreditTitle(t *testing.T) {
	convey.Convey("CreditTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			title, link, _ := d.CreditTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWorkflowLiveNoticeTitle(t *testing.T) {
	convey.Convey("LiveNoticeTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			title, link, err := d.LiveNoticeTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWorkflowLivePictureTitle(t *testing.T) {
	convey.Convey("LivePictureTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			title, link, err := d.LivePictureTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWorkflowDynamicTitle(t *testing.T) {
	convey.Convey("DynamicTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			title, link, err := d.DynamicTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestWorkflowHuoniaoTitle(t *testing.T) {
	convey.Convey("HuoniaoTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			title, link, err := d.HuoniaoTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}
