package view

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

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
}

func Test_PingMc(t *testing.T) {
	Convey("PingMc", t, func() {
		err := d.PingMc(context.TODO())
		So(err, ShouldBeNil)
	})
}
