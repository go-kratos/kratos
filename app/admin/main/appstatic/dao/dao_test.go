package dao

import (
	"flag"
	"os"
	"strings"

	"go-common/app/admin/main/appstatic/conf"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

var d *Dao

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.web-svr.appstatic-admin")
		flag.Set("conf_token", "26eec0d4137fac469c03e5ae147ed101")
		flag.Set("tree_id", "22976")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/appstatic-admin-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.client.SetTransport(gock.DefaultTransport)
	return r
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}
