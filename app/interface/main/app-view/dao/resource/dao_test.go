package resource

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

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
}

func Test_Paster(t *testing.T) {
	Convey("should get banner", t, func() {
		_, err := d.Paster(context.Background(), 1, 2, "", "", "")
		So(err, ShouldBeNil)
	})
}

func Test_PlayerIcon(t *testing.T) {
	Convey("should get player icon", t, func() {
		_, err := d.PlayerIcon(context.Background())
		So(err, ShouldBeNil)
	})
}
