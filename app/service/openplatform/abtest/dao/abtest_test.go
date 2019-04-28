package dao

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"testing"

	"go-common/app/service/openplatform/abtest/conf"
	"go-common/library/cache/redis"

	. "github.com/smartystreets/goconvey/convey"
)

var testID int

func init() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(fmt.Errorf("conf.Init() error(%v)", err))
	}
	d = New(conf.Conf)
}

func TestActByGroup(t *testing.T) {
	Convey("TestActByGroup: ", t, func() {
		_, err := d.ActByGroup(context.TODO(), 1)
		So(err, ShouldBeNil)
	})
}

func TestListAb(t *testing.T) {
	Convey("TestListAb: ", t, func() {
		_, _, err := d.ListAb(context.TODO(), 1, 2, "0,1,2")
		So(err, ShouldBeNil)
	})
}

var testCase = "testxx"

func TestAddAb(t *testing.T) {
	Convey("TestAddAb: ", t, func() {
		var err error
		testID, err = redis.Int(d.AddAb(context.TODO(), "test", "test", `{"precision":100,"ratio":[80,20]}`, rand.Intn(10000000), 0, 1, "test"))
		So(err, ShouldBeNil)
		So(testID, ShouldNotEqual, 0)
	})
}

func TestUpAb(t *testing.T) {
	Convey("TestUpAb: ", t, func() {
		res, err := d.UpAb(context.TODO(), testID, "test", testCase, `{"precision":100,"ratio":[80,20]}`, 0, "update", 1, 0, 1)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 1)
	})
}

func TestAb(t *testing.T) {
	Convey("TestAb: ", t, func() {
		res, err := d.Ab(context.TODO(), testID)
		So(err, ShouldBeNil)
		So(res.Desc, ShouldEqual, testCase)
	})
}

func TestUpStatus(t *testing.T) {
	Convey("TestUpSta: ", t, func() {
		res, err := d.UpStatus(context.TODO(), testID, 1, "test", 1)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 1)
	})
	Convey("TestUpSta: ", t, func() {
		res, err := d.UpStatus(context.TODO(), testID, 0, "test", 1)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 1)
	})
}

func TestDelStatus(t *testing.T) {
	Convey("TestDelSta: ", t, func() {
		res, err := d.DelAb(context.TODO(), testID)
		So(err, ShouldBeNil)
		So(res, ShouldEqual, 1)
	})
}
