package service

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	"testing"
	"time"
)

func getTestRecordCoinStreamItem(coinType string, coinNum int64, tid string, orgCoinNum int64, serviceType string) *model.RecordCoinStreamItem {

	return &model.RecordCoinStreamItem{
		CoinType:      coinType,
		CoinNum:       coinNum,
		ExtendTid:     getTestExtendTid(),
		Timestamp:     time.Now().Unix(),
		TransactionId: tid,
		Type:          serviceType,
		OrgCoinNum:    orgCoinNum,
	}
}

func getDetalWithSnapShot(u *TestUser) *model.DetailWithSnapShot {
	var detail2 *model.DetailWithSnapShot
	tx, _ := s.dao.BeginTx(ctx)

	detail2, _ = s.dao.WalletForUpdate(tx, u.uid)
	tx.Rollback()
	return detail2
}

func queryRecordCoinStream(t *testing.T, uid int64, platform string, tid string, streamList []*model.RecordCoinStreamItem) error {

	bp := getTestDefaultBasicParam(tid)
	bp.Version = 23
	bp.Source = "kkk"

	steamListBytes, _ := json.Marshal(streamList)

	form := &model.RecordCoinStreamForm{Uid: uid, Data: string(steamListBytes)}

	_, err := s.RecordCoinStream(ctx, bp, uid, platform, form)
	return err
}

func TestService_RecordCoinStream(t *testing.T) {
	Convey("normal", t, testWithTestUser(func(u *TestUser) {
		tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
		platform := "pc"
		coinType := "gold"
		var num int64 = 100

		d1 := getDetalWithSnapShot(u)

		streamList := make([]*model.RecordCoinStreamItem, 0)

		streamList = append(streamList, getTestRecordCoinStreamItem(coinType, 0-num, tid, u.coin, "pay"))
		streamList = append(streamList, getTestRecordCoinStreamItem(coinType, num, tid, u.coin-num, "recharge"))
		streamList = append(streamList, getTestRecordCoinStreamItem("silver", 1024, tid, u.coin, "recharge"))
		streamList = append(streamList, getTestRecordCoinStreamItem("silver", -1023, tid, u.coin, "pay"))

		streamList[0].Reserved1 = 1
		streamList[1].Reserved1 = 2

		err := queryRecordCoinStream(t, u.uid, platform, tid, streamList)
		So(err, ShouldBeNil)

		record := assertNormalCoinStream(t, platform, "gold", tid, 3, -1*num, model.PAYTYPE, model.STREAM_OP_RESULT_SUB_SUCC, TX_STATUS_SUCC)
		t.Logf("record:%+v", record)
		So(record.OrgCoinNum, ShouldEqual, u.coin)
		So(record.Version, ShouldEqual, 23)

		assertNormalCoinStream(t, platform, "gold", tid, 2, num, model.RECHARGETYPE, model.STREAM_OP_RESULT_ADD_SUCC, TX_STATUS_SUCC)

		record, err = s.dao.GetCoinStreamByTidAndOffset(ctx, tid, 2)
		So(err, ShouldBeNil)
		So(record.OpResult, ShouldEqual, 2)
		So(record.Platform, ShouldEqual, model.GetPlatformNo(platform))
		So(record.CoinType, ShouldEqual, model.GetCoinTypeNumber(model.GetSysCoinType(coinType, platform)))
		So(record.Reserved1, ShouldEqual, 2)

		record, err = s.dao.GetCoinStreamByTidAndOffset(ctx, tid, 3)
		So(err, ShouldBeNil)
		So(record.OpResult, ShouldEqual, 1)
		So(record.Platform, ShouldEqual, model.GetPlatformNo(platform))
		So(record.CoinType, ShouldEqual, model.GetCoinTypeNumber(model.GetSysCoinType(coinType, platform)))
		So(record.Reserved1, ShouldEqual, 1)

		d2 := getDetalWithSnapShot(u)
		So(num, ShouldEqual, d2.GoldRechargeCnt-d1.GoldRechargeCnt)
		So(1023, ShouldEqual, d2.SilverPayCnt-d1.SilverPayCnt)

	}))

	Convey("params", t, testWith(func() {
		Convey("data err", testWithRandUid(func(uid int64) {
			tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
			platform := "pc"
			data := "unknown"

			form := &model.RecordCoinStreamForm{Uid: uid, Data: data}
			bp := getTestDefaultBasicParam(tid)

			_, err := s.RecordCoinStream(ctx, bp, uid, platform, form)
			So(err, ShouldEqual, ecode.RequestErr)

		}))

		Convey("data empty", testWithRandUid(func(uid int64) {
			tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
			platform := "pc"

			streamList := make([]*model.RecordCoinStreamItem, 0)
			err := queryRecordCoinStream(t, uid, platform, tid, streamList)
			So(err, ShouldEqual, ecode.RequestErr)
		}))
		Convey("data too much", testWithRandUid(func(uid int64) {
			tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
			platform := "pc"
			coinType := "gold"
			var num int64 = 100

			streamList := make([]*model.RecordCoinStreamItem, 0)

			for i := 0; i < 30; i++ {
				streamList = append(streamList, getTestRecordCoinStreamItem(coinType, 0-num, tid, 1000, "pay"))
				streamList = append(streamList, getTestRecordCoinStreamItem(coinType, num, tid, 900, "recharge"))
			}

			err := queryRecordCoinStream(t, uid, platform, tid, streamList)
			So(err, ShouldEqual, ecode.RequestErr)

		}))

	}))

	Convey("logic", t, testWith(func() {
		Convey("wrong orgCoinNum", testWithRandUid(func(uid int64) {
			tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
			platform := "pc"
			coinType := "gold"
			var num int64 = 100

			streamList := make([]*model.RecordCoinStreamItem, 0)
			streamList = append(streamList, getTestRecordCoinStreamItem(coinType, 0-num, tid, -100, "pay"))
			err := queryRecordCoinStream(t, uid, platform, tid, streamList)
			So(err, ShouldEqual, ecode.RequestErr)
		}))

		Convey("wrong validType", testWithRandUid(func(uid int64) {
			tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
			platform := "pc"
			coinType := "gold"
			var num int64 = 100

			streamList := make([]*model.RecordCoinStreamItem, 0)
			streamList = append(streamList, getTestRecordCoinStreamItem(coinType, 0-num, tid, 100, "exchange"))
			err := queryRecordCoinStream(t, uid, platform, tid, streamList)
			So(err, ShouldEqual, ecode.RequestErr)
		}))

		Convey("wrong validCoinType", testWithRandUid(func(uid int64) {
			tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
			platform := "pc"
			coinType := "unknown"
			var num int64 = 100

			streamList := make([]*model.RecordCoinStreamItem, 0)
			streamList = append(streamList, getTestRecordCoinStreamItem(coinType, 0-num, tid, 100, "pay"))
			err := queryRecordCoinStream(t, uid, platform, tid, streamList)
			So(err, ShouldEqual, ecode.RequestErr)
		}))

		Convey("pay but num greater or equal 0", testWithRandUid(func(uid int64) {
			tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
			platform := "pc"
			coinType := "gold"

			for _, num := range []int64{0, 1, 100} {

				streamList := make([]*model.RecordCoinStreamItem, 0)
				streamList = append(streamList, getTestRecordCoinStreamItem(coinType, num, tid, 100, "pay"))
				err := queryRecordCoinStream(t, uid, platform, tid, streamList)
				So(err, ShouldEqual, ecode.RequestErr)
			}

		}))

		Convey("recharge but num less or equal 0", testWithRandUid(func(uid int64) {
			tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
			platform := "pc"
			coinType := "gold"

			for _, num := range []int64{0, -1, -100} {

				streamList := make([]*model.RecordCoinStreamItem, 0)
				streamList = append(streamList, getTestRecordCoinStreamItem(coinType, num, tid, 100, "recharge"))
				err := queryRecordCoinStream(t, uid, platform, tid, streamList)
				So(err, ShouldEqual, ecode.RequestErr)
			}

		}))

	}))
}
