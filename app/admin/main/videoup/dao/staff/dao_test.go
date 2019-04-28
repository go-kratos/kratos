package staff

import (
	"flag"
	"go-common/app/admin/main/videoup/conf"
	"go-common/app/admin/main/videoup/model/archive"
	"testing"

	"context"
	"os"
	"strings"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
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

func TestStaffs(t *testing.T) {
	Convey("Staffs", t, WithDao(func(d *Dao) {
		httpMock("GET", d.staffURI).Reply(200).JSON(`{"code":0,"message":"0","ttl":1,"data":null}`)
		_, err := d.Staffs(context.TODO(), 1)
		So(err, ShouldBeNil)
	}))
	Convey("StaffApplyBatchSubmit", t, WithDao(func(d *Dao) {
		var ap *archive.StaffBatchParam
		httpMock("Do", d.staffURI).Reply(200).JSON(`{"code":0,"message":"0","ttl":1,"data":null}`)
		err := d.StaffApplyBatchSubmit(context.TODO(), ap)
		So(err, ShouldNotBeNil)
	}))
}
