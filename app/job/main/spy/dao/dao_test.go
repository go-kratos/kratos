package dao

import (
	"context"
	"flag"
	"os"
	"strings"
	"testing"

	"go-common/app/job/main/spy/conf"

	. "github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
	c = context.Background()
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account-law.spy-job")
		flag.Set("conf_appid", "main.account-law.spy-job")
		flag.Set("conf_token", "1404609bda3db4cd20ad9823bfd20a83")
		flag.Set("tree_id", "2856")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_env", "10")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/spy-job-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.httpClient.SetTransport(gock.DefaultTransport)
	os.Exit(m.Run())
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func TestPing(t *testing.T) {
	Convey("TestPing", t, func() {
		err := d.Ping(c)
		So(err, ShouldBeNil)
	})

}
