package bplus

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-resource/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-resource-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestUserCheck(t *testing.T) {
	Convey("get UserCheck all", t, func() {
		res, err := d.UserCheck(ctx(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
