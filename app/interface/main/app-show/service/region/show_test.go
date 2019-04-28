package region

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model"

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
	dir, _ := filepath.Abs("../../cmd/app-show-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestChildShow(t *testing.T) {
	Convey("get ChildShow data", t, WithService(func(s *Service) {
		res := s.ChildShow(context.TODO(), model.PlatIPhone, 0, 1, 0, 111, "", "iphone", time.Now())
		So(res, ShouldNotBeEmpty)
	}))
}

func TestShow(t *testing.T) {
	Convey("get Show data", t, WithService(func(s *Service) {
		res := s.Show(context.TODO(), model.PlatIPhone, 1, 0, 0, "channel", "", "", "iphone", "phone", "")
		So(res, ShouldNotBeEmpty)
	}))
}

func TestShowDynamic(t *testing.T) {
	Convey("get ShowDynamic data", t, WithService(func(s *Service) {
		res := s.ShowDynamic(context.TODO(), model.PlatIPhone, 0, 1, 1, 20)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestChildListShow(t *testing.T) {
	Convey("get ChildListShow data", t, WithService(func(s *Service) {
		res := s.ChildListShow(context.TODO(), model.PlatIPhone, 1, 0, 1, 20, 0, 0, "new", "ios", "iphone", "phone")
		So(res, ShouldNotBeEmpty)
	}))
}

func TestDynamic(t *testing.T) {
	Convey("get Dynamic data", t, WithService(func(s *Service) {
		res := s.Dynamic(context.TODO(), model.PlatIPhone, 1, 0, 0, "channel", "", "", "iphone", "phone", "", time.Now())
		So(res, ShouldNotBeEmpty)
	}))
}

func TestDynamicList(t *testing.T) {
	Convey("get DynamicList data", t, WithService(func(s *Service) {
		res := s.DynamicList(context.TODO(), model.PlatIPhone, 1, true, 0, 0, time.Now())
		So(res, ShouldNotBeEmpty)
	}))
}

func TestDynamicChild(t *testing.T) {
	Convey("get DynamicChild data", t, WithService(func(s *Service) {
		res := s.DynamicChild(context.TODO(), model.PlatIPhone, 1, 0, 0, 0, "iphone", time.Now())
		So(res, ShouldNotBeEmpty)
	}))
}
func TestDynamicListChild(t *testing.T) {
	Convey("get DynamicListChild data", t, WithService(func(s *Service) {
		res := s.DynamicListChild(context.TODO(), model.PlatIPhone, 1, 0, 0, true, 0, 0, time.Now())
		So(res, ShouldNotBeEmpty)
	}))
}

func TestAudit(t *testing.T) {
	Convey("get Audit data", t, WithService(func(s *Service) {
		res, _ := s.Audit(context.TODO(), "iphone", model.PlatIPhone, 1, 1, false)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestAuditChild(t *testing.T) {
	Convey("get AuditChild data", t, WithService(func(s *Service) {
		res, _ := s.AuditChild(context.TODO(), "iphone", "", model.PlatIPhone, 1, 1, 0)
		So(res, ShouldNotBeEmpty)
	}))
}

func TestAuditChildList(t *testing.T) {
	Convey("get AuditChildList data", t, WithService(func(s *Service) {
		res, _ := s.AuditChildList(context.TODO(), "iphone", "", model.PlatIPhone, 1, 1, 0, 1, 1)
		So(res, ShouldNotBeEmpty)
	}))
}
