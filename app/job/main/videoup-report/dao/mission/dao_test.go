package mission

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/videoup-report/conf"
	"go-common/app/job/main/videoup-report/model/mission"

	. "github.com/smartystreets/goconvey/convey"
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
		Reset(func() {})
		f(d)
	}
}

func Test_Missions(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		mm  map[int]*mission.Mission
	)
	Convey("Missions", t, WithDao(func(d *Dao) {
		mm, err = d.Missions(c)
		So(err, ShouldBeNil)
		So(mm, ShouldNotBeNil)
		So(len(mm), ShouldBeGreaterThanOrEqualTo, 0)
	}))
}
