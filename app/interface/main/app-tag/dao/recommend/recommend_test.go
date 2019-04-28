package recommend

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

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-tag-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestFeedDynamic(t *testing.T) {
	Convey("FeedDynamic", t, func() {
		res, _, _, _, err := d.FeedDynamic(ctx(), true, 1, 1, 1, 0, 0, time.Now())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
