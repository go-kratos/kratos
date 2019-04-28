package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/conf"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	"sync"
	"testing"
	"time"
)

func TestService_Modify(t *testing.T) {
	Convey("normal add gold", t, testWith(func() {
		for _, platform := range testValidPlatform {
			u := initTestUser()
			beforeDetail := getTestAll(t, u.uid, platform)
			tid := getTestTidForCall(t, int32(model.RECHARGETYPE))
			bp := getTestDefaultBasicParam(tid)
			form := getTestRechargeOrPayForm(u.uid, "gold", 100, tid)
			_, err := s.Modify(ctx, bp, u.uid, platform, form)
			So(err, ShouldBeNil)
			afterDetail := getTestAll(t, u.uid, platform)
			So(atoiForTest(afterDetail.GoldRechargeCnt)-atoiForTest(beforeDetail.GoldRechargeCnt), ShouldEqual, 0)
			So(atoiForTest(afterDetail.Gold)-atoiForTest(beforeDetail.Gold), ShouldEqual, 100)

			record, recordErr := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			So(recordErr, ShouldBeNil)

			So(record.OrgCoinNum, ShouldEqual, u.coin)
			So(record.DeltaCoinNum, ShouldEqual, 100)
			So(record.OpType, ShouldEqual, int32(model.RECHARGETYPE))
			So(record.OpResult, ShouldEqual, model.STREAM_OP_RESULT_ADD_SUCC)
			success, resp, err := queryQueryWithUid(t, tid, u.uid)
			So(success, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(resp.Status, ShouldEqual, TX_STATUS_SUCC)

			_, err = s.Modify(ctx, bp, u.uid, platform, form)
			So(err, ShouldEqual, ecode.TargetBlocked)

		}

	}))

	Convey("normal decrement gold", t, testWith(func() {
		for _, platform := range testValidPlatform {
			u := initTestUser()
			beforeDetail := getTestAll(t, u.uid, platform)
			tid := getTestTidForCall(t, int32(model.PAYTYPE))
			bp := getTestDefaultBasicParam(tid)
			form := getTestRechargeOrPayForm(u.uid, "gold", -100, tid)
			_, err := s.Modify(ctx, bp, u.uid, platform, form)
			So(err, ShouldBeNil)
			afterDetail := getTestAll(t, u.uid, platform)
			So(atoiForTest(afterDetail.GoldPayCnt)-atoiForTest(beforeDetail.GoldPayCnt), ShouldEqual, 0)
			So(atoiForTest(afterDetail.Gold)-atoiForTest(beforeDetail.Gold), ShouldEqual, -100)

			record, recordErr := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			So(recordErr, ShouldBeNil)

			So(record.OrgCoinNum, ShouldEqual, u.coin)
			So(record.DeltaCoinNum, ShouldEqual, -100)
			So(record.OpType, ShouldEqual, int32(model.PAYTYPE))
			So(record.OpResult, ShouldEqual, model.STREAM_OP_RESULT_SUB_SUCC)
			success, resp, err := queryQueryWithUid(t, tid, u.uid)
			So(success, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(resp.Status, ShouldEqual, TX_STATUS_SUCC)

		}

	}))

	Convey("normal add silver", t, testWith(func() {
		for _, platform := range testValidPlatform {
			u := initTestUser()
			beforeDetail := getTestAll(t, u.uid, platform)
			tid := getTestTidForCall(t, int32(model.RECHARGETYPE))
			bp := getTestDefaultBasicParam(tid)
			form := getTestRechargeOrPayForm(u.uid, "silver", 100, tid)
			_, err := s.Modify(ctx, bp, u.uid, platform, form)
			So(err, ShouldBeNil)
			afterDetail := getTestAll(t, u.uid, platform)
			So(atoiForTest(afterDetail.SilverPayCnt)-atoiForTest(beforeDetail.SilverPayCnt), ShouldEqual, 0)
			So(atoiForTest(afterDetail.Silver)-atoiForTest(beforeDetail.Silver), ShouldEqual, 100)

			record, recordErr := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			So(recordErr, ShouldBeNil)

			So(record.OrgCoinNum, ShouldEqual, u.coin)
			So(record.DeltaCoinNum, ShouldEqual, 100)
			So(record.OpType, ShouldEqual, int32(model.RECHARGETYPE))
			So(record.CoinType, ShouldEqual, 0)
			So(record.OpResult, ShouldEqual, model.STREAM_OP_RESULT_ADD_SUCC)
			success, resp, err := queryQueryWithUid(t, tid, u.uid)
			So(success, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(resp.Status, ShouldEqual, TX_STATUS_SUCC)

		}

	}))

	Convey("normal decrement silver", t, testWith(func() {
		for _, platform := range testValidPlatform {
			u := initTestUser()
			beforeDetail := getTestAll(t, u.uid, platform)
			tid := getTestTidForCall(t, int32(model.PAYTYPE))
			bp := getTestDefaultBasicParam(tid)
			form := getTestRechargeOrPayForm(u.uid, "silver", -100, tid)
			_, err := s.Modify(ctx, bp, u.uid, platform, form)
			So(err, ShouldBeNil)
			afterDetail := getTestAll(t, u.uid, platform)
			So(atoiForTest(afterDetail.SilverPayCnt)-atoiForTest(beforeDetail.SilverPayCnt), ShouldEqual, 0)
			So(atoiForTest(afterDetail.Silver)-atoiForTest(beforeDetail.Silver), ShouldEqual, -100)

			record, recordErr := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			So(recordErr, ShouldBeNil)

			So(record.OrgCoinNum, ShouldEqual, u.coin)
			So(record.DeltaCoinNum, ShouldEqual, -100)
			So(record.OpType, ShouldEqual, int32(model.PAYTYPE))
			So(record.CoinType, ShouldEqual, 0)
			So(record.OpResult, ShouldEqual, model.STREAM_OP_RESULT_SUB_SUCC)
			success, resp, err := queryQueryWithUid(t, tid, u.uid)
			So(success, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(resp.Status, ShouldEqual, TX_STATUS_SUCC)

		}

	}))

	Convey("params", t, testWith(func() {
		Convey("coinNum equal 0", testWithTestUser(func(u *TestUser) {
			tid := getTestTidForCall(t, int32(model.RECHARGETYPE))
			platform := "pc"
			bp := getTestDefaultBasicParam(tid)
			form := getTestRechargeOrPayForm(u.uid, "silver", 0, tid)
			_, err := s.Modify(ctx, bp, u.uid, platform, form)
			So(err, ShouldEqual, ecode.RequestErr)
			_, recordErr := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			So(recordErr, ShouldNotBeNil)

		}))

		Convey("not local coin", testWithTestUser(func(u *TestUser) {
			tid := getTestTidForCall(t, int32(model.RECHARGETYPE))
			platform := "pc"
			bp := getTestDefaultBasicParam(tid)
			form := getTestRechargeOrPayForm(u.uid, "metal", 100, tid)
			_, err := s.Modify(ctx, bp, u.uid, platform, form)
			So(err, ShouldEqual, ecode.RequestErr)
			_, recordErr := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			So(recordErr, ShouldNotBeNil)
		}))
	}))

	Convey("logic", t, testWith(func() {
		Convey("not enough", testWith(func() {
			for _, platform := range testValidPlatform {
				for _, coinType := range testLocalValidCoinType {
					u := initTestUser()
					tid := getTestTidForCall(t, int32(model.PAYTYPE))
					bp := getTestDefaultBasicParam(tid)
					form := getTestRechargeOrPayForm(u.uid, coinType, -1001, tid)
					_, err := s.Modify(ctx, bp, u.uid, platform, form)
					So(err, ShouldEqual, ecode.CoinNotEnough)
					coinRecord, recordErr := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
					So(recordErr, ShouldBeNil)
					So(coinRecord.OpResult, ShouldEqual, model.STREAM_OP_RESULT_SUB_FAILED)
					So(coinRecord.OpReason, ShouldEqual, model.STREAM_OP_REASON_NOT_ENOUGH_COIN)
					success, resp, err := queryQueryWithUid(t, tid, u.uid)
					So(success, ShouldBeTrue)
					So(err, ShouldBeNil)
					So(resp.Status, ShouldEqual, TX_STATUS_FAILED)

					wallet := getTestWallet(t, u.uid, platform)
					So(wallet.Gold, ShouldEqual, "1000")

				}
			}

		}))
	}))

	Convey("multi", t, testWith(func() {
		for _, platform := range []string{"pc", "ios"} {
			for _, coinType := range []string{"gold", "silver"} {

				u := initTestUser()

				var detail1 *model.DetailWithSnapShot
				s.dao.UpdateSnapShotTime(ctx, u.uid, "")
				tx, _ := s.dao.BeginTx(ctx)
				detail1, _ = s.dao.WalletForUpdate(tx, u.uid)
				t.Logf("detail1: %+v", detail1)
				So(detail1.SnapShotTime, ShouldEqual, "0001-01-01 00:00:00")
				tx.Rollback()

				beforeDetail := getTestAll(t, u.uid, platform)

				var step int64 = 10
				num := u.coin / step
				times := int(step * 2)

				var wg sync.WaitGroup
				var lock sync.Mutex

				modifyRespMap := make(map[string]*ServiceResForTest)

				for i := 0; i < times; i++ {
					wg.Add(1)
					go func() {
						localService := New(conf.Conf)
						tid := getTestTidForCall(t, int32(model.PAYTYPE))
						bp := getTestDefaultBasicParam(tid)
						form := getTestRechargeOrPayForm(u.uid, coinType, num*-1, tid)
						v, err := localService.Modify(ctx, bp, u.uid, platform, form)
						lock.Lock()
						modifyRespMap[tid] = &ServiceResForTest{Err: err, V: v}
						lock.Unlock()
						wg.Done()

					}()
					time.Sleep(time.Millisecond * 30)
				}
				wg.Wait()

				successNum := 0

				validErrs := map[error]bool{ecode.TargetBlocked: true, nil: true, ecode.CoinNotEnough: true}
				for tid, res := range modifyRespMap {
					_, ok := validErrs[res.Err]
					So(ok, ShouldBeTrue)
					if res.Err == nil {
						assertNormalCoinStream(t, platform, coinType, tid, 0, num*-1, model.PAYTYPE, model.STREAM_OP_RESULT_SUB_SUCC, TX_STATUS_SUCC)
						successNum++
					} else if res.Err == ecode.TargetBlocked {
						record := assertNormalCoinStream(t, platform, coinType, tid, 0, num*-1, model.PAYTYPE, model.STREAM_OP_RESULT_SUB_FAILED, TX_STATUS_FAILED)
						So(record.OpReason, ShouldEqual, model.STREAM_OP_REASON_LOCK_FAILED)
					} else {
						record := assertNormalCoinStream(t, platform, coinType, tid, 0, num*-1, model.PAYTYPE, model.STREAM_OP_RESULT_SUB_FAILED, TX_STATUS_FAILED)
						So(record.OpReason, ShouldEqual, model.STREAM_OP_REASON_NOT_ENOUGH_COIN)
					}
				}

				t.Logf("successNum:%d", successNum)

				var detail2 *model.DetailWithSnapShot
				tx, _ = s.dao.BeginTx(ctx)

				detail2, _ = s.dao.WalletForUpdate(tx, u.uid)
				t.Logf("detail2: %+v", detail2)
				So(detail2.SnapShotGold, ShouldEqual, detail1.Gold)
				So(detail2.SnapShotIapGold, ShouldEqual, detail1.IapGold)
				So(detail2.SnapShotSilver, ShouldEqual, detail1.Silver)
				So(detail2.SnapShotTime, ShouldNotEqual, "0001-01-01 00:00:00")

				tx.Rollback()

				afterDetail := getTestAll(t, u.uid, platform)

				So(getCoinNumByDetailForTest(afterDetail, coinType)-getCoinNumByDetailForTest(beforeDetail, coinType), ShouldEqual, -1*successNum*int(num))
				So(atoiForTest(afterDetail.GoldPayCnt)-atoiForTest(beforeDetail.GoldPayCnt), ShouldEqual, 0)
				So(atoiForTest(afterDetail.GoldRechargeCnt)-atoiForTest(beforeDetail.GoldRechargeCnt), ShouldEqual, 0)
				So(atoiForTest(afterDetail.SilverPayCnt)-atoiForTest(beforeDetail.SilverPayCnt), ShouldEqual, 0)

				u.rest()

			}
		}
	}))
}

func TestService_ModifyBasicParam(t *testing.T) {
	Convey("default", t, testWithTestUser(func(u *TestUser) {
		for _, platform := range testValidPlatform {
			tid := getTestTidForCall(t, int32(model.RECHARGETYPE))
			bp := getTestDefaultBasicParam(tid)
			form := getTestRechargeOrPayForm(u.uid, "gold", 1, tid)
			_, err := s.Modify(ctx, bp, u.uid, platform, form)
			So(err, ShouldBeNil)
			record, err := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			oP := model.GetPlatformNo(platform)
			So(record.Platform, ShouldEqual, oP)
			So(record.Reserved1, ShouldEqual, 0)
			So(record.Uid, ShouldEqual, u.uid)
			So(record.CoinType, ShouldEqual, model.GetCoinTypeNumber(model.GetSysCoinType("gold", platform)))
			So(record.OpResult, ShouldEqual, 2)
			So(record.BizSource, ShouldEqual, "")
			So(record.DeltaCoinNum, ShouldEqual, 1)
		}

	}))

	Convey("advanced", t, testWithTestUser(func(u *TestUser) {
		for _, platform := range testValidPlatform {
			tid := getTestTidForCall(t, int32(model.RECHARGETYPE))
			bp := getTestDefaultBasicParam(tid)
			bp.Reason = 1
			bp.BizCode = "send-gift"
			form := getTestRechargeOrPayForm(u.uid, "gold", 1, tid)
			_, err := s.Modify(ctx, bp, u.uid, platform, form)
			So(err, ShouldBeNil)
			record, err := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			oP := model.GetPlatformNo(platform)
			So(record.Platform, ShouldEqual, oP)
			So(record.Reserved1, ShouldEqual, bp.Reason)
			So(record.Uid, ShouldEqual, u.uid)
			So(record.CoinType, ShouldEqual, model.GetCoinTypeNumber(model.GetSysCoinType("gold", platform)))
			So(record.OpResult, ShouldEqual, 2)
			So(record.BizSource, ShouldEqual, "")
			So(record.DeltaCoinNum, ShouldEqual, 1)
			So(record.BizCode, ShouldEqual, bp.BizCode)
		}

	}))
}
