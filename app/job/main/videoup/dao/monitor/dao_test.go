package monitor

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup/conf"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/videoup-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		f(d)
	}
}

func Test_History(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("History", t, WithDao(func(d *Dao) {
		err = d.Send(c, "test-msg")
		So(err, ShouldBeNil)
	}))
}
