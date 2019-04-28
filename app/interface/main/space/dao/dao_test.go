package dao

import (
	"flag"
	"os"
	"strings"
	"testing"

	"go-common/app/interface/main/space/conf"

	gock "gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "dev" {
		flag.Set("app_id", "main.web-svr.space-interface")
		flag.Set("conf_token", "445a6be6947d52172ca5c4f3a55e3993")
		flag.Set("tree_id", "5242")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/space-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.httpR.SetTransport(gock.DefaultTransport)
	d.httpGame.SetTransport(gock.DefaultTransport)
	os.Exit(m.Run())
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}
