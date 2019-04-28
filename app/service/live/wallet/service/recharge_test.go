package service

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"

	"go-common/app/service/live/wallet/conf"
	"go-common/app/service/live/wallet/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/sync/errgroup"
)

func getTestExtendTid() string {
	return fmt.Sprintf("test-service:ex:%d", r.Int31n(1000000))
}

func getTestRechargeOrPayForm(uid int64, coinType string, coinNum int64, tid string) *model.RechargeOrPayForm {

	return &model.RechargeOrPayForm{
		Uid:           uid,
		CoinType:      coinType,
		CoinNum:       coinNum,
		ExtendTid:     getTestExtendTid(),
		Timestamp:     time.Now().Unix(),
		TransactionId: tid,
	}
}
func getTestExchangeForm(uid int64, srcCoinType string, srcCoinNum int64, destCoinType string, destCoinNum int64, tid string) *model.ExchangeForm {

	return &model.ExchangeForm{
		Uid:           uid,
		SrcCoinType:   srcCoinType,
		SrcCoinNum:    srcCoinNum,
		ExtendTid:     getTestExtendTid(),
		Timestamp:     time.Now().Unix(),
		TransactionId: tid,
		DestCoinNum:   destCoinNum,
		DestCoinType:  destCoinType,
	}
}
func getTestTidForCall(t *testing.T, serviceType int32) string {
	bp := getTestDefaultBasicParam("")
	v, err := s.GetTid(ctx, bp, 0, serviceType, getTestParamsJson())
	if err != nil {
		t.FailNow()
	}
	tidResp := v.(*model.TidResp)
	tid := tidResp.TransactionId
	return tid
}

var (
	testLocalValidCoinType = []string{"gold", "silver"}
)

// 增加一个初始化用户的功能　方便测试

type TestUser struct {
	uid           int64
	coin          int64
	originGold    int64
	originSilver  int64
	originIapGold int64
}

func initTestUser() *TestUser {
	uid := getTestRandUid()
	var coin int64 = 1000
	wallet, _ := s.dao.Melonseed(ctx, uid)
	needGold := coin - wallet.Gold
	s.dao.AddGold(ctx, uid, int(needGold))
	needIapGold := coin - wallet.IapGold
	s.dao.AddIapGold(ctx, uid, int(needIapGold))
	needSilver := coin - wallet.Silver
	s.dao.AddSilver(ctx, uid, int(needSilver))
	s.dao.ChangeCostBase(ctx, uid, coin)
	s.dao.DelWalletCache(ctx, uid)
	return &TestUser{
		uid:           uid,
		coin:          coin,
		originGold:    wallet.Gold,
		originSilver:  wallet.Silver,
		originIapGold: wallet.IapGold,
	}
}

func (u *TestUser) rest() {
	uid := u.uid
	wallet, _ := s.dao.Melonseed(ctx, uid)
	needGold := u.originGold - wallet.Gold
	s.dao.AddGold(ctx, uid, int(needGold))
	needIapGold := u.originIapGold - wallet.IapGold
	s.dao.AddIapGold(ctx, uid, int(needIapGold))
	needSilver := u.originSilver - wallet.Silver
	s.dao.AddSilver(ctx, uid, int(needSilver))
	s.dao.DelWalletCache(ctx, uid)
}

func testWithTestUser(f func(u *TestUser)) func() {
	once.Do(startService)
	u := initTestUser()
	return func() {
		f(u)
		u.rest()
	}
}

func queryRecharge(t *testing.T, uid int64, platform string, coinType string, coinNum int64) (success bool, tid string) {
	tid = getTestTidForCall(t, int32(model.RECHARGETYPE))

	bp := getTestDefaultBasicParam(tid)
	_, err := s.Recharge(ctx, bp, uid, platform, getTestRechargeOrPayForm(uid, coinType, coinNum, tid))
	if err != nil {
		t.FailNow()
	}
	success = true
	return

}

func TestService_Recharge(t *testing.T) {
	Convey("normal", t, testWith(func() {

		Convey("basic", func() {
			var num int64 = 1000
			for _, platform := range testValidPlatform {
				uid := getTestRandUid()
				beforeWallet := getTestWallet(t, uid, platform)
				beforeDetail := getTestAll(t, uid, platform)
				for _, coinType := range testLocalValidCoinType {
					tid := getTestTidForCall(t, int32(model.RECHARGETYPE))

					bp := getTestDefaultBasicParam(tid)

					v, err := s.Recharge(ctx, bp, uid, platform, getTestRechargeOrPayForm(uid, coinType, num, tid))
					So(err, ShouldBeNil)
					So(v, ShouldNotBeNil)

				}

				afterWallet := getTestWallet(t, uid, platform)
				afterDetail := getTestAll(t, uid, platform)

				So(atoiForTest(afterWallet.Gold)-atoiForTest(beforeWallet.Gold), ShouldEqual, num)
				So(atoiForTest(afterWallet.Silver)-atoiForTest(beforeWallet.Silver), ShouldEqual, num)

				So(atoiForTest(afterDetail.Gold)-atoiForTest(beforeDetail.Gold), ShouldEqual, num)
				So(atoiForTest(afterDetail.Silver)-atoiForTest(beforeDetail.Silver), ShouldEqual, num)
				So(atoiForTest(afterDetail.GoldRechargeCnt)-atoiForTest(beforeDetail.GoldRechargeCnt), ShouldEqual, num)
				So(atoiForTest(afterDetail.GoldPayCnt)-atoiForTest(beforeDetail.GoldPayCnt), ShouldEqual, 0)
				So(atoiForTest(afterDetail.SilverPayCnt)-atoiForTest(beforeDetail.SilverPayCnt), ShouldEqual, 0)

			}
		})
	}))

	Convey("params", t, testWith(func() {

		Convey("invalid uid", func() {
			invalidUids := []int64{-1, 0, -100}
			platform := "ios"
			coinType := "gold"
			var num int64 = 100
			for _, uid := range invalidUids {
				tid := getTestTidForCall(t, int32(model.RECHARGETYPE))

				bp := getTestDefaultBasicParam(tid)

				v, err := s.Recharge(ctx, bp, uid, platform, getTestRechargeOrPayForm(uid, coinType, num, tid))
				So(err, ShouldNotBeNil)
				So(errors.Cause(err).(ecode.Codes), ShouldEqual, ecode.RequestErr)
				So(v, ShouldBeNil)
			}
		})

		Convey("invalid coinNum", func() {
			uid := getTestRandUid()
			platform := "pc"
			coinType := "silver"
			invalidCoinNums := []int64{0, -1}

			for _, num := range invalidCoinNums {
				tid := getTestTidForCall(t, int32(model.RECHARGETYPE))

				bp := getTestDefaultBasicParam(tid)

				v, err := s.Recharge(ctx, bp, uid, platform, getTestRechargeOrPayForm(uid, coinType, num, tid))
				So(err, ShouldNotBeNil)
				So(errors.Cause(err).(ecode.Codes), ShouldEqual, ecode.RequestErr)
				So(v, ShouldBeNil)
			}
		})

		Convey("invalid coinType", func() {
			uid := getTestRandUid()
			platform := "pc"
			var num int64 = 100
			invalidCoinTypes := []string{"unknown", "metal"}

			for _, coinType := range invalidCoinTypes {
				tid := getTestTidForCall(t, int32(model.RECHARGETYPE))

				bp := getTestDefaultBasicParam(tid)

				v, err := s.Recharge(ctx, bp, uid, platform, getTestRechargeOrPayForm(uid, coinType, num, tid))
				So(err, ShouldNotBeNil)
				So(errors.Cause(err).(ecode.Codes), ShouldEqual, ecode.RequestErr)
				So(v, ShouldBeNil)
			}
		})

		Convey("invalid ExtendTid", func() {
			uid := getTestRandUid()
			platform := "pc"
			var num int64 = 100
			coinType := "gold"
			tid := getTestTidForCall(t, int32(model.RECHARGETYPE))

			bp := getTestDefaultBasicParam(tid)

			form := getTestRechargeOrPayForm(uid, coinType, num, tid)
			form.ExtendTid = ""

			v, err := s.Recharge(ctx, bp, uid, platform, form)
			So(err, ShouldNotBeNil)
			So(errors.Cause(err).(ecode.Codes), ShouldEqual, ecode.RequestErr)
			So(v, ShouldBeNil)
		})

		Convey("invalid tid", func() {
			uid := getTestRandUid()
			platform := "pc"
			var num int64 = 100
			coinType := "gold"
			tid := getTestTidForCall(t, int32(model.RECHARGETYPE))

			bp := getTestDefaultBasicParam(tid)

			bp.TransactionId = ""

			form := getTestRechargeOrPayForm(uid, coinType, num, tid)
			form.TransactionId = ""

			v, err := s.Recharge(ctx, bp, uid, platform, form)
			So(err, ShouldNotBeNil)
			So(errors.Cause(err).(ecode.Codes), ShouldEqual, ecode.RequestErr)
			So(v, ShouldBeNil)
		})

		Convey("same tid twice", func() {
			uid := getTestRandUid()
			platform := "pc"
			var num int64 = 100
			coinType := "gold"
			tid := getTestTidForCall(t, int32(model.RECHARGETYPE))

			bp := getTestDefaultBasicParam(tid)
			form := getTestRechargeOrPayForm(uid, coinType, num, tid)
			_, err := s.Recharge(ctx, bp, uid, platform, form)
			So(err, ShouldBeNil)

			_, err = s.Recharge(ctx, bp, uid, platform, form)
			So(errors.Cause(err).(ecode.Codes), ShouldEqual, ecode.TargetBlocked)

		})
	}))

	Convey("multi", t, func() {
		Convey("basic", func() {
			var num int64 = 1000
			times := 10
			uid := getTestRandUid()
			for _, platform := range testValidPlatform {
				for _, coinType := range testLocalValidCoinType {
					beforeWallet := getTestWallet(t, uid, platform)
					beforeDetail := getTestAll(t, uid, platform)
					s.dao.UpdateSnapShotTime(ctx, uid, "")

					var detail1 *model.DetailWithSnapShot
					tx, _ := s.dao.BeginTx(ctx)
					detail1, _ = s.dao.WalletForUpdate(tx, uid)
					t.Logf("detail1: %+v", detail1)
					So(detail1.SnapShotTime, ShouldEqual, "0001-01-01 00:00:00")
					tx.Rollback()

					var lock sync.Mutex
					var wg sync.WaitGroup

					rechargeRespMap := make(map[string]*ServiceResForTest)

					for i := 0; i < times; i++ {
						wg.Add(1)
						go func(index int) {
							localService := New(conf.Conf)
							bp := getTestDefaultBasicParam("")
							v, _ := localService.GetTid(ctx, bp, 0, int32(model.RECHARGETYPE), getTestParamsJson())
							tidResp := v.(*model.TidResp)
							tid := tidResp.TransactionId
							bp = getTestDefaultBasicParam(tid)
							v, err := s.Recharge(ctx, bp, uid, platform, getTestRechargeOrPayForm(uid, coinType, num, tid))
							lock.Lock()
							rechargeRespMap[tid] = &ServiceResForTest{V: v, Err: err}
							lock.Unlock()
							wg.Done()
						}(i)
						time.Sleep(time.Millisecond * 30)
					}

					wg.Wait()

					validErrs := map[error]bool{ecode.TargetBlocked: true, nil: true}

					successNum := 0
					for tid, resp := range rechargeRespMap {
						_, ok := validErrs[resp.Err]
						So(ok, ShouldBeTrue)

						record, recordErr := s.dao.GetCoinStreamByTid(ctx, tid)
						So(recordErr, ShouldBeNil)

						So(record.DeltaCoinNum, ShouldEqual, num)
						So(record.OpType, ShouldEqual, int32(model.RECHARGETYPE))
						sysCoinType := 0
						if coinType == "gold" {
							if platform == "ios" {
								sysCoinType = 2
							} else {
								sysCoinType = 1
							}
						}
						So(record.CoinType, ShouldEqual, sysCoinType)
						success, queryResp, err := queryQueryWithUid(t, tid, uid)
						So(success, ShouldBeTrue)
						So(err, ShouldBeNil)

						if resp.Err == nil {
							So(record.OpResult, ShouldEqual, model.STREAM_OP_RESULT_ADD_SUCC)
							So(queryResp.Status, ShouldEqual, TX_STATUS_SUCC)

							successNum++
						} else {
							So(record.OpResult, ShouldEqual, model.STREAM_OP_RESULT_ADD_FAILED)
							So(record.OpReason, ShouldEqual, model.STREAM_OP_REASON_LOCK_FAILED)
							So(queryResp.Status, ShouldEqual, TX_STATUS_FAILED)
						}
					}

					So(len(rechargeRespMap), ShouldEqual, times)
					So(successNum, ShouldBeGreaterThan, 0)
					t.Logf("multi successNum:%d", successNum)

					var detail2 *model.DetailWithSnapShot
					tx, _ = s.dao.BeginTx(ctx)

					detail2, _ = s.dao.WalletForUpdate(tx, uid)
					t.Logf("detail2: %+v", detail2)
					So(detail2.SnapShotGold, ShouldEqual, detail1.Gold)
					So(detail2.SnapShotIapGold, ShouldEqual, detail1.IapGold)
					So(detail2.SnapShotSilver, ShouldEqual, detail1.Silver)
					So(detail2.SnapShotTime, ShouldNotEqual, "0001-01-01 00:00:00")

					tx.Rollback()

					afterWallet := getTestWallet(t, uid, platform)
					afterDetail := getTestAll(t, uid, platform)
					successCount := successNum * int(num)

					if coinType == "gold" {
						So(atoiForTest(afterWallet.Gold)-atoiForTest(beforeWallet.Gold), ShouldEqual, successCount)
						So(atoiForTest(afterDetail.Gold)-atoiForTest(beforeDetail.Gold), ShouldEqual, successCount)
						So(atoiForTest(afterDetail.Gold)-atoiForTest(beforeDetail.Gold), ShouldEqual, successCount)
						So(atoiForTest(afterDetail.GoldRechargeCnt)-atoiForTest(beforeDetail.GoldRechargeCnt), ShouldEqual, successCount)
						So(atoiForTest(afterDetail.GoldPayCnt)-atoiForTest(beforeDetail.GoldPayCnt), ShouldEqual, 0)
						So(atoiForTest(afterDetail.SilverPayCnt)-atoiForTest(beforeDetail.SilverPayCnt), ShouldEqual, 0)

					} else if coinType == "silver" {
						So(atoiForTest(afterWallet.Silver)-atoiForTest(beforeWallet.Silver), ShouldEqual, successCount)
						So(atoiForTest(afterDetail.Silver)-atoiForTest(beforeDetail.Silver), ShouldEqual, successCount)
						// silver 不统计充值
						So(atoiForTest(afterDetail.GoldRechargeCnt)-atoiForTest(beforeDetail.GoldRechargeCnt), ShouldEqual, 0)
						So(atoiForTest(afterDetail.GoldPayCnt)-atoiForTest(beforeDetail.GoldPayCnt), ShouldEqual, 0)
						So(atoiForTest(afterDetail.SilverPayCnt)-atoiForTest(beforeDetail.SilverPayCnt), ShouldEqual, 0)
					}

				}

			}
		})
	})

	Convey("iphone", t, func() {
		uid := getTestRandUid()
		platform := "pc"
		coinType := "gold"
		beforeWallet, err := s.dao.Melonseed(ctx, uid)
		beforePcResp := getTestWallet(t, uid, platform)
		beforeIosResp := getTestWallet(t, uid, "ios")
		So(err, ShouldBeNil)
		var num int64 = 1000

		tid := getTestTidForCall(t, int32(model.RECHARGETYPE))

		bp := getTestDefaultBasicParam(tid)

		v, err := s.Recharge(ctx, bp, uid, platform, getTestRechargeOrPayForm(uid, coinType, num, tid))
		So(v, ShouldNotBeNil)
		So(err, ShouldBeNil)

		afterAddPcWallet, err := s.dao.Melonseed(ctx, uid)
		So(err, ShouldBeNil)
		So(beforeWallet.IapGold, ShouldEqual, afterAddPcWallet.IapGold)

		resp := getTestWallet(t, uid, platform)
		So(atoiForTest(resp.Gold)-atoiForTest(beforePcResp.Gold), ShouldEqual, int(num))

		platform = "ios"

		tid = getTestTidForCall(t, int32(model.RECHARGETYPE))
		bp = getTestDefaultBasicParam(tid)
		v, err = s.Recharge(ctx, bp, uid, platform, getTestRechargeOrPayForm(uid, coinType, num, tid))
		So(v, ShouldNotBeNil)
		So(err, ShouldBeNil)
		afterAddIosWallet, err := s.dao.Melonseed(ctx, uid)
		So(err, ShouldBeNil)
		So(beforeWallet.IapGold+num, ShouldEqual, afterAddIosWallet.IapGold) // ios 则加到数据库的iap_gold字段

		resp = getTestWallet(t, uid, platform)

		So(atoiForTest(resp.Gold)-atoiForTest(beforeIosResp.Gold), ShouldEqual, int(num))

	})

}

func Test_TestUser(t *testing.T) {
	Convey("set and reset", t, testWith(func() {
		u := initTestUser()
		platform := "pc"
		resp := getTestWallet(t, u.uid, platform)
		So(u.coin, ShouldEqual, atoiForTest(resp.Gold))
		So(u.coin, ShouldEqual, atoiForTest(resp.Silver))
		platform = "ios"
		resp = getTestWallet(t, u.uid, platform)
		So(u.coin, ShouldEqual, atoiForTest(resp.Gold))
		u.rest()

	}))

	Convey("testWithTestUser", t, testWithTestUser(func(u *TestUser) {
		platform := "pc"
		resp := getTestWallet(t, u.uid, platform)
		So(u.coin, ShouldEqual, atoiForTest(resp.Gold))
		So(u.coin, ShouldEqual, atoiForTest(resp.Silver))
		platform = "ios"
		resp = getTestWallet(t, u.uid, platform)
		So(u.coin, ShouldEqual, atoiForTest(resp.Gold))

		var num int64 = 100
		success, _ := queryRecharge(t, u.uid, platform, "gold", num)
		So(success, ShouldBeTrue)

		resp = getTestWallet(t, u.uid, platform)
		So(u.coin+num, ShouldEqual, atoiForTest(resp.Gold))

	}))
}

func TestService_RechargeDoubleCoin(t *testing.T) {
	times := []string{""}
	times = append(times, time.Now().Add(-time.Second*86400).Format("2006-01-02 15:04:05"))
	for _, htime := range times {
		Convey("default_"+htime, t, testWithTestUser(func(u *TestUser) {
			platform := "pc"
			tid := getTestTidForCall(t, int32(model.RECHARGETYPE))
			s.dao.UpdateSnapShotTime(ctx, u.uid, htime)
			t.Logf("htime:%s", htime)

			bp := getTestDefaultBasicParam(tid)

			v, err := s.Recharge(ctx, bp, u.uid, platform, getTestRechargeOrPayForm(u.uid, "gold", 1, tid))
			So(err, ShouldBeNil)
			So(v, ShouldNotBeNil)

			var detail1 *model.DetailWithSnapShot
			tx, _ := s.dao.BeginTx(ctx)

			detail1, _ = s.dao.WalletForUpdate(tx, u.uid)
			t.Logf("%+v", detail1)
			So(detail1.SnapShotGold, ShouldEqual, u.coin)
			So(detail1.SnapShotIapGold, ShouldEqual, u.coin)
			So(detail1.SnapShotSilver, ShouldEqual, u.coin)
			So(detail1.SnapShotTime, ShouldNotEqual, "")

			tx.Rollback()

			tid = getTestTidForCall(t, int32(model.RECHARGETYPE))

			platform = "ios"
			v, err = s.Recharge(ctx, bp, u.uid, platform, getTestRechargeOrPayForm(u.uid, "gold", 1, tid))
			So(err, ShouldBeNil)
			So(v, ShouldNotBeNil)

			tid = getTestTidForCall(t, int32(model.RECHARGETYPE))
			platform = "android"
			v, err = s.Recharge(ctx, bp, u.uid, platform, getTestRechargeOrPayForm(u.uid, "silver", 1, tid))
			So(err, ShouldBeNil)
			So(v, ShouldNotBeNil)

			var detail2 *model.DetailWithSnapShot
			tx, _ = s.dao.BeginTx(ctx)

			detail2, _ = s.dao.WalletForUpdate(tx, u.uid)
			t.Logf("%+v", detail2)
			So(detail2.SnapShotGold, ShouldEqual, u.coin)
			So(detail2.SnapShotIapGold, ShouldEqual, u.coin)
			So(detail2.SnapShotSilver, ShouldEqual, u.coin)
			So(detail2.SnapShotTime, ShouldEqual, detail1.SnapShotTime)

			tx.Rollback()

		}))
	}
}

func TestService_RechargeForNewUser(t *testing.T) {
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
		v, err := s.Get(ctx, getTestDefaultBasicParam(""), uid, platform, 0)
		resp := v.(*model.MelonseedResp)
		So(resp, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(atoiForTest(resp.Gold), ShouldEqual, 0)
		So(atoiForTest(resp.Silver), ShouldEqual, 0)

		_, err = s.dao.DetailWithoutDefault(ctx, uid)

		So(err, ShouldEqual, sql.ErrNoRows)

		var wg errgroup.Group
		times := 2
		st := 0
		for i := 0; i < times; i++ {
			wg.Go(func() error {
				tid := getTestTidForCall(t, int32(model.RECHARGETYPE))
				bp := getTestDefaultBasicParam(tid)
				_, err := s.Recharge(context.Background(), bp, uid, platform, getTestRechargeOrPayForm(uid, "gold", 1, tid))
				if err == nil {
					st++
				}
				return err
			})
		}
		wg.Wait()
		So(st, ShouldBeGreaterThan, 0)
		t.Logf("st : %d", st)
		var detail2 *model.DetailWithSnapShot
		tx, _ := s.dao.BeginTx(ctx)
		detail2, _ = s.dao.WalletForUpdate(tx, uid)
		t.Logf("detail2:%+v", detail2)
		So(detail2.Gold, ShouldEqual, st)
		So(detail2.SnapShotGold, ShouldEqual, 0)
		So(detail2.SnapShotIapGold, ShouldEqual, 0)
		So(detail2.SnapShotSilver, ShouldEqual, 0)
		So(detail2.SnapShotTime, ShouldNotEqual, "0001-01-01 00:00:00")

		tx.Rollback()

		v, err = s.Get(ctx, getTestDefaultBasicParam(""), uid, platform, 0)
		resp = v.(*model.MelonseedResp)
		So(resp, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(atoiForTest(resp.Gold), ShouldEqual, st)
		So(atoiForTest(resp.Silver), ShouldEqual, 0)

	}))

}

func TestService_RechargeBasicParam(t *testing.T) {

	Convey("default", t, testWithTestUser(func(u *TestUser) {
		for _, platform := range testValidPlatform {
			tid := getTestTidForCall(t, int32(model.RECHARGETYPE))

			bp := getTestDefaultBasicParam(tid)

			v, err := s.Recharge(ctx, bp, u.uid, platform, getTestRechargeOrPayForm(u.uid, "gold", 1, tid))
			So(err, ShouldBeNil)
			So(v, ShouldNotBeNil)

			record, err := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			So(err, ShouldBeNil)
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
			bp.BizSource = "live"
			bp.BizCode = "send-gift"
			bp.Area = 1
			bp.MetaData = "zz"

			v, err := s.Recharge(ctx, bp, u.uid, platform, getTestRechargeOrPayForm(u.uid, "gold", 1, tid))
			So(err, ShouldBeNil)
			So(v, ShouldNotBeNil)

			record, err := s.dao.GetCoinStreamByUidTid(ctx, u.uid, tid)
			So(err, ShouldBeNil)
			oP := model.GetPlatformNo(platform)
			So(record.Platform, ShouldEqual, oP)
			So(record.Reserved1, ShouldEqual, 1)
			So(record.Uid, ShouldEqual, u.uid)
			So(record.CoinType, ShouldEqual, model.GetCoinTypeNumber(model.GetSysCoinType("gold", platform)))
			So(record.OpResult, ShouldEqual, 2)
			So(record.BizSource, ShouldEqual, bp.BizSource)
			So(record.BizCode, ShouldEqual, bp.BizCode)
			So(record.Area, ShouldEqual, bp.Area)
			So(record.MetaData, ShouldEqual, bp.MetaData)
			So(record.DeltaCoinNum, ShouldEqual, 1)

		}

	}))
}
