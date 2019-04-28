package archive

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

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
}

func Test_Archives(t *testing.T) {
	Convey("should get Archives", t, func() {
		_, err := d.Archives(context.Background(), []int64{1, 2, 3, 4})
		So(err, ShouldBeNil)
	})
}
