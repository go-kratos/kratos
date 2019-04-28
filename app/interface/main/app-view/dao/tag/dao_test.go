package tag

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-view/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-view-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func ctx() context.Context {
	return context.Background()
}

func TestArcTags(t *testing.T) {
	Convey("get ArcTags all", t, func() {
		res, err := d.ArcTags(ctx(), 1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
