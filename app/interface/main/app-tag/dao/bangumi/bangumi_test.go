package bangumi

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-tag/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-tag-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func ctx() context.Context {
	return context.Background()
}

func TestSeasonid(t *testing.T) {
	Convey("get Seasonid all", t, func() {
		res, err := d.Seasonid([]int64{1}, time.Now())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
