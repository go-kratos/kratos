package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpdateMask(t *testing.T) {
	var (
		c              = context.TODO()
		cid      int64 = 2386051
		maskTime int64 = 60
		fps      int32 = 25
		list           = "26777486_s0_0_1499,26777486_s1_1500_2999,26777486_s2_3000_4499,26777486_s3_4500_5999,26777486_s4_6000_7499,26777486_s5_7500_7897"
	)
	Convey("test update mask", t, func() {
		err := testDao.UpdateMask(c, cid, maskTime, fps, model.MaskPlatMbl, list)
		So(err, ShouldBeNil)
	})
}

func TestMaskList(t *testing.T) {
	var (
		c         = context.TODO()
		cid int64 = 1352
	)
	Convey("test mobile mask list", t, func() {
		res, err := testDao.MaskList(c, cid, model.MaskPlatMbl)
		t.Logf("==============%+v", res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
