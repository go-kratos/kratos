package region

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-channel/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-channel-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestAllList(t *testing.T) {
	Convey("get AllList all", t, func() {
		res, err := d.AllList(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestLimit(t *testing.T) {
	Convey("get Limit all", t, func() {
		res, err := d.Limit(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestConfig(t *testing.T) {
	Convey("get Config all", t, func() {
		res, err := d.Config(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
