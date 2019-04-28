package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/job/main/upload/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../cmd/bfs-upload-job.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func TestPing(t *testing.T) {
	Convey("Ping", t, func() {
		err := d.Ping(context.TODO())
		So(err, ShouldBeNil)
	})
}

func TestClose(t *testing.T) {
	Convey("Ping", t, func() {
		err := d.Close
		So(err, ShouldBeNil)
	})
}
