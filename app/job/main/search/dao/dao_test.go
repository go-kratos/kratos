package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/search/conf"

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

func Test_Reply(t *testing.T) {
	Convey("open app", t, WithDao(func(d *Dao) {
		var (
			err error
			c   = context.TODO()
		)
		err = d.Ping(c)
		So(err, ShouldBeNil)
	}))
}

func Test_SetRecover(t *testing.T) {
	Convey("set recover", t, WithDao(func(d *Dao) {
		var (
			err error
			c   = context.TODO()
		)
		d.SetRecover(c, "archive_video", 1000, "", 0)
		So(err, ShouldBeNil)
	}))
}

func Test_Close(t *testing.T) {
	Convey("test close", t, WithDao(func(d *Dao) {
		d.Close()
	}))
}

func Test_SendSMS(t *testing.T) {
	Convey("test send sms", t, WithDao(func(d *Dao) {
		var err error
		err = d.SendSMS("test sms")
		So(err, ShouldBeNil)
	}))
}
