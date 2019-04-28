package service

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/service/main/vip/conf"
	"go-common/app/service/main/vip/model"
	xtime "go-common/library/time"

	"math/rand"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c = context.TODO()
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") == "" {
		flag.Set("app_id", "main.account.vip-service")
		flag.Set("conf_token", "48c1c43b703d238f0cba6d8d63ec9463")
		flag.Set("tree_id", "10964")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/vip-service-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
	os.Exit(m.Run())
}

func TestService_BatchInfo(t *testing.T) {

	var (
		c            = context.TODO()
		id     int64 = 5
		appkey       = "ad4bb9b8f5d9d4a7"
	)
	Convey("should return true where err == nil and v != nil ", t, func() {
		v, bis, err := s.BatchInfo(c, id, appkey)
		t.Logf("bi:%+v bis:%+v", v, bis)
		So(err, ShouldBeNil)
	})
}

func TestService_OpenCode(t *testing.T) {
	Convey(" open code", t, func() {
		code, err := s.OpenCode(context.TODO(), "d037cfbc8c7d2477", 123)
		So(err, ShouldBeNil)
		t.Logf("code info(%+v)", code)
	})
}

func TestService_WebToken(t *testing.T) {
	Convey("web token", t, func() {
		token, err := s.WebToken(context.TODO())
		So(token, ShouldNotBeNil)
		So(err, ShouldBeNil)

	})
}

func TestService_CodeInfo(t *testing.T) {
	Convey("code info", t, func() {
		code, err := s.CodeInfo(context.TODO(), "8339116653c30d82")
		So(code, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestService_Verify(t *testing.T) {
	Convey("veryfi code", t, func() {
		t, err := s.Verify(context.TODO(), "8339116653c30d82", "f5556e7237330bfd")
		So(t, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestService_BusinessByPool(t *testing.T) {
	Convey("business by pool", t, func() {
		r, err := s.BusinessByPool(context.TODO(), 12)
		So(r, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestService_CodeInfos(t *testing.T) {
	Convey("sel code infos", t, func() {
		var codes = []string{"a482bb5a1d679fe4", "798ae11c6c395893", "0b3dd60f69b9b0d9"}
		cs, err := s.CodeInfos(context.TODO(), codes)
		t.Logf("cs(%+v)", cs)
		So(err, ShouldBeNil)
	})
}

func TestService_Actives(t *testing.T) {
	Convey("sel actives", t, func() {
		var relationIds = []string{"1806141200336386401", "1806141156245707979"}
		res, err := s.Actives(context.TODO(), relationIds)
		t.Logf("res(%+v)", res)
		So(err, ShouldBeNil)
	})
}

func TestService_PayNotify2(t *testing.T) {
	Convey("pay notify2 err == nil", t, func() {
		n := new(model.PayNotifyContent)
		n.PayChannelID = 3
		n.PayChannel = "wechart"
		n.ExpiredTime = time.Now().UnixNano() / 1e6
		n.CustomerID = 10004
		n.TxID = 2997443593637572605
		n.PayStatus = "SUCCESS"
		n.OrderID = "1807021929153562939"
		n.PayAmount = 2500
		err := s.PayNotify2(c, n)
		So(err, ShouldBeNil)
	})
}

func TestService_PaySignNotify(t *testing.T) {
	Convey("pay sign notify err == nil", t, func() {
		n := new(model.PaySignNotify)
		n.CustomerID = 10004
		n.PayChannel = "wechart"
		n.ChangeType = "ADD"
		n.UID = 10
		err := s.PaySignNotify(c, n)
		So(err, ShouldBeNil)
	})
}

func TestService_RefundNotify(t *testing.T) {
	Convey("refund notify ", t, func() {
		arg := new(model.PayRefundNotify)
		arg.CustomerID = 10004
		arg.OrderID = "93035846180822184133"
		arg.RefundCount = 1
		refundList := make([]*model.PayRefundList, 0)
		payRefund := new(model.PayRefundList)
		payRefund.CustomerRefundID = "930358461808221121"
		payRefund.RefundStatus = "REFUND_SUCCESS"
		payRefund.RefundAmount = 6800
		refundList = append(refundList, payRefund)
		arg.BatchRefundList = refundList
		err := s.RefundNotify(c, arg)
		So(err, ShouldBeNil)
	})
}

func TestService_CodeOpened(t *testing.T) {
	Convey("code opened", t, func() {
		arg := new(model.ArgCodeOpened)
		arg.BisAppkey = "b2cf4e9dbe9fd2e3"
		arg.StartTime = xtime.Time(1528041600)
		arg.EndTime = xtime.Time(1535990400)
		arg.Cursor = 0
		arg.BisTs = 1532519139425
		arg.BisSign = "005a96054946166accea285c1c033fbd"
		res, err := s.CodeOpened(c, arg)
		t.Logf("%+v", res)
		So(err, ShouldBeNil)
	})
}

func TestResourceBatchOpenVip(t *testing.T) {
	Convey("test resource batch open vip", t, func() {
		//非大会员->普通大会员
		mid := rand.NewSource(99999).Int63()
		arg := new(model.ArgUseBatch)
		arg.BatchID = 1
		arg.Appkey = "7d9f6f6fe2a898e8"
		arg.Remark = "开通大会员"
		arg.OrderNo = uuid.New().String()[0:24]
		arg.Mid = mid
		err := s.ResourceBatchOpenVip(c, arg)
		So(err, ShouldBeNil)
		//普通大会员->年度大会员
		arg = new(model.ArgUseBatch)
		arg.BatchID = 33
		arg.Appkey = "7d9f6f6fe2a898e8"
		arg.Remark = "开通大会员"
		arg.OrderNo = uuid.New().String()[0:24]
		arg.Mid = mid
		err = s.ResourceBatchOpenVip(c, arg)
		So(err, ShouldBeNil)
		//年度大会员->年度大会员
		arg = new(model.ArgUseBatch)
		arg.BatchID = 34
		arg.Appkey = "7d9f6f6fe2a898e8"
		arg.Remark = "开通大会员"
		arg.OrderNo = uuid.New().String()[0:24]
		arg.Mid = mid
		err = s.ResourceBatchOpenVip(c, arg)
		So(err, ShouldBeNil)
	})
}
