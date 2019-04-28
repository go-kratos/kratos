package manager

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup-report/conf"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/videoup-report-job.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		f(d)
	}
}

func Test_User(t *testing.T) {
	var (
		c = context.TODO()
	)

	Convey("User", t, WithDao(func(d *Dao) {
		um, err := d.User(c, 421)
		So(err, ShouldBeNil)
		So(um, ShouldNotBeNil)
	}))
}
