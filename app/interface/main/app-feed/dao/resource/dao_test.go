package resource

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/app-feed/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-feed-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func Test_Banner(t *testing.T) {
	Convey("should get banner", t, func() {
		_, _, err := d.Banner(context.Background(), 1, 2, 3, "", "", "", "", "", "", true, "", "", "")
		So(err, ShouldBeNil)
	})
}
