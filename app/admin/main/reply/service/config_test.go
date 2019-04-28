package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/reply/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	var (
		m = &model.Config{
			Oid:       1,
			Type:      4,
			ShowEntry: 1,
			ShowAdmin: 1,
		}
		c = context.Background()
	)
	Convey("test config ", t, WithService(func(s *Service) {
		_, err := s.AddReplyConfig(c, m)
		So(err, ShouldBeNil)
		cfg, err := s.LoadReplyConfig(c, m.Type, m.Category, m.Oid)
		So(err, ShouldBeNil)
		So(cfg, ShouldNotBeNil)
		list, total, pages, err := s.PaginateReplyConfig(c, m.Type, m.Category, m.Oid, "", 0, 10)
		So(err, ShouldBeNil)
		So(len(list), ShouldNotEqual, 0)
		So(total, ShouldNotEqual, 0)
		So(pages, ShouldNotEqual, 0)
		ok, err := s.RenewReplyConfig(c, cfg.ID)
		So(err, ShouldBeNil)
		So(ok, ShouldEqual, true)
	}))
}
