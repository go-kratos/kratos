package upcrmservice

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmserviceAddAuditLog(t *testing.T) {
	convey.Convey("AddAuditLog", t, func(ctx convey.C) {
		var (
			bizID   = int(0)
			tp      = int8(0)
			action  = ""
			uid     = int64(0)
			uname   = ""
			oids    = []int64{}
			index   = []interface{}{}
			content map[string]interface{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.AddAuditLog(bizID, tp, action, uid, uname, oids, index, content)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
