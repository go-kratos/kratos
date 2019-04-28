package business

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/search/conf"
	"go-common/app/job/main/search/dao"

	. "github.com/smartystreets/goconvey/convey"
)

func WithBusinessArv(f func(d *Avr)) func() {
	return func() {
		dir, _ := filepath.Abs("../dao/cmd/goconvey.toml")
		flag.Set("conf", dir)
		flag.Parse()
		conf.Init()
		d := dao.New(conf.Conf)
		bsn := NewAvr(d, "archive_video", conf.Conf)
		f(bsn)
	}
}

func Test_AvrRecover(t *testing.T) {
	Convey("set recover", t, WithBusinessArv(func(d *Avr) {
		var (
			err error
			c   = context.TODO()
		)
		d.SetRecover(c, 1000, "", 0)
		So(err, ShouldBeNil)
	}))
}

func Test_AvrInitOffset(t *testing.T) {
	Convey("test close", t, WithBusinessArv(func(d *Avr) {
		d.InitOffset(context.TODO())
	}))
}

func WithBusinessDmDate(f func(d *DmDate)) func() {
	return func() {
		dir, _ := filepath.Abs("../dao/cmd/goconvey.toml")
		flag.Set("conf", dir)
		flag.Parse()
		conf.Init()
		d := dao.New(conf.Conf)
		bsn := NewDmDate(d, "dm_search")
		f(bsn)
	}
}

func Test_DmDateRecover(t *testing.T) {
	Convey("set recover", t, WithBusinessDmDate(func(d *DmDate) {
		var (
			err error
			c   = context.TODO()
		)
		d.SetRecover(c, 1000, "", 0)
		So(err, ShouldBeNil)
	}))
}

func Test_DmDateInitOffset(t *testing.T) {
	Convey("test close", t, WithBusinessDmDate(func(d *DmDate) {
		d.InitOffset(context.TODO())
	}))
}

func WithBusinessLog(f func(d *Log)) func() {
	return func() {
		dir, _ := filepath.Abs("../dao/cmd/goconvey.toml")
		flag.Set("conf", dir)
		flag.Parse()
		conf.Init()
		d := dao.New(conf.Conf)
		bsn := NewLog(d, "log_audit")
		f(bsn)
	}
}

func Test_LogRecover(t *testing.T) {
	Convey("set recover", t, WithBusinessLog(func(d *Log) {
		var (
			err error
			c   = context.TODO()
		)
		d.SetRecover(c, 1000, "", 0)
		So(err, ShouldBeNil)
	}))
}

func Test_LogInitOffset(t *testing.T) {
	Convey("test close", t, WithBusinessLog(func(d *Log) {
		d.InitOffset(context.TODO())
	}))
}

func Test_LogInitIndex(t *testing.T) {
	Convey("test init index", t, WithBusinessLog(func(d *Log) {
		d.InitIndex(context.TODO())
	}))
}

func Test_LogOffset(t *testing.T) {
	Convey("test offset", t, WithBusinessLog(func(d *Log) {
		d.Offset(context.TODO())
	}))
}

func Test_LogSetRecover(t *testing.T) {
	Convey("test set recover", t, WithBusinessLog(func(d *Log) {
		d.SetRecover(context.TODO(), 0, "", 0)
	}))
}

func Test_LogAllMessages(t *testing.T) {
	Convey("test set recover", t, WithBusinessLog(func(d *Log) {
		d.AllMessages(context.TODO())
	}))
}
