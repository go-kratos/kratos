package offer

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"go-common/app/job/main/app-wall/conf"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-wall-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestActive(t *testing.T) {
	Convey("Active", t, func() {
		err := d.Active(ctx(), "", "", "", "", "")
		So(err, ShouldBeNil)
	})
}
