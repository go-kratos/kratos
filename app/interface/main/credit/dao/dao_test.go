package dao

import (
	"flag"
	"os"
	"strings"
	"testing"

	gock "gopkg.in/h2non/gock.v1"

	"go-common/app/interface/main/credit/conf"

	_ "github.com/go-sql-driver/mysql"
	// . "github.com/smartystreets/goconvey/convey"
)

var d *Dao

// func CleanCache() {
// 	c := context.TODO()
// 	pool := redis.NewPool(conf.Conf.Redis.Config)
// 	pool.Get(c).Do("FLUSHDB")
// }

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account-law.credit")
		flag.Set("conf_appid", "main.account-law.credit")
		flag.Set("conf_token", "aX4znxOXioonmhCtY5Piod6XLCHDPUKt")
		flag.Set("tree_id", "5659")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_env", "10")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/credit-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.client.SetTransport(gock.DefaultTransport)
	os.Exit(m.Run())
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

// func WithDao(f func(d *Dao)) func() {
// 	return func() {
// 		Reset(func() { CleanCache() })
// 		f(d)
// 	}
// }

// func WithMysql(f func(d *Dao)) func() {
// 	return func() {
// 		Reset(func() { CleanCache() })
// 		f(d)
// 	}
// }

// func WithCleanCache(f func()) func() {
// 	return func() {
// 		Reset(func() { CleanCache() })
// 	}
// }

// func httpMock(method, url string) *gock.Request {
// 	r := gock.New(url)
// 	r.Method = strings.ToUpper(method)
// 	return r
// }

// func ctx() context.Context {
// 	return context.Background()
// }
