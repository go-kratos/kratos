package dao

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/service/main/push/conf"
	"go-common/app/service/main/push/model"
	"go-common/library/cache/redis"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.web-svr.push-service")
		flag.Set("conf_token", "668329d872842a0079691e868e0fa12d")
		flag.Set("tree_id", "35083")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		dir, _ := filepath.Abs("../cmd/push-service-test.toml")
		flag.Set("conf", dir)
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(d)
	}
}

func CleanCache() {
	c := context.TODO()
	redisPool := redis.NewPool(conf.Conf.Redis.Config)
	redisPool.Get(c).Do("FLUSHDB")
}

func Test_MiInvalidTokens(t *testing.T) {
	Convey("fetch mi invalid tokens", t, WithDao(func(d *Dao) {
		// 用的时候打开，消息消费完了就没了
		// err := d.DelInvalidMiReports(context.TODO())
		// So(err, ShouldBeNil)
	}))
}

func Test_buildAPNS(t *testing.T) {
	Convey("build apns", t, func() {
		info := &model.PushInfo{
			TaskID:   model.TempTaskID(),
			APPID:    1,
			Title:    model.DefaultMessageTitle,
			Summary:  "bilibili",
			LinkType: 8,
		}
		item := &model.PushItem{
			Platform: 2,
			Token:    "sdfsdfewfsadfsdfsdf",
			Mid:      888,
		}
		apns := buildAPNS(info, item)
		t.Logf("apns(%+v)", apns)
	})
}
func Test_RefreshAuth(t *testing.T) {
	Convey("ping redis", t, WithDao(func(d *Dao) {
		err := d.Ping(context.Background())
		So(err, ShouldBeNil)
	}))
}
