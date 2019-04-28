package archive

import (
	"context"
	"flag"
	"go-common/app/admin/main/videoup/conf"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
	"os"
	"strings"
)

var (
	d *Dao
)

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}

func TestWeightLog(t *testing.T) {
	Convey("WeightLog", t, WithDao(func(d *Dao) {
		d.WeightLog(context.TODO(), 2604)
	}))
}

func TestGetNextTask(t *testing.T) {
	Convey("GetNextTask", t, WithDao(func(d *Dao) {
		_, err := d.GetNextTask(context.TODO(), 102)
		So(err, ShouldBeNil)
	}))
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.clientR.SetTransport(gock.DefaultTransport)
	return r
}

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.videoup-admin")
		flag.Set("conf_token", "gRSfeavV7kJdY9875Gf29pbd2wrdKZ1a")
		flag.Set("tree_id", "2307")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/videoup-admin.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
