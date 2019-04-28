package tag

import (
	"context"
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"path/filepath"
	"testing"
	"time"

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
	rand.Seed(time.Now().UnixNano())
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}

func TestDao_UpBind(t *testing.T) {
	Convey("upbind", t, WithDao(func(d *Dao) {
		err := d.UpBind(context.TODO(), 176, 1, "haha", "日常", "")
		So(err, ShouldBeNil)
	}))
}

func TestDao_AdminBind(t *testing.T) {
	Convey("adminbind", t, WithDao(func(d *Dao) {
		err := d.AdminBind(context.TODO(), 176, 2, "haha", "日常", "")
		So(err, ShouldBeNil)
	}))
}
