package resource

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-show/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-show-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestResBanner(t *testing.T) {
	Convey("ResBanner", t, func() {
		res, err := d.ResBanner(ctx(), 1, 1, 1, "", "", "", "", "", "", "", "", false)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
