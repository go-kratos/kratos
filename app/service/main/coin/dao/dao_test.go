package dao

import (
	"flag"
	"os"
	"strings"
	"testing"

	"go-common/app/service/main/coin/conf"

	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.coin-service")
		flag.Set("conf_appid", "coin-service")
		flag.Set("conf_token", "G1Rp73FzNnV0k2EdGdL7BZO5YL49n1Mt")
		flag.Set("tree_id", "2131")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/coin-service-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.httpClient.SetTransport(gock.DefaultTransport)
	m.Run()
	os.Exit(0)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}
