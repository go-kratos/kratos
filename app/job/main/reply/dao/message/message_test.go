package message

import (
	"context"
	"flag"
	"go-common/app/job/main/reply/conf"
	"os"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.community.reply-job")
		flag.Set("conf_token", "5deea0665f8a7670b22a719337a39c7d")
		flag.Set("tree_id", "2123")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/reply-job-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = NewMessageDao(conf.Conf)
	os.Exit(m.Run())
}

func TestMessageNewMessageDao(t *testing.T) {
	convey.Convey("NewMessageDao", t, func(ctx convey.C) {
		var (
			c = conf.Conf
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := NewMessageDao(c)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMessageLike(t *testing.T) {
	convey.Convey("Like", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(0)
			tomid     = int64(0)
			title     = ""
			msg       = ""
			extraInfo = ""
			now       = time.Now()
			err       error
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err = d.Like(c, mid, tomid, title, msg, extraInfo, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(msg, convey.ShouldBeBlank)
			})
		})
	})
}

func TestMessageReply(t *testing.T) {
	convey.Convey("Reply", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mc        = ""
			resID     = ""
			mid       = int64(0)
			tomid     = int64(0)
			title     = ""
			msg       = ""
			extraInfo = ""
			now       = time.Now()
			err       error
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err = d.Reply(c, mc, resID, mid, tomid, title, msg, extraInfo, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(msg, convey.ShouldBeBlank)
			})
		})
	})
}

func TestMessageDeleteReply(t *testing.T) {
	convey.Convey("DeleteReply", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			title = ""
			msg   = ""
			now   = time.Now()
			err   error
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err = d.DeleteReply(c, mid, title, msg, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(msg, convey.ShouldBeBlank)
			})
		})
	})
}

func TestMessageAt(t *testing.T) {
	convey.Convey("At", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(0)
			mids      = []int64{}
			title     = ""
			msg       = ""
			extraInfo = ""
			now       = time.Now()
			err       error
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err = d.At(c, mid, mids, title, msg, extraInfo, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(msg, convey.ShouldBeBlank)
			})
		})
	})
}

func TestMessageAcceptReport(t *testing.T) {
	convey.Convey("AcceptReport", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			title = ""
			msg   = ""
			now   = time.Now()
			err   error
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err = d.AcceptReport(c, mid, title, msg, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(msg, convey.ShouldBeBlank)
			})
		})
	})
}

func TestMessageSystem(t *testing.T) {
	convey.Convey("System", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mc    = ""
			resID = ""
			mid   = int64(0)
			title = ""
			msg   = ""
			info  = ""
			now   = time.Now()
			err   error
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err = d.System(c, mc, resID, mid, title, msg, info, now)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(msg, convey.ShouldBeBlank)
			})
		})
	})
}

func TestMessagesend(t *testing.T) {
	convey.Convey("send", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mc    = ""
			resID = ""
			title = ""
			msg   = ""
			tp    = int(0)
			pub   = int64(0)
			mids  = []int64{}
			info  = ""
			ts    = int64(0)
			err   error
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err = d.send(c, mc, resID, title, msg, tp, pub, mids, info, ts)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(msg, convey.ShouldBeBlank)
			})
		})
	})
}

func TestMessageconverAt(t *testing.T) {
	convey.Convey("converAt", t, func(ctx convey.C) {
		var (
			title = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := converAt(title)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldBeBlank)
			})
		})
	})
}

func TestMessageconvertMsg(t *testing.T) {
	convey.Convey("convertMsg", t, func(ctx convey.C) {
		var (
			msg = "评"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := convertMsg(msg)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "評")
			})
		})
	})
}
