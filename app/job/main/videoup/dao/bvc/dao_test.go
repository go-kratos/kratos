package bvc

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup/conf"
)

func Test_VideoCapable(t *testing.T) {
	Convey("Tool", t, WithDao(func(d *Dao) {
		err := d.VideoCapable(context.TODO(), int64(2), []int64{21477219}, 0)
		So(err, ShouldBeNil)
	}))
}

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
		Reset(func() {})
		f(d)
	}
}
