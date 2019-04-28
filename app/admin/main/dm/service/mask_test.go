package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMaskState(t *testing.T) {
	convey.Convey("mask state", t, func() {
		open, mobile, web, err := svr.MaskState(context.TODO(), 1, 1352)
		convey.So(err, convey.ShouldBeNil)
		t.Logf("===%d,%d,%d", open, mobile, web)
	})
}

func TestUpdateMaskState(t *testing.T) {
	convey.Convey("open mask", t, func() {
		err := svr.UpdateMaskState(context.TODO(), 1, 1352, 1, 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestGenerateMask(t *testing.T) {
	convey.Convey("generate mask", t, func() {
		err := svr.GenerateMask(context.TODO(), 1, 1352, 1)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestMaskUps(t *testing.T) {
	convey.Convey("test mask ups", t, func() {
		res, err := svr.MaskUps(context.Background(), 1, 50)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
		t.Log(res, res.Result, res.Page)
		for _, v := range res.Result {
			t.Log(v)
		}
	})
}

func TestMaskUpOpen(t *testing.T) {
	convey.Convey("test mask up open", t, func() {
		err := svr.MaskUpOpen(context.Background(), []int64{1111, 142341123}, 1, "")
		convey.So(err, convey.ShouldBeNil)
	})
}
