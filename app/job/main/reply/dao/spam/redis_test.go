package spam

import (
	"context"
	"flag"
	"go-common/app/job/main/reply/conf"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Cache
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
	d = NewCache(conf.Conf.Redis.Config)
	os.Exit(m.Run())
}

func TestSpamkeyRcntCnt(t *testing.T) {
	convey.Convey("keyRcntCnt", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.keyRcntCnt(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSpamkeyUpRcntCnt(t *testing.T) {
	convey.Convey("keyUpRcntCnt", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.keyUpRcntCnt(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSpamkeyDailyCnt(t *testing.T) {
	convey.Convey("keyDailyCnt", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.keyDailyCnt(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSpamkeyActRec(t *testing.T) {
	convey.Convey("keyActRec", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.keyActRec(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSpamkeySpamRpRec(t *testing.T) {
	convey.Convey("keySpamRpRec", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.keySpamRpRec(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSpamkeySpamRpDaily(t *testing.T) {
	convey.Convey("keySpamRpDaily", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.keySpamRpDaily(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSpamkeySpamActRec(t *testing.T) {
	convey.Convey("keySpamActRec", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.keySpamActRec(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSpamIncrReply(t *testing.T) {
	convey.Convey("IncrReply", t, func(ctx convey.C) {
		var (
			mid  = int64(0)
			isUp bool
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.IncrReply(context.Background(), mid, isUp)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSpamIncrAct(t *testing.T) {
	convey.Convey("IncrAct", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.IncrAct(context.Background(), mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSpamIncrDailyReply(t *testing.T) {
	convey.Convey("IncrDailyReply", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.IncrDailyReply(context.Background(), mid)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSpamTTLDailyReply(t *testing.T) {
	convey.Convey("TTLDailyReply", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ttl, err := d.TTLDailyReply(context.Background(), mid)
			ctx.Convey("Then err should be nil.ttl should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ttl, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSpamExpireDailyReply(t *testing.T) {
	convey.Convey("ExpireDailyReply", t, func(ctx convey.C) {
		var (
			mid = int64(0)
			exp = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.ExpireDailyReply(context.Background(), mid, exp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestSpamSetReplyRecSpam(t *testing.T) {
	convey.Convey("SetReplyRecSpam", t, func(ctx convey.C) {
		var (
			mid  = int64(0)
			code = int(0)
			exp  = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetReplyRecSpam(context.Background(), mid, code, exp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestSpamSetReplyDailySpam(t *testing.T) {
	convey.Convey("SetReplyDailySpam", t, func(ctx convey.C) {
		var (
			mid  = int64(0)
			code = int(0)
			exp  = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetReplyDailySpam(context.Background(), mid, code, exp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestSpamSetActionRecSpam(t *testing.T) {
	convey.Convey("SetActionRecSpam", t, func(ctx convey.C) {
		var (
			mid  = int64(0)
			code = int(0)
			exp  = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetActionRecSpam(context.Background(), mid, code, exp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
