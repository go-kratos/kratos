package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"testing"
)

func TestBasicParams(t *testing.T) {
	Convey("recharge", t, testWith(func() {
		tid := getTestTidForCall(t, int32(model.RECHARGETYPE))
		platform := "pc"
		uid := getTestRandUid()
		coinType := "gold"
		var num int64 = 100

		var area int64 = 10
		bizCode := "test/test"
		source := "test_source"
		bizSource := "test_biz_source"
		metaData := "test_meta"
		bp := getTestBasicParam(tid, area, bizCode, source, bizSource, metaData)

		v, err := s.Recharge(ctx, bp, uid, platform, getTestRechargeOrPayForm(uid, coinType, num, tid))
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)

		coinRecord, recordErr := s.dao.GetCoinStreamByTid(ctx, tid)
		So(recordErr, ShouldBeNil)
		So(coinRecord.MetaData, ShouldEqual, metaData)
		So(coinRecord.Area, ShouldEqual, area)
		So(coinRecord.BizCode, ShouldEqual, bizCode)
		So(coinRecord.Source, ShouldEqual, source)
		So(coinRecord.BizSource, ShouldEqual, bizSource)

	}))
}
