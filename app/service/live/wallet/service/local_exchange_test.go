package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"testing"
	"time"
)

func getWalletService() *WalletService {
	ws := new(WalletService)
	ws.c = ctx
	ws.s = s
	return ws
}

func TestPayHandler_LocalPay(t *testing.T) {
	Convey("normal", t, testWithTestUser(func(u *TestUser) {

		nd, _ := s.dao.Detail(ctx, u.uid)

		h := getLocalPayHandler(getWalletService())
		s.dao.UpdateSnapShotTime(ctx, u.uid, "")
		ef := model.ExchangeParam{
			Uid:           u.uid,
			SrcCoinType:   "gold",
			SrcCoinNum:    1,
			DestCoinType:  "silver",
			DestCoinNum:   10,
			ExtendTid:     getTestExtendTid(),
			TransactionId: getTestTidForCall(t, 2),
			Timestamp:     time.Now().Unix(),
			Platform:      "pc",
			Reason:        1,
		}
		t.Logf("uid:%d tid:%s", u.uid, ef.TransactionId)

		srcCoinStream := model.CoinStreamRecord{}
		model.InjectFieldToCoinStream(&srcCoinStream, &ef)
		srcCoinStream.DeltaCoinNum = -1
		srcCoinStream.CoinType = model.SysCoinTypeGold
		srcCoinStream.OpType = int32(model.EXCHANGETYPE)

		destCoinStream := model.CoinStreamRecord{}
		model.InjectFieldToCoinStream(&destCoinStream, &ef)
		destCoinStream.DeltaCoinNum = 10
		destCoinStream.CoinType = model.SysCoinTypeSilver
		destCoinStream.OpType = int32(model.EXCHANGETYPE)

		exchangeRecord := model.NewExchangeSteam(ef.Uid, ef.TransactionId, 1, int32(ef.SrcCoinNum), 0, int32(ef.DestCoinNum), ef.Timestamp, 0)

		resp, err := h.exchange(ctx, u.uid, ef.Platform, model.SysCoinTypeGold, ef.SrcCoinNum, model.SysCoinTypeSilver, 10, &srcCoinStream, &destCoinStream, exchangeRecord)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		t.Logf("%+v", resp)

		m1 := getTestWallet(t, u.uid, "pc")
		So(m1.Gold, ShouldEqual, "999")
		So(m1.Silver, ShouldEqual, "1010")

		var detail2 *model.DetailWithSnapShot
		tx, _ := s.dao.BeginTx(ctx)

		detail2, _ = s.dao.WalletForUpdate(tx, u.uid)
		t.Logf("detail2: %+v", detail2)
		So(detail2.SnapShotGold, ShouldEqual, 1000)
		So(detail2.SnapShotIapGold, ShouldEqual, 1000)
		So(detail2.SnapShotSilver, ShouldEqual, 1000)
		So(detail2.SnapShotTime, ShouldNotEqual, "0001-01-01 00:00:00")
		So(detail2.GoldPayCnt-nd.GoldPayCnt, ShouldEqual, 1)
		So(detail2.SilverPayCnt-nd.SilverPayCnt, ShouldEqual, 0)

		tx.Rollback()

		ef.TransactionId = getTestTidForCall(t, 2)
		model.InjectFieldToCoinStream(&srcCoinStream, &ef)
		srcCoinStream.DeltaCoinNum = -1
		srcCoinStream.CoinType = model.SysCoinTypeGold
		srcCoinStream.OpType = int32(model.EXCHANGETYPE)
		model.InjectFieldToCoinStream(&destCoinStream, &ef)
		destCoinStream.DeltaCoinNum = 10
		destCoinStream.CoinType = model.SysCoinTypeSilver
		destCoinStream.OpType = int32(model.EXCHANGETYPE)
		exchangeRecord = model.NewExchangeSteam(ef.Uid, ef.TransactionId, 1, int32(ef.SrcCoinNum), 0, int32(ef.DestCoinNum), ef.Timestamp, 0)
		resp, err = h.exchange(ctx, u.uid, ef.Platform, model.SysCoinTypeGold, ef.SrcCoinNum, model.SysCoinTypeSilver, 10, &srcCoinStream, &destCoinStream, exchangeRecord)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		t.Logf("%+v", resp)

		d1, _ := s.dao.Detail(ctx, u.uid)
		So(d1.Gold, ShouldEqual, 998)
		So(d1.Silver, ShouldEqual, 1020)

		tx, _ = s.dao.BeginTx(ctx)

		detail2, _ = s.dao.WalletForUpdate(tx, u.uid)
		t.Logf("detail2: %+v", detail2)
		So(detail2.SnapShotGold, ShouldEqual, 1000)
		So(detail2.SnapShotIapGold, ShouldEqual, 1000)
		So(detail2.SnapShotSilver, ShouldEqual, 1000)
		So(detail2.SnapShotTime, ShouldNotEqual, "0001-01-01 00:00:00")

		tx.Rollback()

		s.dao.UpdateSnapShotTime(ctx, u.uid, "")

		ef.TransactionId = getTestTidForCall(t, 2)
		model.InjectFieldToCoinStream(&srcCoinStream, &ef)
		srcCoinStream.DeltaCoinNum = -1
		srcCoinStream.CoinType = model.SysCoinTypeGold
		srcCoinStream.OpType = int32(model.EXCHANGETYPE)
		model.InjectFieldToCoinStream(&destCoinStream, &ef)
		destCoinStream.DeltaCoinNum = 10
		destCoinStream.CoinType = model.SysCoinTypeSilver
		destCoinStream.OpType = int32(model.EXCHANGETYPE)
		exchangeRecord = model.NewExchangeSteam(ef.Uid, ef.TransactionId, 1, int32(ef.SrcCoinNum), 0, int32(ef.DestCoinNum), ef.Timestamp, 0)
		resp, err = h.exchange(ctx, u.uid, ef.Platform, model.SysCoinTypeGold, ef.SrcCoinNum, model.SysCoinTypeSilver, 10, &srcCoinStream, &destCoinStream, exchangeRecord)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		t.Logf("%+v", resp)

		d1, _ = s.dao.Detail(ctx, u.uid)
		So(d1.Gold, ShouldEqual, 997)
		So(d1.Silver, ShouldEqual, 1030)

		tx, _ = s.dao.BeginTx(ctx)

		detail2, _ = s.dao.WalletForUpdate(tx, u.uid)
		t.Logf("detail2: %+v", detail2)
		So(detail2.SnapShotGold, ShouldEqual, 998)
		So(detail2.SnapShotIapGold, ShouldEqual, 1000)
		So(detail2.SnapShotSilver, ShouldEqual, 1020)
		So(detail2.SnapShotTime, ShouldNotEqual, "0001-01-01 00:00:00")
		tx.Rollback()

	}))

}
