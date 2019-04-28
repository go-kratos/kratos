package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/point/conf"
	"go-common/app/admin/main/point/model"
	pointmol "go-common/app/service/main/point/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c = context.TODO()
)

func init() {
	dir, _ := filepath.Abs("../cmd/point-admin.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	fmt.Printf("%+v", conf.Conf)
	time.Sleep(time.Second)
}

func Test_PointHistory(t *testing.T) {
	Convey("Test_PointHistory", t, func() {
		arg := &model.ArgPointHistory{}
		res, err := s.PointHistory(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_PointCoinInfo(t *testing.T) {
	Convey("Test_PointCoinInfo", t, func() {
		res, err := s.PointCoinInfo(context.TODO(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_PointCoin(t *testing.T) {
	Convey("Test_PointCoin", t, func() {
		var (
			id  int64
			err error
			res *model.PointConf
		)
		pc := &model.PointConf{AppID: 1}
		id, err = s.PointCoinAdd(context.TODO(), pc)
		So(err, ShouldBeNil)
		So(id, ShouldBeGreaterThanOrEqualTo, 1)
		res, err = s.PointCoinInfo(context.TODO(), id)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.AppID, ShouldEqual, pc.AppID)
	})
}

func Test_PointConfList(t *testing.T) {
	Convey("Test_PointConfList", t, func() {
		res, err := s.PointConfList(context.TODO())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestPointAdd(t *testing.T) {
	Convey("TestPointAdd", t, func() {
		arg := new(pointmol.ArgPoint)
		arg.Mid = int64(2222)
		arg.Point = int64(1)
		arg.Remark = "系统发放"
		arg.Operator = "yubaihai"
		arg.ChangeType = model.PointSystem
		err := s.PointAdd(c, arg)
		So(err, ShouldBeNil)
	})
}
