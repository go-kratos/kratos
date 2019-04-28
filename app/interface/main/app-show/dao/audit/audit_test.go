package audit

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-show/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-show-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestAudits(t *testing.T) {
	Convey("Audits", t, func() {
		res, err := d.Audits(ctx())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestClose(t *testing.T) {
	Convey("Close", t, func() {
		d.Close()
	})
}
