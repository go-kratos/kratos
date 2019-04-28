package dynamic

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/app-show/conf"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-show-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func Test_RegionDynamic(t *testing.T) {
	Convey("should get RegionDynamic", t, func() {
		_, _, err := d.RegionDynamic(context.Background(), 1, 2, 3)
		So(err, ShouldBeNil)
	})
}
