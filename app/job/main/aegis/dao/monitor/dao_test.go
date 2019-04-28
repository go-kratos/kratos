package monitor

import (
	"flag"
	"go-common/app/job/main/aegis/conf"
	"gopkg.in/h2non/gock.v1"
	"os"
	"strings"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.aegis-job")
		flag.Set("conf_token", "aed3cc21ca345ffc284c6036da32352b")
		flag.Set("tree_id", "61819")
		flag.Set("conf_version", "1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/aegis-job.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.http.SetTransport(gock.DefaultTransport)
	return r
}
