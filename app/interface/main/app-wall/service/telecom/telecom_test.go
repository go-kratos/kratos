package telecom

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-wall/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func WithService(f func(s *Service)) func() {
	return func() {
		f(s)
	}
}

func init() {
	dir, _ := filepath.Abs("../../cmd/app-wall-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestTelecomPay(t *testing.T) {
	Convey("TelecomPay", t, WithService(func(s *Service) {
		res, _, err := s.TelecomPay(context.TODO(), 1, 1, 1, 1, 1, "")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestCancelRepeatOrder(t *testing.T) {
	Convey("CancelRepeatOrder", t, WithService(func(s *Service) {
		res, err := s.CancelRepeatOrder(context.TODO(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestOrderList(t *testing.T) {
	Convey("OrderList", t, WithService(func(s *Service) {
		res, _, err := s.OrderList(context.TODO(), 1, 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestPhoneFlow(t *testing.T) {
	Convey("PhoneFlow", t, WithService(func(s *Service) {
		res, _, err := s.PhoneFlow(context.TODO(), 1, 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestOrderConsent(t *testing.T) {
	Convey("OrderConsent", t, WithService(func(s *Service) {
		res, _, err := s.OrderConsent(context.TODO(), 1, 1, "")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestPhoneCode(t *testing.T) {
	Convey("PhoneCode", t, WithService(func(s *Service) {
		res, _, err := s.PhoneCode(context.TODO(), 1, "", time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestPhoneSendSMS(t *testing.T) {
	Convey("PhoneSendSMS", t, WithService(func(s *Service) {
		err := s.PhoneSendSMS(context.TODO(), 1)
		So(err, ShouldBeNil)
	}))
}

func TestOrderState(t *testing.T) {
	Convey("OrderState", t, WithService(func(s *Service) {
		res, _, err := s.OrderState(context.TODO(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}
