package telecom

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-wall/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func ctx() context.Context {
	return context.Background()
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-wall-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestOrdersUserFlow(t *testing.T) {
	Convey("OrdersUserFlow", t, func() {
		res, err := d.OrdersUserFlow(ctx(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestOrdersUserByOrderID(t *testing.T) {
	Convey("OrdersUserByOrderID", t, func() {
		res, err := d.OrdersUserByOrderID(ctx(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	})
}
