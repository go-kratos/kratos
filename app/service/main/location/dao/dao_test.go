package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/location/conf"
	"go-common/library/cache/redis"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.location-service")
		flag.Set("conf_token", "bdc7976c3e1dbf8adeb1cdb7b1e823af")
		flag.Set("tree_id", "3170")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/location-example.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func CleanCache() {
	pool := redis.NewPool(conf.Conf.Redis.Zlimit.Config)
	pool.Get(context.TODO()).Do("FLUSHDB")
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		convey.Reset(func() { CleanCache() })
		f(d)
	}
}
