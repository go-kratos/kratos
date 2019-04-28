package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"go-common/app/service/openplatform/anti-fraud/model"
)

func TestGetQusBankInfo(t *testing.T) {
	Convey("TestGetQusInfo", t, func() {
		_, err := svr.GetQusBankInfo(context.TODO(), 1527479929734)
		So(err, ShouldBeNil)
	})
}

func TestGetQusBanklist(t *testing.T) {
	Convey("TestGetQusInfo", t, func() {
		res, err := svr.GetQusBanklist(context.TODO(), 0, 10, "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestQusBankCheck(t *testing.T) {
	Convey("QusBankCheck", t, func() {
		data := &model.ArgCheckQus{
			Cnt:    3,
			QusIDs: []int64{1527065586769, 1527065707515, 2},
		}
		res, err := svr.QusBankCheck(context.TODO(), data)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
