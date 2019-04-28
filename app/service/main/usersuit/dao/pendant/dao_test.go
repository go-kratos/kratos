package pendant

import (
	"context"
	"flag"
	"os"
	"strings"
	"testing"

	"go-common/app/service/main/usersuit/conf"

	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
	c = context.Background()
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.usersuit-service")
		flag.Set("conf_token", "BVWgBtBvS2pkTBbmxAl0frX6KRA14d5P")
		flag.Set("tree_id", "6813")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/convey-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.client.SetTransport(gock.DefaultTransport)
	m.Run()
	os.Exit(0)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}
