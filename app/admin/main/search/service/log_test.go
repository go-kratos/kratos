package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/search/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_LogAudit(t *testing.T) {
	var (
		err error
		c   = context.Background()
		p   = &model.LogParams{
			Bsp: &model.BasicSearchParams{
				AppID: "log_audit",
			},
		}
		params map[string][]string
	)

	Convey("LogAudit", t, WithService(func(s *Service) {
		business, ok := svr.Check("log_audit", 0)
		if !ok {
			return
		}
		_, err = s.LogAudit(c, params, p, business)
		So(err, ShouldBeNil)
	}))
}

//func Test_LogAuditGroupBy(t *testing.T) {
//	var (
//		err error
//		c   = context.Background()
//		p   = &model.LogParams{
//			Bsp: &model.BasicSearchParams{
//				AppID: "log_audit_group",
//			},
//		}
//		params map[string][]string
//	)
//	params = map[string][]string{
//		"group": {"oid"},
//	}
//	Convey("LogAuditGroupBy", t, WithService(func(s *Service) {
//		indexMapping, indexFmt, ok := svr.Check("log_audit", p.Business)
//		if !ok {
//			return
//		}
//		_, err = s.LogAuditGroupBy(c, params, p, indexMapping, indexFmt)
//		Printf("---------%v", err)
//		So(err, ShouldBeNil)
//	}))
//}

func Test_LogUserAction(t *testing.T) {
	var (
		err error
		c   = context.Background()
		p   = &model.LogParams{
			Bsp: &model.BasicSearchParams{
				AppID: "log_user_action",
			},
		}
		params map[string][]string
	)
	Convey("LogUserAction", t, WithService(func(s *Service) {
		business, ok := svr.Check("log_user_action", 0)
		if !ok {
			return
		}
		_, err = s.LogUserAction(c, params, p, business)
		So(err, ShouldBeNil)
	}))
}
