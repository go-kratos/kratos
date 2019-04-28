package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoJuryApply(t *testing.T) {
	convey.Convey("JuryApply", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			expired = time.Now()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.JuryApply(c, mid, expired)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddUserVoteTotal(t *testing.T) {
	convey.Convey("AddUserVoteTotal", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddUserVoteTotal(c, mid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoJuryInfo(t *testing.T) {
	convey.Convey("JuryInfo", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.JuryInfo(c, mid)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoJuryInfos(t *testing.T) {
	convey.Convey("JuryInfos", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{111, 88889017}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			mbj, err := d.JuryInfos(c, mids)
			convCtx.Convey("Then err should be nil.mbj should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(mbj, convey.ShouldNotBeNil)
			})
		})
	})
}
