package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/sms/conf"
	pb "go-common/app/service/main/sms/api"

	. "github.com/smartystreets/goconvey/convey"
)

var s *Service

func init() {
	dir, _ := filepath.Abs("../cmd/sms-admin-test.toml")
	flag.Set("conf", dir)
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}

func TestAddTemplate(t *testing.T) {
	Convey("add tpl", t, func() {
		req := &pb.AddTemplateReq{
			Tcode:     "test",
			Template:  "test",
			Stype:     2,
			Submitter: "wj",
		}
		_, err := s.AddTemplate(context.TODO(), req)
		t.Log(err)
	})
}

func TestUpdateTemplate(t *testing.T) {
	Convey("update tpl", t, func() {
		req := &pb.UpdateTemplateReq{
			Tcode:     "test",
			Template:  "test",
			Stype:     2,
			Status:    1,
			Submitter: "wj",
		}
		_, err := s.UpdateTemplate(context.TODO(), req)
		So(err, ShouldNotBeNil)
	})
}

func TestTemplateList(t *testing.T) {
	Convey("tpl list", t, func() {
		req := &pb.TemplateListReq{Pn: 1, Ps: 10}
		res, err := s.TemplateList(context.TODO(), req)
		So(err, ShouldBeNil)
		So(res.Total, ShouldBeGreaterThan, 0)
		So(res.List, ShouldNotBeEmpty)
	})
}
