package ads

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/resource/conf"
	"go-common/library/cache/redis"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.resource-service")
		flag.Set("conf_token", "y79sErNhxggjvULS0O8Czas9PaxHBF5o")
		flag.Set("tree_id", "3232")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/resource-service-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

// func CleanCache() {
// 	pool := redis.NewPool(conf.Conf.Redis.Ads.Config)
// 	pool.Get(context.TODO()).Do("FLUSHDB")
// }

func WithDao(f func(d *Dao)) func() {
	return func() {
		convey.Reset(func() {
			pool := redis.NewPool(conf.Conf.Redis.Ads.Config)
			pool.Get(context.TODO()).Do("FLUSHDB")
		})
		f(d)
	}
}

func TestDaoPing(t *testing.T) {
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(context.TODO())
		ctx.Convey("Err should be nil", func() {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
