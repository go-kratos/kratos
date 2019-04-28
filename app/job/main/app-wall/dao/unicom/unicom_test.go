package unicom

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/app-wall/conf"
	"go-common/app/job/main/app-wall/model/unicom"

	. "github.com/smartystreets/goconvey/convey"
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

func TestUpUserIntegral(t *testing.T) {
	Convey("UpUserIntegral", t, func() {
		res, err := d.UpUserIntegral(ctx(), &unicom.UserBind{})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestOrdersUserFlow(t *testing.T) {
	Convey("OrdersUserFlow", t, func() {
		res, err := d.OrdersUserFlow(ctx(), "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestBindAll(t *testing.T) {
	Convey("BindAll", t, func() {
		res, err := d.BindAll(ctx(), 1, 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestUserBind(t *testing.T) {
	Convey("UserBind", t, func() {
		res, err := d.UserBind(ctx(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestIPSync(t *testing.T) {
	Convey("IPSync", t, func() {
		res, err := d.IPSync(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestInUserPackLog(t *testing.T) {
	Convey("InUserPackLog", t, func() {
		res, err := d.InUserPackLog(ctx(), nil)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestInUserIntegralLog(t *testing.T) {
	Convey("InUserIntegralLog", t, func() {
		res, err := d.InUserIntegralLog(ctx(), nil)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestUserBindCache(t *testing.T) {
	Convey("UserBindCache", t, func() {
		res, err := d.UserBindCache(ctx(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestAddUserBindCache(t *testing.T) {
	Convey("AddUserBindCache", t, func() {
		err := d.AddUserBindCache(ctx(), 1, nil)
		So(err, ShouldBeNil)
	})
}

func TestDeleteUserPackReceiveCache(t *testing.T) {
	Convey("DeleteUserPackReceiveCache", t, func() {
		err := d.DeleteUserPackReceiveCache(ctx(), 1)
		So(err, ShouldBeNil)
	})
}
