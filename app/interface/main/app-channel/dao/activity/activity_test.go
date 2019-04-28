package activity

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

func TestActivitys(t *testing.T) {
	Convey("get Activitys all", t, func() {
		res, err := d.Activitys(ctx(), []int64{1})
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
