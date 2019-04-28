package search

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

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.httpClient.SetTransport(gock.DefaultTransport)
	return r
}
func TestOutTime(t *testing.T) {
	Convey("OutTime", t, WithDao(func(d *Dao) {
		httpMock("GET", d.URI).Reply(200).JSON(`{"code":0,"message":"0","ttl":1,"data":{"order":"ctime","sort":"desc","result":[],"debug":"","page":{"num":1,"size":10,"total":37}}}`)
		_, err := d.OutTime(context.TODO(), []int64{481, 6, 75, 248, 74, 246})
		So(err, ShouldBeNil)
	}))
}

func TestInQuitList(t *testing.T) {
	Convey("InQuitList", t, WithDao(func(d *Dao) {
		httpMock("GET", d.URI).Reply(200).JSON(`{"code":0,"message":"0","ttl":1,"data":{"order":"ctime","sort":"desc","result":[{"uid":0,"action":"0"}],"debug":"","page":{"num":1,"size":10,"total":37}}}`)
		_, err := d.InQuitList(context.TODO(), []int64{481}, "bt", "et")
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
