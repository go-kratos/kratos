package resource

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-interface/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-interface-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestResBanner(t *testing.T) {
	Convey("Banner", t, func() {
		res, err := d.Banner(ctx(), "", "", "", "", "", "", "", 1, 1, 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
