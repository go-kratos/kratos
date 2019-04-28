package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpdateMask(t *testing.T) {
	var (
		c              = context.TODO()
		cid      int64 = 1352
		maskTime int64 = 60
		fps      int32 = 25
		list           = ""
	)
	Convey("test update mask", t, func() {
		err := svr.UpdateMask(c, cid, maskTime, fps, model.MaskPlatMbl, list)
		So(err, ShouldBeNil)
	})
}

func TestMaskList(t *testing.T) {
	var (
		c         = context.TODO()
		cid int64 = 1352
	)
	Convey("test mask list", t, func() {
		res, err := svr.MaskList(c, cid, model.MaskPlatMbl)
		t.Logf("==============%+v", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
