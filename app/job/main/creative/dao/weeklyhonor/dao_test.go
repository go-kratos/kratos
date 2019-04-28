package weeklyhonor

import (
	"flag"
	"os"
	"strings"
	"testing"

	"go-common/app/job/main/creative/conf"

	"github.com/golang/mock/gomock"
	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.creative-job")
		flag.Set("conf_token", "43943fda0bb311e8865c66d44b23cda7")
		flag.Set("tree_id", "16037")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/creative-job.toml")
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
	d.httpClient.SetTransport(gock.DefaultTransport)
	return r
}

func WithMock(t *testing.T, f func(controller *gomock.Controller)) func() {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	return func() {
		f(ctrl)
	}
}
