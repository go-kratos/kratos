package audit

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/model/audit"

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
	dir, _ := filepath.Abs("../../cmd/app-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func TestAudits(t *testing.T) {
	Convey("get Audits", t, WithService(func(s *Service) {
		res, err := s.Audits(context.TODO())
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestAuditByID(t *testing.T) {
	Convey("get AuditByID", t, WithService(func(s *Service) {
		res, err := s.AuditByID(context.TODO(), 1)
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
	}))
}

func TestAddAudit(t *testing.T) {
	Convey("add Audit", t, WithService(func(s *Service) {
		a := &audit.Param{
			Build:   222,
			MobiApp: "iphone",
		}
		err := s.AddAudit(context.TODO(), a, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestUpdateAudit(t *testing.T) {
	Convey("update Audit", t, WithService(func(s *Service) {
		a := &audit.Param{
			ID:      1,
			Build:   222,
			Remark:  "asdsd",
			MobiApp: "iphone",
		}
		err := s.UpdateAudit(context.TODO(), a, time.Now())
		So(err, ShouldBeNil)
	}))
}

func TestDelAudit(t *testing.T) {
	Convey("del Audit", t, WithService(func(s *Service) {
		err := s.DelAudit(context.TODO(), 19)
		So(err, ShouldBeNil)
	}))
}
