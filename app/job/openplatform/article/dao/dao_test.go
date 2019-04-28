package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/openplatform/article/conf"

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
	Convey("open reply", t, WithDao(func(d *Dao) {
		var (
			err error
			c   = context.TODO()
		)
		err = d.OpenReply(c, 88, 88)
		So(err, ShouldBeNil)
	}))
}
