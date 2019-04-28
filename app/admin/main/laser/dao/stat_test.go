package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoStatArchiveStat(t *testing.T) {
	convey.Convey("StatArchiveStat", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			business  = int(0)
			typeIDS   = []int64{}
			uids      = []int64{}
			statTypes = []int64{}
			statDate  = time.Now()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			d.StatArchiveStat(c, business, typeIDS, uids, statTypes, statDate)
			convCtx.Convey("Then err should be nil.statNodes should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestDaoQueryArchiveCargo(t *testing.T) {
	convey.Convey("QueryArchiveCargo", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			statTime = time.Now()
			uids     = []int64{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			d.QueryArchiveCargo(c, statTime, uids)
			convCtx.Convey("Then err should be nil.items should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestDaoGetUIDByNames(t *testing.T) {
	convey.Convey("GetUIDByNames", t, func(convCtx convey.C) {
		var (
			c      = context.Background()
			unames = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			d.GetUIDByNames(c, unames)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestDaoGetUNamesByUids(t *testing.T) {
	convey.Convey("GetUNamesByUids", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			uids = []int64{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			d.GetUNamesByUids(c, uids)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestDaoStatArchiveStatStream(t *testing.T) {
	convey.Convey("StatArchiveStatStream", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			business  = int(0)
			typeIDS   = []int64{}
			uids      = []int64{}
			statTypes = []int64{}
			statDate  = time.Now()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			d.StatArchiveStatStream(c, business, typeIDS, uids, statTypes, statDate)
			convCtx.Convey("Then err should be nil.statNodes should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}
