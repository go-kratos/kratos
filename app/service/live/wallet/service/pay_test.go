package service

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/conf"
	"go-common/app/service/live/wallet/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/sync/errgroup"
	"sync"
	"testing"
	"time"
)

func queryPay(t *testing.T, uid int64, platform string, coinType string, coinNum int64) (success bool, tid string, err error) {
	tid = getTestTidForCall(t, int32(model.PAYTYPE))

	bp := getTestDefaultBasicParam(tid)
	_, err = s.Pay(ctx, bp, uid, platform, getTestRechargeOrPayForm(uid, coinType, coinNum, tid))
	if err == nil {
		success = true
	}
	return

}

func queryPayWithReason(t *testing.T, uid int64, platform string, coinType string, coinNum int64, reason string) (success bool, tid string, err error) {
	tid = getTestTidForCall(t, int32(model.PAYTYPE))

	bp := getTestDefaultBasicParam(tid)
	_, err = s.Pay(ctx, bp, uid, platform, getTestRechargeOrPayForm(uid, coinType, coinNum, tid), reason)
	if err == nil {
		success = true
	}
	return

}

func testWithRandUid(f func(uid int64)) func() {
	return testWith(func() {
		uid := getTestRandUid()
		f(uid)
	})
}

func TestService_Pay(t *testing.T) {
	Convey("normal", t, testWithTestUser(func(u *TestUser) {
		for _, platform := range testValidPlatform {
			beforeWallet := getTestWallet(t, u.uid, platform)
			beforeDetail := getTestAll(t, u.uid, platform)

			var num int64 = 100
			// 加金瓜子
			success, tid, err := queryPay(t, u.uid, platform, "gold", num)
			So(success, ShouldBeTrue)
			So(err, ShouldBeNil)
			success, queryResp, err := queryQueryWithUid(t, tid, u.uid)
			t.Logf("tid:%s", tid)
			So(success, ShouldBeTrue)
			So(queryResp.Status, ShouldEqual, TX_STATUS_SUCC)
			So(err, ShouldBeNil)

			// 加银瓜子
			success, tid, err = queryPay(t, u.uid, platform, "silver", num)
			So(success, ShouldBeTrue)
			So(err, ShouldBeNil)
			success, queryResp, err = queryQueryWithUid(t, tid, u.uid)
			So(success, ShouldBeTrue)
			So(queryResp.Status, ShouldEqual, TX_STATUS_SUCC)
			So(err, ShouldBeNil)

			afterWallet := getTestWallet(t, u.uid, platform)
			afterDetail := getTestAll(t, u.uid, platform)

			So(atoiForTest(afterWallet.Gold)-atoiForTest(beforeWallet.Gold), ShouldEqual, 0-num)
			So(atoiForTest(afterWallet.Silver)-atoiForTest(beforeWallet.Silver), ShouldEqual, 0-num)

			So(atoiForTest(afterDetail.Gold)-atoiForTest(beforeDetail.Gold), ShouldEqual, 0-num)
			So(atoiForTest(afterDetail.Silver)-atoiForTest(beforeDetail.Silver), ShouldEqual, 0-num)
			So(atoiForTest(afterDetail.GoldPayCnt)-atoiForTest(beforeDetail.GoldPayCnt), ShouldEqual, num)
			So(atoiForTest(afterDetail.GoldRechargeCnt)-atoiForTest(beforeDetail.GoldRechargeCnt), ShouldEqual, 0)
			So(atoiForTest(afterDetail.SilverPayCnt)-atoiForTest(beforeDetail.SilverPayCnt), ShouldEqual, num)

		}
	}))

	Convey("logic", t, testWith(func() {

		Convey("coin not enough", testWithTestUser(func(u *TestUser) {
			var num int64 = 100
			for _, platform := range testValidPlatform {
				for _, coinType := range testLocalValidCoinType {
					_, tid, err := queryPay(t, u.uid, platform, coinType, u.coin+num)
					So(err, ShouldEqual, ecode.CoinNotEnough)
					record, err := s.dao.GetCoinStreamByTid(ctx, tid)
					So(err, ShouldBeNil)
					So(record.OpResult, ShouldEqual, model.STREAM_OP_RESULT_SUB_FAILED)
					So(record.OpReason, ShouldEqual, model.STREAM_OP_REASON_NOT_ENOUGH_COIN)
					success, queryResp, err := queryQueryWithUid(t, tid, u.uid)
					So(success, ShouldBeTrue)
					So(queryResp.Status, ShouldEqual, TX_STATUS_FAILED)
					So(err, ShouldBeNil)

				}
			}
		}))

	}))

	Convey("params", t, testWith(func() {

		Convey("coinNum invalid", testWithRandUid(func(uid int64) {
			invalidCoinNums := []int64{-1, 0, -100}
			for _, coinNum := range invalidCoinNums {
				success, tid, err := queryPay(t, uid, "pc", "gold", coinNum)
				So(success, ShouldEqual, false)
				So(err, ShouldEqual, ecode.RequestErr)

				success, resp, err := queryQueryWithUid(t, tid, uid)
				So(success, ShouldEqual, false)
				So(resp, ShouldBeNil)
				So(err, ShouldEqual, ecode.NothingFound)
			}
		}))

		Convey("uid invalid", testWith(func() {
			invalidUIDs := []int64{-1, 0, -2}
			for _, uid := range invalidUIDs {
				success, _, err := queryPay(t, uid, "pc", "gold", 100)
				So(success, ShouldEqual, false)
				So(err, ShouldEqual, ecode.RequestErr)
			}
		}))

	}))
}

func TestService_PayMulti(t *testing.T) {

	Convey("multi", t, testWith(func() {
		rate := 10
		times := rate * 2
		platforms := []string{"pc", "ios"}
		coinTypes := []string{"silver", "gold"}
		for _, platform := range platforms {
			for _, coinType := range coinTypes {
				u := initTestUser()

				s.dao.UpdateSnapShotTime(ctx, u.uid, "")
				var detail1 *model.DetailWithSnapShot
				tx, _ := s.dao.BeginTx(ctx)
				detail1, _ = s.dao.WalletForUpdate(tx, u.uid)
				t.Logf("detail1: %+v", detail1)
				So(detail1.SnapShotTime, ShouldEqual, "0001-01-01 00:00:00")
				tx.Rollback()

				payRespMap := make(map[string]*ServiceResForTest)
				coinNum := u.coin / int64(rate)
				var wg sync.WaitGroup
				var lock sync.Mutex
				for i := 0; i < times; i++ {
					wg.Add(1)
					go func() {
						localService := New(conf.Conf)
						tid := getTestTidForCall(t, int32(model.PAYTYPE))

						bp := getTestDefaultBasicParam(tid)
						v, err := localService.Pay(ctx, bp, u.uid, platform, getTestRechargeOrPayForm(u.uid, coinType, coinNum, tid))

						lock.Lock()
						payRespMap[tid] = &ServiceResForTest{V: v, Err: err}
						lock.Unlock()
						wg.Done()
					}()

					time.Sleep(time.Millisecond * 30)
				}

				wg.Wait()

				successNum := 0
				validErrs := map[error]bool{ecode.TargetBlocked: true, nil: true, ecode.CoinNotEnough: true}

				for tid, res := range payRespMap {
					_, ok := validErrs[res.Err]
					if !ok {
						t.Logf("%v", res.Err)
						t.FailNow()
					}
					So(ok, ShouldBeTrue)
					querySuccess, queryResp, queryErr := queryQueryWithUid(t, tid, u.uid)
					So(querySuccess, ShouldBeTrue)
					So(queryErr, ShouldBeNil)
					coinRecord, recordErr := s.dao.GetCoinStreamByTid(ctx, tid)
					So(recordErr, ShouldBeNil)
					if res.Err == nil {
						successNum++
						So(queryResp.Status, ShouldEqual, TX_STATUS_SUCC)
						So(coinRecord.OpResult, ShouldEqual, model.STREAM_OP_RESULT_SUB_SUCC)
						So(coinRecord.OpReason, ShouldEqual, 0)
					} else {
						So(queryResp.Status, ShouldEqual, TX_STATUS_FAILED)
						So(coinRecord.OpResult, ShouldEqual, model.STREAM_OP_RESULT_SUB_FAILED)
						if res.Err == ecode.TargetBlocked {
							So(coinRecord.OpReason, ShouldEqual, model.STREAM_OP_REASON_LOCK_FAILED)
						} else {
							So(coinRecord.OpReason, ShouldEqual, model.STREAM_OP_REASON_NOT_ENOUGH_COIN)
						}
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

				subCoin := successNum * int(coinNum)
				wallet := getTestWallet(t, u.uid, platform)
				var nowCoin int
				if coinType == "gold" {
					nowCoin = atoiForTest(wallet.Gold)
				} else if coinType == "silver" {
					nowCoin = atoiForTest(wallet.Silver)
				}

				So(nowCoin, ShouldBeGreaterThanOrEqualTo, 0)

				So(nowCoin, ShouldEqual, int(u.coin)-subCoin)

			}

		}

	}))

}

func TestService_PayForNewUser(t *testing.T) {
	Convey("new user", t, testWith(func() {
		var uid int64
		for {
			uid = getTestRandUid()
			_, err := s.dao.DetailWithoutDefault(ctx, uid)
			if err != nil {
				if err == sql.ErrNoRows {
					break
				} else {
					t.Errorf("select wrong :%s", err.Error())
					t.FailNow()
				}
			}

		}

		platform := "pc"
		var wg errgroup.Group
		times := 2
		for i := 0; i < times; i++ {
			wg.Go(func() error {
				tid := getTestTidForCall(t, int32(model.PAYTYPE))
				bp := getTestDefaultBasicParam(tid)
				_, err := s.Pay(context.Background(), bp, uid, platform, getTestRechargeOrPayForm(uid, "gold", 1, tid))
				So(err, ShouldEqual, ecode.CoinNotEnough)
				return nil
			})
		}
		wg.Wait()
		var detail2 *model.DetailWithSnapShot
		tx, _ := s.dao.BeginTx(ctx)
		detail2, _ = s.dao.WalletForUpdate(tx, uid)
		t.Logf("detail2:%+v", detail2)
		So(detail2.Gold, ShouldEqual, 0)
		So(detail2.SnapShotGold, ShouldEqual, 0)
		So(detail2.SnapShotIapGold, ShouldEqual, 0)
		So(detail2.SnapShotSilver, ShouldEqual, 0)
		So(detail2.SnapShotTime, ShouldEqual, "0001-01-01 00:00:00")

		tx.Rollback()

	}))

}

func TestService_PayMetal(t *testing.T) {
	Convey("normal", t, testWith(func() {
		var uid int64 = 1
		s.dao.ModifyMetal(ctx, uid, 2, 0, nil)
		t.Logf("test uid is %d", uid)
		platform := "pc"
		om, _ := s.dao.GetMetal(ctx, uid)
		t.Logf("before pay metal : %f", om)

		var num int64 = 1
		// 加金瓜子
		success, tid, err := queryPayWithReason(t, uid, platform, "metal", num, "ut")
		if success {
			So(err, ShouldBeNil)
			success, queryResp, queryErr := queryQueryWithUid(t, tid, uid)
			So(success, ShouldBeTrue)
			So(queryResp.Status, ShouldEqual, TX_STATUS_SUCC)
			So(queryErr, ShouldBeNil)
		} else {
			So(err, ShouldEqual, ecode.CoinNotEnough)
			success, queryResp, queryErr := queryQueryWithUid(t, tid, uid)
			So(success, ShouldBeTrue)
			So(queryResp.Status, ShouldEqual, TX_STATUS_FAILED)
			So(queryErr, ShouldBeNil)
		}

		om, _ = s.dao.GetMetal(ctx, uid)
		t.Logf("after pay metal : %f", om)
	}))
}

func TestService_PayBasicParam(t *testing.T) {
	Convey("default", t, testWithTestUser(func(u *TestUser) {
		for _, platform := range testValidPlatform {

			var num int64 = 100
			// 加金瓜子
			tid := getTestTidForCall(t, int32(model.PAYTYPE))

			bp := getTestDefaultBasicParam(tid)
			_, err := s.Pay(ctx, bp, u.uid, platform, getTestRechargeOrPayForm(u.uid, "gold", num, tid))
			So(err, ShouldBeNil)

			record, err := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			So(err, ShouldBeNil)
			oP := model.GetPlatformNo(platform)
			So(record.Platform, ShouldEqual, oP)
			So(record.Reserved1, ShouldEqual, 0)
			So(record.Uid, ShouldEqual, u.uid)
			So(record.CoinType, ShouldEqual, model.GetCoinTypeNumber(model.GetSysCoinType("gold", platform)))
			So(record.OpResult, ShouldEqual, 1)
			So(record.BizSource, ShouldEqual, "")
			So(record.DeltaCoinNum, ShouldEqual, -1*num)

		}
	}))

	Convey("advanced", t, testWithTestUser(func(u *TestUser) {
		for _, platform := range testValidPlatform {

			var num int64 = 100
			// 加金瓜子
			tid := getTestTidForCall(t, int32(model.PAYTYPE))

			bp := getTestDefaultBasicParam(tid)
			bp.Reason = 1
			bp.BizCode = "send-gift-pay"
			bp.Version = 2
			_, err := s.Pay(ctx, bp, u.uid, platform, getTestRechargeOrPayForm(u.uid, "gold", num, tid))
			So(err, ShouldBeNil)

			record, err := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			So(err, ShouldBeNil)
			oP := model.GetPlatformNo(platform)
			So(record.Platform, ShouldEqual, oP)
			So(record.Uid, ShouldEqual, u.uid)
			So(record.CoinType, ShouldEqual, model.GetCoinTypeNumber(model.GetSysCoinType("gold", platform)))
			So(record.OpResult, ShouldEqual, 1)
			So(record.BizSource, ShouldEqual, "")
			So(record.DeltaCoinNum, ShouldEqual, -1*num)
			So(record.Reserved1, ShouldEqual, bp.Reason)
			So(record.BizCode, ShouldEqual, bp.BizCode)
			So(record.Version, ShouldEqual, 2)

		}
	}))
}
