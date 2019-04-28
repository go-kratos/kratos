package activity

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup/conf"
	"go-common/app/job/main/videoup/model/archive"
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

func Test_AddVideo(t *testing.T) {
	var (
		c   = context.TODO()
		a   = new(archive.Archive)
		err error
	)
	Convey("AddVideo", t, WithDao(func(d *Dao) {
		err = d.AddVideo(c, a, 10086)
		So(err, ShouldBeNil)
	}))
}
