package service

import (
	"context"
	"flag"
	"testing"
	"time"

	"go-common/app/job/main/vip/conf"
	"go-common/app/job/main/vip/model"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c context.Context
)

func init() {
	flag.Set("conf", "../cmd/vip-job-test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	c = context.TODO()
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	s = New(conf.Conf)
	time.Sleep(time.Second * 2)
}

func Test_ScanUserInfo(t *testing.T) {
	Convey("should return true err == nil", t, func() {
		err := s.ScanUserInfo(context.TODO())
		So(err, ShouldBeNil)
	})
}

func TestService_HadExpiredMsgJob(t *testing.T) {
	Convey("had expireMsg job", t, func() {
		s.hadExpiredMsgJob()
	})
}

func TestService_WillExpiredMsgJob(t *testing.T) {
	Convey("had expire msg job", t, func() {
		s.willExpiredMsgJob()
	})
}

func TestService_SendMessageJob(t *testing.T) {
	Convey("send message job", t, func() {
		s.sendMessageJob()
	})
}
func TestService_SendBcoinJob(t *testing.T) {
	Convey("send bcoin job", t, func() {
		s.sendBcoinJob()
	})
}
func TestSalaryVideoCouponJob(t *testing.T) {
	Convey("salaryVideoCouponJob err == nil", t, func() {
		s.salaryVideoCouponJob()
		s.salaryVideoCouponJob()
	})
}

func TestService_HandlerVipChangeHistory(t *testing.T) {
	Convey("handlervip change history ", t, func() {
		err := s.HandlerVipChangeHistory()
		So(err, ShouldBeNil)
	})
}

func TestService_HandlerBcoin(t *testing.T) {
	Convey(" handler bcoin history ", t, func() {
		err := s.HandlerBcoin()
		So(err, ShouldBeNil)
	})
}

func TestService_HandlerPayOrder(t *testing.T) {
	Convey("handler pay order", t, func() {
		err := s.HandlerPayOrder()
		So(err, ShouldBeNil)
	})
}

func Test_push(t *testing.T) {
	Convey("handler push data err should be nil", t, func() {

		err := s.pushData(context.TODO())
		So(err, ShouldBeNil)
	})
}
func TestService_CheckBcoinData(t *testing.T) {
	Convey("check bcoin data", t, func() {
		mids, err := s.CheckBcoinData(context.TODO())
		So(mids, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestService_CheckChangeHistory(t *testing.T) {
	Convey("check change history", t, func() {
		mids, err := s.CheckChangeHistory(context.TODO())
		So(mids, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})
}
func Test_HandlerAutoRenewLogInfo(t *testing.T) {
	Convey("err should be nil", t, func() {
		err := s.handlerAutoRenewLogInfo(context.TODO(), &model.VipUserInfo{Mid: 2089809, PayType: model.AutoRenew, PayChannelID: 100})
		So(err, ShouldBeNil)
	})
}
