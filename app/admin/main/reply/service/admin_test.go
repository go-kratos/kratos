package service

import (
	"context"
	"go-common/app/admin/main/reply/model"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdminLog(t *testing.T) {
	var (
		ok      bool
		oid     = int64(1)
		rpID    = int64(1)
		adminID = int64(1)
		typ     = int32(4)
		now     = time.Now()
		c       = context.Background()
	)
	Convey("action set", t, WithService(func(s *Service) {
		err := s.addAdminLog(c, oid, rpID, adminID, typ, model.AdminIsNew, model.AdminIsReport, model.AdminOperDelete, "test", "remark", now)
		So(err, ShouldBeNil)
		list, err := s.LogsByRpID(c, rpID)
		So(err, ShouldBeNil)
		So(len(list), ShouldNotEqual, 0)
		for _, log := range list {
			if log.ReplyID == rpID {
				ok = true
			}
		}
		So(ok, ShouldEqual, true)
	}))

}
