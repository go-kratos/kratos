package location

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-channel/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-channel-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestAuthPIDs(t *testing.T) {
	Convey("get AuthPIDs", t, func() {
		_, err := d.AuthPIDs(ctx(), "417,1521", "127.0.0.0")
		So(err, ShouldBeNil)
	})
}
