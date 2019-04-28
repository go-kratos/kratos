package archive

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveStaff(t *testing.T) {
	convey.Convey("Staff", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(27515258)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			aids, err := d.Staff(c, mid)
			ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(aids), convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestArchiveStaffs(t *testing.T) {
	convey.Convey("Staffs", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{27515258}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			aidm, err := d.Staffs(c, mids)
			ctx.Convey("Then err should be nil.aidm should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(aidm, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestArchiveStaffAid(t *testing.T) {
	convey.Convey("StaffAid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(10110188)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mids, err := d.StaffAid(c, aid)
			ctx.Convey("Then err should be nil.mids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mids, convey.ShouldNotBeNil)
			})
		})
	})
}
