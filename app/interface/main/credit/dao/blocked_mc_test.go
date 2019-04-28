package dao

import (
	"context"
	"go-common/app/interface/main/credit/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaouserBlockedListKey(t *testing.T) {
	convey.Convey("userBlockedListKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := userBlockedListKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoblockedInfoKey(t *testing.T) {
	convey.Convey("blockedInfoKey", t, func(convCtx convey.C) {
		var (
			id = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := blockedInfoKey(id)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBlockedUserListCache(t *testing.T) {
	convey.Convey("BlockedUserListCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(-1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ls, err := d.BlockedUserListCache(c, mid)
			convCtx.Convey("Then err should be nil.ls should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(len(ls), convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoSetBlockedUserListCache(t *testing.T) {
	convey.Convey("SetBlockedUserListCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ls  = []*model.BlockedInfo{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetBlockedUserListCache(c, mid, ls)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBlockedInfoCache(t *testing.T) {
	convey.Convey("BlockedInfoCache", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(-1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			info, err := d.BlockedInfoCache(c, id)
			convCtx.Convey("Then err should be nil.info should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(info, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetBlockedInfoCache(t *testing.T) {
	convey.Convey("SetBlockedInfoCache", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			id   = int64(0)
			info = &model.BlockedInfo{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetBlockedInfoCache(c, id, info)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoBlockedInfosCache(t *testing.T) {
	convey.Convey("BlockedInfosCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{234}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			infos, miss, err := d.BlockedInfosCache(c, ids)
			convCtx.Convey("Then err should not be nil.infos,miss should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldNotBeNil)
				convCtx.So(miss, convey.ShouldBeNil)
				convCtx.So(infos, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetBlockedInfosCache(t *testing.T) {
	convey.Convey("SetBlockedInfosCache", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			infos = []*model.BlockedInfo{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetBlockedInfosCache(c, infos)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
