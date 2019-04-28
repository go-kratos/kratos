package data

import (
	"context"
	"flag"
	"go-common/app/admin/main/videoup/conf"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"os"
)

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}

func TestArchiveRelated(t *testing.T) {
	Convey("ArchiveRelated", t, WithDao(func(d *Dao) {
		httpMock("GET", d.relatedURI).Reply(200).JSON(`{"code":0,"data":[{"key":"123","value":"123"}]}`)
		_, err := d.ArchiveRelated(context.TODO(), []int64{10010, 10086})
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

func TestMonitorOids(t *testing.T) {
	Convey("MonitorOids", t, WithDao(func(d *Dao) {
		httpMock("GET", d.moniOidsURI).Reply(200).JSON(`{"code":0,"data":[{"oid":123,"time":123}]}`)
		_, err := d.MonitorOids(context.TODO(), 1)
		So(err, ShouldBeNil)
	}))
}
