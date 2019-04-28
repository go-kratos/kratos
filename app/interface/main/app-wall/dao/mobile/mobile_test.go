package mobile

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-wall/conf"
	"go-common/app/interface/main/app-wall/model/mobile"

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

func TestInOrdersSync(t *testing.T) {
	Convey("unicom InOrdersSync", t, func() {
		res, err := d.InOrdersSync(ctx(), &mobile.MobileXML{})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestFlowSync(t *testing.T) {
	Convey("unicom FlowSync", t, func() {
		res, err := d.FlowSync(ctx(), &mobile.MobileXML{})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestOrdersUserFlow(t *testing.T) {
	Convey("unicom OrdersUserFlow", t, func() {
		res, err := d.OrdersUserFlow(ctx(), "", time.Now())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestMobileCache(t *testing.T) {
	Convey("unicom MobileCache", t, func() {
		res, err := d.MobileCache(ctx(), "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
