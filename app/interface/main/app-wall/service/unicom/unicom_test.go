package unicom

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

func TestUserFlow(t *testing.T) {
	Convey("Unicom UserFlow", t, WithService(func(s *Service) {
		res, _, err := s.UserFlow(context.TODO(), "", "iphone", "127.0.0.1", 111, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestUserState(t *testing.T) {
	Convey("Unicom UserState", t, WithService(func(s *Service) {
		res, _, err := s.UserState(context.TODO(), "", "iphone", "127.0.0.1", 111, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestUnicomState(t *testing.T) {
	Convey("Unicom UnicomState", t, WithService(func(s *Service) {
		res, err := s.UnicomState(context.TODO(), "", "iphone", "127.0.0.1", 111, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestUserFlowState(t *testing.T) {
	Convey("Unicom UserFlowState", t, WithService(func(s *Service) {
		res, err := s.UserFlowState(context.TODO(), "", time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestIsUnciomIP(t *testing.T) {
	Convey("Unicom IsUnciomIP", t, WithService(func(s *Service) {
		err := s.IsUnciomIP(0, "127.0.0.1", "iphone", 0, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestUserUnciomIP(t *testing.T) {
	Convey("Unicom UserUnciomIP", t, WithService(func(s *Service) {
		err := s.UserUnciomIP(0, "127.0.0.1", "", "iphone", 0, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestOrder(t *testing.T) {
	Convey("Unicom Order", t, WithService(func(s *Service) {
		res, _, err := s.Order(context.TODO(), "", "", 0, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestCancelOrder(t *testing.T) {
	Convey("Unicom CancelOrder", t, WithService(func(s *Service) {
		res, _, err := s.CancelOrder(context.TODO(), "")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestUnicomSMSCode(t *testing.T) {
	Convey("Unicom UnicomSMSCode", t, WithService(func(s *Service) {
		res, err := s.UnicomSMSCode(context.TODO(), "", time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestAddUnicomBind(t *testing.T) {
	Convey("Unicom AddUnicomBind", t, WithService(func(s *Service) {
		res, err := s.AddUnicomBind(context.TODO(), "", 1, 1, time.Now())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestReleaseUnicomBind(t *testing.T) {
	Convey("Unicom ReleaseUnicomBind", t, WithService(func(s *Service) {
		res, err := s.ReleaseUnicomBind(context.TODO(), 1, 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestUserBind(t *testing.T) {
	Convey("Unicom UserBind", t, WithService(func(s *Service) {
		res, _, err := s.UserBind(context.TODO(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestUnicomPackList(t *testing.T) {
	Convey("Unicom UnicomPackList", t, WithService(func(s *Service) {
		res := s.UnicomPackList()
		So(res, ShouldNotBeEmpty)
	}))
}

func TestUnicomPackReceive(t *testing.T) {
	Convey("Unicom UnicomPackReceive", t, WithService(func(s *Service) {
		res, _ := s.UnicomPackReceive(context.TODO(), 1, 1, time.Now())
		So(res, ShouldNotBeEmpty)
	}))
}

func TestUnicomUnicomFlowPack(t *testing.T) {
	Convey("Unicom UnicomFlowPack", t, WithService(func(s *Service) {
		res, _ := s.UnicomFlowPack(context.TODO(), 1, "", time.Now())
		So(res, ShouldNotBeEmpty)
	}))
}

func TestUnicomloadUnicomFlow(t *testing.T) {
	Convey("Unicom loadUnicomFlow", t, WithService(func(s *Service) {
		s.loadUnicomFlow()
	}))
}

func TestUnicomloadUnicomIPOrder(t *testing.T) {
	Convey("Unicom loadUnicomIPOrder", t, WithService(func(s *Service) {
		s.loadUnicomIPOrder(time.Now())
	}))
}
