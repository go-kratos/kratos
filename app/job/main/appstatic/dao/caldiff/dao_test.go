package caldiff

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"go-common/app/job/main/appstatic/conf"

	"gopkg.in/h2non/gock.v1"
)

var (
	d   *Dao
	ctx = context.Background()
)

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.web-svr.appstatic-job")
		flag.Set("conf_token", "bb62e0021e8d0c7996ccf1fb8d267cee")
		flag.Set("tree_id", "40719")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		dir, _ := filepath.Abs("../cmd/appstatic-job-test.toml")
		flag.Set("conf", dir)
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.client.SetTransport(gock.DefaultTransport)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		f(d)
	}
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}
