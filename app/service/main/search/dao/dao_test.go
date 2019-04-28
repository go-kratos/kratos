package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/service/main/search/conf"
	"go-common/app/service/main/search/model"

	. "github.com/smartystreets/goconvey/convey"
)

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../cmd/goconvey.toml")
		flag.Set("conf", dir)
		flag.Parse()
		conf.Init()
		d := New(conf.Conf)
		f(d)
	}
}

func Test_PgcMedia(t *testing.T) {
	Convey("open app", t, WithDao(func(d *Dao) {
		var (
			err error
			c   = context.TODO()
			p   *model.PgcMediaParams
		)
		_, err = d.PgcMedia(c, p)
		So(err, ShouldBeNil)
	}))
}

func Test_ReplyRecord(t *testing.T) {
	Convey("reply record", t, WithDao(func(d *Dao) {
		var (
			err error
			c   = context.TODO()
			p   *model.ReplyRecordParams
		)
		_, err = d.ReplyRecord(c, p)
		So(err, ShouldBeNil)
	}))
}

func Test_DmHistory(t *testing.T) {
	Convey("DmHistory", t, WithDao(func(d *Dao) {
		var (
			err error
			c   = context.TODO()
			p   *model.DmHistoryParams
		)
		_, err = d.DmHistory(c, p)
		So(err, ShouldBeNil)
	}))
}
