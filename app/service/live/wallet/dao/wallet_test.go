package dao

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"go-common/app/service/live/wallet/conf"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"math/rand"
)

var (
	once    sync.Once
	d       *Dao
	ctx     = context.Background()
	testUID int64
	r       *rand.Rand
)

func initConf() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
}

func TestMain(m *testing.M) {
	conf.ConfPath = "../cmd/live-wallet-test.toml"
	once.Do(startService)
	os.Exit(m.Run())
}

func startService() {
	initConf()
	d = New(conf.Conf)
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	testUID = r.Int63n(100000000)
	log.Info("test uid:%ld", testUID)
	time.Sleep(time.Second * 1)
}

func TestInitWallet(t *testing.T) {
	Convey("Init Wallet", t, func() {
		once.Do(startService)
		_, err := d.InitWallet(ctx, testUID, 0, 0, 0)
		So(err, ShouldBeNil)
	})
}

func TestMelonseed(t *testing.T) {
	Convey("Melon seed", t, func() {
		once.Do(startService)
		var wallet *model.Melonseed
		wallet, err := d.Melonseed(ctx, testUID)
		So(err, ShouldBeNil)
		So(wallet.Gold, ShouldEqual, 0)
	})
}

func TestDetail(t *testing.T) {
	Convey("Detail", t, func() {
		once.Do(startService)
		var detail *model.Detail
		detail, err := d.Detail(ctx, testUID)
		So(err, ShouldBeNil)
		So(detail.Gold, ShouldEqual, 0)
		So(detail.Silver, ShouldEqual, 0)
		So(detail.SilverPayCnt, ShouldEqual, 0)
		So(detail.IapGold, ShouldEqual, 0)
		So(detail.GoldRechargeCnt, ShouldEqual, 0)
		So(detail.GoldPayCnt, ShouldEqual, 0)

	})
}

func TestDao_AddGold(t *testing.T) {
	Convey("Add Gold", t, func() {
		once.Do(startService)
		var testUID = r.Int63n(100000000)
		var odetail *model.Detail
		odetail, err := d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		num := 1000
		affect, err := d.AddGold(ctx, testUID, num)
		So(err, ShouldBeNil)
		So(affect, ShouldBeGreaterThanOrEqualTo, 1) // 插入affect=1 更新affect=2

		var ndetail *model.Detail
		ndetail, err = d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		So(ndetail.Gold-odetail.Gold, ShouldEqual, num)
		So(ndetail.GoldRechargeCnt, ShouldEqual, odetail.GoldRechargeCnt)
	})
}
func TestDao_RechargeGold(t *testing.T) {
	Convey("recharge gold", t, func() {
		once.Do(startService)
		var testUID = r.Int63n(100000000)
		var odetail *model.Detail
		odetail, err := d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		num := 1000
		affect, err := d.RechargeGold(ctx, testUID, num)
		So(err, ShouldBeNil)
		So(affect, ShouldBeGreaterThanOrEqualTo, 1) // 插入affect=1 更新affect=2

		var ndetail *model.Detail
		ndetail, err = d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		So(ndetail.Gold-odetail.Gold, ShouldEqual, num)
		So(ndetail.IapGold-odetail.IapGold, ShouldEqual, 0)
		So(ndetail.GoldRechargeCnt-odetail.GoldRechargeCnt, ShouldEqual, num)
	})
}

func TestDao_RechargeIapGold(t *testing.T) {
	Convey("recharge iap gold", t, func() {
		once.Do(startService)
		var testUID = r.Int63n(100000000)
		var odetail *model.Detail
		odetail, err := d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		num := 1000
		affect, err := d.RechargeIapGold(ctx, testUID, num)
		So(err, ShouldBeNil)
		So(affect, ShouldBeGreaterThanOrEqualTo, 1) // 插入affect=1 更新affect=2

		var ndetail *model.Detail
		ndetail, err = d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		So(ndetail.Gold-odetail.Gold, ShouldEqual, 0)
		So(ndetail.IapGold-odetail.IapGold, ShouldEqual, num)
		So(ndetail.GoldRechargeCnt-odetail.GoldRechargeCnt, ShouldEqual, num)
	})
}

func TestDao_AddIapGold(t *testing.T) {
	Convey("Add iap Gold", t, func() {
		once.Do(startService)
		var testUID = r.Int63n(100000000)
		var odetail *model.Detail
		odetail, err := d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		num := 1000
		affect, err := d.AddIapGold(ctx, testUID, num)
		So(err, ShouldBeNil)
		So(affect, ShouldBeGreaterThanOrEqualTo, 1) // 插入affect=1 更新affect=2

		var ndetail *model.Detail
		ndetail, err = d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		So(ndetail.Gold-odetail.Gold, ShouldEqual, 0)
		So(ndetail.IapGold-odetail.IapGold, ShouldEqual, num)
		So(ndetail.GoldRechargeCnt, ShouldEqual, odetail.GoldRechargeCnt)
	})
}

func TestDao_AddSilver(t *testing.T) {
	Convey("Add silver", t, func() {
		once.Do(startService)
		var testUID = r.Int63n(100000000)
		var odetail *model.Detail
		odetail, err := d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		num := 1000
		affect, err := d.AddSilver(ctx, testUID, num)
		So(err, ShouldBeNil)
		So(affect, ShouldBeGreaterThanOrEqualTo, 1) // 插入affect=1 更新affect=2

		var ndetail *model.Detail
		ndetail, err = d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		So(ndetail.Gold-odetail.Gold, ShouldEqual, 0)
		So(ndetail.IapGold-odetail.IapGold, ShouldEqual, 0)
		So(ndetail.Silver-odetail.Silver, ShouldEqual, num)
	})
}

func TestDao_ConsumeGold(t *testing.T) {
	Convey("Consume Gold", t, func() {
		once.Do(startService)
		var testUID = r.Int63n(100000000)
		var odetail *model.Detail
		odetail, err := d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		num := 1000
		affect, err := d.RechargeGold(ctx, testUID, num)
		So(err, ShouldBeNil)
		So(affect, ShouldBeGreaterThanOrEqualTo, 1) // 插入affect=1 更新affect=2

		pay := 500
		affect, err = d.ConsumeGold(ctx, testUID, pay)
		So(err, ShouldBeNil)
		So(affect, ShouldEqual, 1)

		var ndetail *model.Detail
		ndetail, err = d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		So(ndetail.Gold-odetail.Gold, ShouldEqual, num-pay)
		So(ndetail.IapGold-odetail.IapGold, ShouldEqual, 0)
		So(ndetail.GoldRechargeCnt-odetail.GoldRechargeCnt, ShouldEqual, num)
		So(ndetail.GoldPayCnt-odetail.GoldPayCnt, ShouldEqual, pay)
	})
}

func TestDao_ConsumeIapGold(t *testing.T) {
	Convey("Consume Iap Gold", t, func() {
		once.Do(startService)
		var testUID = r.Int63n(100000000)
		var odetail *model.Detail
		odetail, err := d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		num := 1000
		affect, err := d.RechargeIapGold(ctx, testUID, num)
		So(err, ShouldBeNil)
		So(affect, ShouldBeGreaterThanOrEqualTo, 1) // 插入affect=1 更新affect=2

		pay := 500
		affect, err = d.ConsumeIapGold(ctx, testUID, pay)
		So(err, ShouldBeNil)
		So(affect, ShouldEqual, 1)

		var ndetail *model.Detail
		ndetail, err = d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		So(ndetail.IapGold-odetail.IapGold, ShouldEqual, num-pay)
		So(ndetail.Gold-odetail.Gold, ShouldEqual, 0)
		So(ndetail.GoldRechargeCnt-odetail.GoldRechargeCnt, ShouldEqual, num)
		So(ndetail.GoldPayCnt-odetail.GoldPayCnt, ShouldEqual, pay)
	})
}

func TestDao_ConsumeSilver(t *testing.T) {
	Convey("Consume Silver", t, func() {
		once.Do(startService)
		var testUID = r.Int63n(100000000)
		var odetail *model.Detail
		odetail, err := d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		num := 1000
		affect, err := d.AddSilver(ctx, testUID, num)
		So(err, ShouldBeNil)
		So(affect, ShouldBeGreaterThanOrEqualTo, 1) // 插入affect=1 更新affect=2

		pay := 500
		affect, err = d.ConsumeSilver(ctx, testUID, pay)
		So(err, ShouldBeNil)
		So(affect, ShouldEqual, 1)

		var ndetail *model.Detail
		ndetail, err = d.Detail(ctx, testUID)
		So(err, ShouldBeNil)

		So(ndetail.IapGold-odetail.IapGold, ShouldEqual, 0)
		So(ndetail.Gold-odetail.Gold, ShouldEqual, 0)
		So(ndetail.Silver-odetail.Silver, ShouldEqual, num-pay)
		So(ndetail.GoldRechargeCnt-odetail.GoldRechargeCnt, ShouldEqual, 0)
		So(ndetail.GoldPayCnt-odetail.GoldPayCnt, ShouldEqual, 0)
		So(ndetail.SilverPayCnt-odetail.SilverPayCnt, ShouldEqual, pay)
	})
}
