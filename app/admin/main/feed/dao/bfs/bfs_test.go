package bfs

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/feed/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestUpload(t *testing.T) {
	Convey("pull file bfs", t, func() {
		res, err := d.Upload(ctx(), "image/jpeg", nil)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
