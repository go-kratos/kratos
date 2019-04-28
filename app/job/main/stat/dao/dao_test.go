package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/stat/conf"
	"go-common/app/service/main/archive/api"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../cmd/stat-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func Test_SendQiyeWX(t *testing.T) {
	Convey("SendQiyeWX", t, func() {
		d.SendQiyeWX("日志报警:stat-job topic(test) 不消费了")
	})
}

func Test_Ping(t *testing.T) {
	Convey("Ping", t, func() {
		d.Ping(context.TODO())
	})
}

func Test_Stat(t *testing.T) {
	Convey("Stat", t, func() {
		st, err := d.Stat(context.TODO(), 10989901)
		So(err, ShouldBeNil)
		Printf("%+v", st)

	})
}

func Test_Update(t *testing.T) {
	Convey("Update", t, func() {
		_, err := d.Update(context.TODO(), &api.Stat{Aid: 10989901, Fav: 100, DisLike: 10, Like: 20})
		So(err, ShouldBeNil)
	})
}

func Test_MaxAID(t *testing.T) {
	Convey("MaxAID", t, func() {
		d.MaxAID(context.TODO())
	})
}
