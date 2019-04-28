package coin

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

func TestAddCoins(t *testing.T) {
	Convey("get AddCoins all", t, func() {
		err := d.AddCoins(ctx(), 1, 1, 1, 1, 1, 1, 1, 0)
		So(err, ShouldBeNil)
	})
}
