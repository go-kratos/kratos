package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"go-common/app/admin/main/tv/conf"

	"flag"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

var d *Dao

func init() {
	// dir, _ := filepath.Abs("../cmd/tv-admin-test.toml")
	// flag.Set("conf", dir)
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.web-svr.tv-admin")
		flag.Set("conf_token", "3d446a004187a6572d656bab1dbff1b0")
		flag.Set("tree_id", "15310")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
}

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
	d.httpSearch.SetTransport(gock.DefaultTransport)
	d.bfsClient.Transport = gock.DefaultTransport
	return r
}

func TestDao_MaxOrder(t *testing.T) {
	Convey("TestDao_MaxOrder", t, WithDao(func(d *Dao) {
		order := d.MaxOrder(context.Background())
		So(order, ShouldBeGreaterThan, 0)
		fmt.Println(order)
	}))
}

func TestDao_MangoRecom(t *testing.T) {
	Convey("TestDao_MangoRecom", t, WithDao(func(d *Dao) {
		err := d.MangoRecom(context.Background(), []int64{3, 4, 5})
		So(err, ShouldBeNil)
		res, err2 := d.GetMRecom(context.Background())
		So(err2, ShouldBeNil)
		data, _ := json.Marshal(res)
		fmt.Println(string(data))
	}))
}
