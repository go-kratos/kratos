package dao

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"go-common/app/service/openplatform/anti-fraud/conf"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(fmt.Errorf("conf.Init() error(%v)", err))
	}
	d = New(conf.Conf)
}

func TestGetQusBankListByIds(t *testing.T) {
	Convey("TestGetQusBankListByIds: ", t, func() {
		testIds := []int64{1527233672941, 2, 3}
		res, err := d.GetQusBankListByIds(context.TODO(), testIds)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestBankSearch(t *testing.T) {
	Convey("TestBankSearch: ", t, func() {
		res, err := d.BankSearch(context.TODO(), "name")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestGetQusBankCount(t *testing.T) {
	Convey("GetQusBankCount: ", t, func() {
		res, err := d.GetQusBankCount(context.TODO(), "wlt")
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func TestUpdateQsBankCnt(t *testing.T) {
	Convey("GetQusBankCount: ", t, func() {
		res, err := d.UpdateQsBankCnt(context.TODO(), 1527233672941)
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThanOrEqualTo, 0)
	})
}
