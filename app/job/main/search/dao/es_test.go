package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/search/conf"
	"go-common/app/job/main/search/model"

	. "github.com/smartystreets/goconvey/convey"
)

func WithES(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../cmd/goconvey.toml")
		flag.Set("conf", dir)
		flag.Parse()
		conf.Init()
		d := New(conf.Conf)
		f(d)
	}
}

func Test_WithES(t *testing.T) {
	Convey("Test_WithES", t, WithES(func(d *Dao) {
		var (
			err error
			c   = context.TODO()
		)
		err = d.Ping(c)
		So(err, ShouldBeNil)
	}))
}

func Test_BulkDatabusData(t *testing.T) {
	Convey("Test_BulkDatabusData", t, WithES(func(d *Dao) {
		var (
			err   error
			c     = context.TODO()
			attrs = &model.Attrs{
				DataSQL: &model.AttrDataSQL{},
			}
		)
		attrs.ESName = "archive"
		attrs.DataSQL.DataIndexSuffix = ""
		d.BulkDatabusData(c, attrs, false)
		So(err, ShouldBeNil)
	}))
}

func Test_BulkDBData(t *testing.T) {
	Convey("Test_BulkDBData", t, WithES(func(d *Dao) {
		var (
			err   error
			c     = context.TODO()
			attrs = &model.Attrs{
				DataSQL: &model.AttrDataSQL{},
			}
		)
		attrs.ESName = "archive"
		attrs.DataSQL.DataIndexSuffix = ""
		d.BulkDBData(c, attrs, false)
		So(err, ShouldBeNil)
	}))
}

func Test_PingESCluster(t *testing.T) {
	Convey("Test_PingESCluster", t, WithES(func(d *Dao) {
		var (
			c   = context.TODO()
			err error
		)
		err = d.pingESCluster(c)
		So(err, ShouldBeNil)
	}))
}
