package account

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

func TestRelations3(t *testing.T) {
	Convey("get Relations3 all", t, func() {
		res := d.Relations3(ctx(), []int64{1}, 1)
		So(res, ShouldNotBeEmpty)
	})
}

func TestCards3(t *testing.T) {
	Convey("get Cards3 all", t, func() {
		res, err := d.Cards3(ctx(), []int64{1})
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
