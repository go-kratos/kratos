package dao

import (
	"flag"
	"os"
	"strings"
	"testing"

	"go-common/app/job/main/member/conf"

	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.member-job")
		flag.Set("conf_token", "VEc5eqZNZHGQi6fsx7J6lJTqOGR9SnEO")
		flag.Set("tree_id", "2134")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/member-job-dev.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.client.SetTransport(gock.DefaultTransport)
	os.Exit(m.Run())
}
