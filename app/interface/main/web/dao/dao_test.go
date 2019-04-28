package dao

import (
	"flag"
	"os"
	"strings"
	"testing"

	"go-common/app/interface/main/web/conf"

	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "local" {
		flag.Set("app_id", "main.web-svr.web-interface")
		flag.Set("conf_token", "bc99caee81f27661bd5ff50d1c0d58d6")
		flag.Set("tree_id", "2685")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/web-interface-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.httpR.SetTransport(gock.DefaultTransport)
	d.httpBigData.SetTransport(gock.DefaultTransport)
	d.httpHelp.SetTransport(gock.DefaultTransport)
	os.Exit(m.Run())
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}
