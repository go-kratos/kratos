package tag

import (
	"context"
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/videoup/conf"
	"gopkg.in/h2non/gock.v1"
	"os"
	"strings"
	"testing"
)

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}
func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.client.SetTransport(gock.DefaultTransport)
	return r
}
func TestDao_AdminBind(t *testing.T) {
	Convey("AdminBind tag同步", t, WithDao(func(d *Dao) {
		httpMock("POST", d.uri).Reply(200).JSON(`{"code":0,"message":"0","ttl":1,"data":{}}`)
		err := d.AdminBind(context.TODO(), 1, 421, "haha", "日常", "")
		So(err, ShouldBeNil)
	}))
}

func TestDao_UpBind(t *testing.T) {
	Convey("UpBind tag同步", t, WithDao(func(d *Dao) {
		httpMock("POST", d.uri).Reply(200).JSON(`{"code":0,"message":"0","ttl":1,"data":{}}`)
		err := d.UpBind(context.TODO(), 1, 421, "haha", "日常", "")
		So(err, ShouldBeNil)
	}))
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
