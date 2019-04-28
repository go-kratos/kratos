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

func queryExchange(t *testing.T, uid int64, platform string, srcCoinType string, srcCoinNum int64, destCoinType string, destCoinNum int64) (success bool, tid string, err error) {
	tid = getTestTidForCall(t, int32(model.EXCHANGETYPE))

	bp := getTestDefaultBasicParam(tid)
	_, err = s.Exchange(ctx, bp, uid, platform, getTestExchangeForm(uid, srcCoinType, srcCoinNum, destCoinType, destCoinNum, tid))
	if err == nil {
		success = true
	}
	return

}

func assertNormalCoinStream(t *testing.T, platform string, coinType string, tid string, offset int, deltaCoinNum int64, serviceType model.ServiceType, result int, status int) *model.CoinStreamRecord {
	record, recordErr := s.dao.GetCoinStreamByTidAndOffset(ctx, tid, offset)
	So(recordErr, ShouldBeNil)

	So(record.DeltaCoinNum, ShouldEqual, deltaCoinNum)
	So(record.OpType, ShouldEqual, int32(serviceType))
	So(record.OpResult, ShouldEqual, result)

	srcSysCoinType := model.GetSysCoinType(coinType, platform)
	srcSysCoinTypeNo := model.GetCoinTypeNumber(srcSysCoinType)
	So(record.CoinType, ShouldEqual, srcSysCoinTypeNo)

	success, resp, err := queryQueryWithUid(t, tid, record.Uid)
	So(success, ShouldBeTrue)
	So(err, ShouldBeNil)
	So(resp.Status, ShouldEqual, status)
	return record

}

func getCoinNumByMelonRespForTest(wallet *model.MelonseedResp, coinType string) int {

	switch coinType {
	case "silver":
		return atoiForTest(wallet.Silver)
	case "gold":
		return atoiForTest(wallet.Gold)
	default:
		return 0
	}
}

func getCoinNumByDetailForTest(wallet *model.DetailResp, coinType string) int {

	switch coinType {
	case "silver":
		return atoiForTest(wallet.Silver)
	case "gold":
		return atoiForTest(wallet.Gold)
	default:
		return 0
	}
}

func TestService_Exchange(t *testing.T) {
	Convey("normal", t, testWith(func() {
		Convey("gold2silver", testWith(func() {
			platforms := []string{"pc", "ios"}
			for _, platform := range platforms {
				u := initTestUser()

				s.dao.UpdateSnapShotTime(ctx, u.uid, "")
				beforeDetail := getTestAll(t, u.uid, platform)
				success, tid, err := queryExchange(t, u.uid, platform, "gold", 10, "silver", 100)
				So(success, ShouldBeTrue)
				So(err, ShouldBeNil)

				wallet := getTestWallet(t, u.uid, platform)

				So(atoiForTest(wallet.Silver), ShouldEqual, u.coin+100)
				So(atoiForTest(wallet.Gold), ShouldEqual, u.coin-10)
				melon, _ := s.dao.Melonseed(ctx, u.uid)
				var gold int64
				if platform == "ios" {
					gold = melon.IapGold
				} else {
					gold = melon.Gold
				}
				So(gold, ShouldEqual, u.coin-10)

				// 检查CoinStream表
				assertNormalCoinStream(t, platform, "gold", tid, 1, -10, model.EXCHANGETYPE, model.STREAM_OP_RESULT_SUB_SUCC, TX_STATUS_SUCC)
				assertNormalCoinStream(t, platform, "silver", tid, 0, 100, model.EXCHANGETYPE, model.STREAM_OP_RESULT_ADD_SUCC, TX_STATUS_SUCC)
				u.rest()

				// check detail
				afterDetail := getTestAll(t, u.uid, platform)

				So(atoiForTest(afterDetail.GoldPayCnt)-atoiForTest(beforeDetail.GoldPayCnt), ShouldEqual, 10)

				var detail2 *model.DetailWithSnapShot
				tx, _ := s.dao.BeginTx(ctx)

				detail2, _ = s.dao.WalletForUpdate(tx, u.uid)
				t.Logf("detail2: %+v", detail2)
				So(detail2.SnapShotGold, ShouldEqual, 1000)
				So(detail2.SnapShotIapGold, ShouldEqual, 1000)
				So(detail2.SnapShotSilver, ShouldEqual, 1000)
				So(detail2.SnapShotTime, ShouldNotEqual, "0001-01-01 00:00:00")

				tx.Rollback()

			}

		}))

		Convey("any2any", testWith(func() {
			platforms := []string{"pc", "ios"}
			coinTypes := []string{"gold", "silver"}

			for _, platform := range platforms {
				for _, srcCoinType := range coinTypes {
					for _, destCoinType := range coinTypes {
						if srcCoinType == destCoinType {
							continue
						}
						u := initTestUser()
						var srcCoinNum int64 = 10
						var destCoinNum int64 = 100
						success, tid, err := queryExchange(t, u.uid, platform, srcCoinType, srcCoinNum, destCoinType, destCoinNum)
						So(success, ShouldBeTrue)
						So(err, ShouldBeNil)

						wallet := getTestWallet(t, u.uid, platform)

						So(getCoinNumByMelonRespForTest(wallet, destCoinType), ShouldEqual, u.coin+destCoinNum)
						So(getCoinNumByMelonRespForTest(wallet, srcCoinType), ShouldEqual, u.coin-srcCoinNum)

						// 检查CoinStream表
						assertNormalCoinStream(t, platform, srcCoinType, tid, 1, srcCoinNum*-1, model.EXCHANGETYPE, model.STREAM_OP_RESULT_SUB_SUCC, TX_STATUS_SUCC)
						assertNormalCoinStream(t, platform, destCoinType, tid, 0, destCoinNum, model.EXCHANGETYPE, model.STREAM_OP_RESULT_ADD_SUCC, TX_STATUS_SUCC)
						u.rest()
					}
				}
			}
		}))
	}))

	Convey("params", t, testWith(func() {
		Convey("same coinType", testWithRandUid(func(uid int64) {
			success, _, err := queryExchange(t, uid, "pc", "gold", 10, "gold", 100)
			So(success, ShouldBeFalse)
			So(err, ShouldEqual, ecode.RequestErr)
		}))

		Convey("src or dest coin less than or equal 0", testWithRandUid(func(uid int64) {
			invalidCoins := []int64{0, -1, -100}
			for _, coin := range invalidCoins {
				success, _, err := queryExchange(t, uid, "pc", "gold", coin, "gold", 100)
				So(success, ShouldBeFalse)
				So(err, ShouldEqual, ecode.RequestErr)
			}

			for _, coin := range invalidCoins {
				success, _, err := queryExchange(t, uid, "pc", "gold", 10, "gold", coin)
				So(success, ShouldBeFalse)
				So(err, ShouldEqual, ecode.RequestErr)
			}
		}))

	}))

	Convey("not enough", t, testWithTestUser(func(u *TestUser) {

		s.dao.UpdateSnapShotTime(ctx, u.uid, "")

		success, tid, err := queryExchange(t, u.uid, "pc", "silver", 1001, "gold", 100)
		So(success, ShouldBeFalse)
		So(err, ShouldEqual, ecode.CoinNotEnough)

		stream := assertNormalCoinStream(t, "pc", "silver", tid, 0, -1001, model.EXCHANGETYPE, model.STREAM_OP_RESULT_SUB_FAILED, TX_STATUS_FAILED)
		So(stream.OpReason, ShouldEqual, 1)

		var detail2 *model.DetailWithSnapShot
		tx, _ := s.dao.BeginTx(ctx)

		detail2, _ = s.dao.WalletForUpdate(tx, u.uid)
		t.Logf("detail2: %+v", detail2)
		So(detail2.SnapShotTime, ShouldEqual, "0001-01-01 00:00:00")
		So(detail2.Silver, ShouldEqual, 1000)
		So(detail2.Gold, ShouldEqual, 1000)
		So(detail2.IapGold, ShouldEqual, 1000)
		tx.Rollback()
		t.Logf("uid:%d,tid:%s", u.uid, tid)

	}))

}

func TestService_ExchangeMetal(t *testing.T) {
	Convey("enough", t, testWithTestUser(func(u *TestUser) {
		_, _, err := s.dao.ModifyMetal(ctx, u.uid, 1, 0, nil)
		if err != nil {
			t.Errorf("add metal failed err : %s", err.Error())
			t.FailNow()
		}

		metal1, _ := s.dao.GetMetal(ctx, u.uid)

		success, tid, err := queryExchange(t, u.uid, "pc", "metal", 1, "gold", 100)
		So(success, ShouldBeTrue)
		So(err, ShouldBeNil)

		assertNormalCoinStream(t, "pc", "metal", tid, 1, -1, model.EXCHANGETYPE, model.STREAM_OP_RESULT_SUB_SUCC, TX_STATUS_SUCC)
		assertNormalCoinStream(t, "pc", "gold", tid, 0, 100, model.EXCHANGETYPE, model.STREAM_OP_RESULT_ADD_SUCC, TX_STATUS_SUCC)

		metal2, _ := s.dao.GetMetal(ctx, u.uid)
		So(metal2-metal1, ShouldEqual, -1)
		t.Logf("uid:%d,tid:%s", u.uid, tid)

	}))

	Convey("not enough", t, testWithTestUser(func(u *TestUser) {

		metal1, _ := s.dao.GetMetal(ctx, u.uid)

		success, tid, err := queryExchange(t, u.uid, "pc", "metal", int64(metal1+1), "gold", 100)
		So(success, ShouldBeFalse)
		So(err, ShouldEqual, ecode.CoinNotEnough)

		assertNormalCoinStream(t, "pc", "metal", tid, 0, -int64(metal1+1), model.EXCHANGETYPE, model.STREAM_OP_RESULT_SUB_FAILED, TX_STATUS_FAILED)

		metal2, _ := s.dao.GetMetal(ctx, u.uid)
		So(metal2-metal1, ShouldEqual, 0)
		t.Logf("uid:%d,tid:%s", u.uid, tid)

	}))
}

func TestExchangeHandler_Multi(t *testing.T) {
	Convey("multi", t, testWith(func() {
		platforms := []string{"pc", "ios"}
		coinTypes := []string{"gold", "silver"}
		var destCoinNum int64 = 1000

		for _, platform := range platforms {
			for _, srcCoinType := range coinTypes {
				for _, destCoinType := range coinTypes {
					if srcCoinType == destCoinType {
						continue
					}
					u := initTestUser()

					s.dao.UpdateSnapShotTime(ctx, u.uid, "")
					var detail1 *model.DetailWithSnapShot
					tx, _ := s.dao.BeginTx(ctx)
					detail1, _ = s.dao.WalletForUpdate(tx, u.uid)
					t.Logf("detail1: %+v", detail1)
					So(detail1.SnapShotTime, ShouldEqual, "0001-01-01 00:00:00")
					tx.Rollback()

					step := 10
					num := u.coin / int64(step)
					times := step * 3

					exchangeRespMap := make(map[string]*ServiceResForTest)
					var wg sync.WaitGroup
					var lock sync.Mutex
					for i := 0; i < times; i++ {
						wg.Add(1)
						go func() {
							localService := New(conf.Conf)
							tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
							bp := getTestDefaultBasicParam(tid)
							v, err := localService.Exchange(ctx, bp, u.uid, platform, getTestExchangeForm(u.uid, srcCoinType, num, destCoinType, destCoinNum, tid))
							lock.Lock()
							exchangeRespMap[tid] = &ServiceResForTest{V: v, Err: err}
							lock.Unlock()
							wg.Done()
						}()
						time.Sleep(time.Millisecond * 30)
					}

					wg.Wait()

					successNum := 0
					validErrs := map[error]bool{ecode.TargetBlocked: true, nil: true, ecode.CoinNotEnough: true}

					for tid, res := range exchangeRespMap {
						_, ok := validErrs[res.Err]
						So(ok, ShouldBeTrue)
						if res.Err == nil {
							successNum++
						}

						if res.Err == nil {
							assertNormalCoinStream(t, platform, srcCoinType, tid, 1, num*-1, model.EXCHANGETYPE, model.STREAM_OP_RESULT_SUB_SUCC, TX_STATUS_SUCC)
							assertNormalCoinStream(t, platform, destCoinType, tid, 0, destCoinNum, model.EXCHANGETYPE, model.STREAM_OP_RESULT_ADD_SUCC, TX_STATUS_SUCC)
						} else if res.Err == ecode.CoinNotEnough {
							record := assertNormalCoinStream(t, platform, srcCoinType, tid, 0, num*-1, model.EXCHANGETYPE, model.STREAM_OP_RESULT_SUB_FAILED, TX_STATUS_FAILED)
							So(record.OpReason, ShouldEqual, model.STREAM_OP_REASON_NOT_ENOUGH_COIN)
						} else {
							record := assertNormalCoinStream(t, platform, srcCoinType, tid, 0, num*-1, model.EXCHANGETYPE, model.STREAM_OP_RESULT_SUB_FAILED, TX_STATUS_FAILED)
							So(record.OpReason, ShouldEqual, model.STREAM_OP_REASON_LOCK_FAILED)
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

					So(successNum, ShouldBeLessThanOrEqualTo, step)
					So(successNum, ShouldBeGreaterThan, 0)
					wallet := getTestWallet(t, u.uid, platform)
					So(getCoinNumByMelonRespForTest(wallet, srcCoinType)+successNum*int(num), ShouldEqual, u.coin)
					So(getCoinNumByMelonRespForTest(wallet, destCoinType)-successNum*int(destCoinNum), ShouldEqual, u.coin)
					u.rest()

				}
			}
		}

	}))
}

func TestService_ExchangeBasicParam(t *testing.T) {
	Convey("default", t, testWithTestUser(func(u *TestUser) {
		for _, platform := range testValidPlatform {
			tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
			bp := getTestDefaultBasicParam(tid)
			form := getTestExchangeForm(u.uid, "gold", 1, "silver", 1, tid)
			_, err := s.Exchange(ctx, bp, u.uid, platform, form)
			So(err, ShouldBeNil)

			record, err := s.dao.GetCoinStreamByTidAndOffset(ctx, tid, 0)
			oP := model.GetPlatformNo(platform)
			So(record.Platform, ShouldEqual, oP)
			So(record.Reserved1, ShouldEqual, 0)
			So(record.Uid, ShouldEqual, u.uid)
			So(record.CoinType, ShouldEqual, model.GetCoinTypeNumber(model.GetSysCoinType("silver", platform)))
			So(record.OpResult, ShouldEqual, 2)
			So(record.BizSource, ShouldEqual, "")
			So(record.DeltaCoinNum, ShouldEqual, 1)

			record, err = s.dao.GetCoinStreamByTidAndOffset(ctx, tid, 1)
			oP = model.GetPlatformNo(platform)
			So(record.Platform, ShouldEqual, oP)
			So(record.Reserved1, ShouldEqual, 0)
			So(record.Uid, ShouldEqual, u.uid)
			So(record.CoinType, ShouldEqual, model.GetCoinTypeNumber(model.GetSysCoinType("gold", platform)))
			So(record.OpResult, ShouldEqual, 1)
			So(record.BizSource, ShouldEqual, "")
			So(record.DeltaCoinNum, ShouldEqual, -1)
		}
	}))

	Convey("advanced", t, testWithTestUser(func(u *TestUser) {
		for _, platform := range testValidPlatform {
			tid := getTestTidForCall(t, int32(model.EXCHANGETYPE))
			bp := getTestDefaultBasicParam(tid)
			bp.BizCode = "exchange"
			bp.Reason = 2
			form := getTestExchangeForm(u.uid, "gold", 1, "silver", 1, tid)
			_, err := s.Exchange(ctx, bp, u.uid, platform, form)
			So(err, ShouldBeNil)

			record, err := s.dao.GetCoinStreamByTidAndOffset(ctx, tid, 0)
			oP := model.GetPlatformNo(platform)
			So(record.Platform, ShouldEqual, oP)
			So(record.Reserved1, ShouldEqual, 2)
			So(record.Uid, ShouldEqual, u.uid)
			So(record.CoinType, ShouldEqual, model.GetCoinTypeNumber(model.GetSysCoinType("silver", platform)))
			So(record.OpResult, ShouldEqual, 2)
			So(record.BizSource, ShouldEqual, "")
			So(record.DeltaCoinNum, ShouldEqual, 1)
			So(record.BizCode, ShouldEqual, bp.BizCode)

			record, err = s.dao.GetCoinStreamByTidAndOffset(ctx, tid, 1)
			oP = model.GetPlatformNo(platform)
			So(record.Platform, ShouldEqual, oP)
			So(record.Reserved1, ShouldEqual, 2)
			So(record.Uid, ShouldEqual, u.uid)
			So(record.CoinType, ShouldEqual, model.GetCoinTypeNumber(model.GetSysCoinType("gold", platform)))
			So(record.OpResult, ShouldEqual, 1)
			So(record.BizSource, ShouldEqual, "")
			So(record.DeltaCoinNum, ShouldEqual, -1)
			So(record.BizCode, ShouldEqual, bp.BizCode)
		}
	}))
}
