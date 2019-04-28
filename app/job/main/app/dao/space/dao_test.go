package space

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/app/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(5 * time.Second)
}

func Test_UpArchives(t *testing.T) {
	Convey("UpArchives", t, func() {
		_, err := d.UpArchives(context.TODO(), 1, 1, 10, "")
		So(err, ShouldBeNil)
	})
}

func Test_UpArcCnt(t *testing.T) {
	Convey("UpArcCnt", t, func() {
		_, err := d.UpArcCnt(context.TODO(), 1, "")
		So(err, ShouldBeNil)
	})
}

func Test_UpArticles(t *testing.T) {
	Convey("UpArticles", t, func() {
		_, _, err := d.UpArticles(context.TODO(), 1, 1, 10)
		So(err, ShouldBeNil)
	})
}
