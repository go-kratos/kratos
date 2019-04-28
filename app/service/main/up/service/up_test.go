package service

import (
	"context"
	"testing"

	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/model"
	bm "go-common/library/net/http/blademaster"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceEdit(t *testing.T) {
	convey.Convey("Edit", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			mid      = int64(0)
			isAuthor = int(0)
			from     = uint8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			row, err := s.Edit(c, mid, isAuthor, from)
			convCtx.Convey("Then err should be nil.row should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(row, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceInfo(t *testing.T) {
	convey.Convey("Info", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			from = uint8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			isAuthor, err := s.Info(c, mid, from)
			convCtx.Convey("Then err should be nil.isAuthor should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(isAuthor, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceIdentifyAll(t *testing.T) {
	convey.Convey("IdentifyAll", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			ip  = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ia, err := s.IdentifyAll(c, mid, ip)
			convCtx.Convey("Then err should be nil.ia should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ia, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpsByGroup(t *testing.T) {
	convey.Convey("UpsByGroup", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			group = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ups := s.UpsByGroup(c, group)
			convCtx.Convey("Then ups should not be nil.", func(convCtx convey.C) {
				convCtx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceSpecialDel(t *testing.T) {
	convey.Convey("SpecialDel", t, func(convCtx convey.C) {
		var (
			c  = &bm.Context{}
			id = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affectedRow, err := s.SpecialDel(c, id)
			convCtx.Convey("Then err should be nil.affectedRow should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceSpecialAdd(t *testing.T) {
	convey.Convey("SpecialAdd", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			adminName = ""
			special   = &model.UpSpecial{}
			mids      = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affectedRow, err := s.SpecialAdd(c, adminName, special, mids)
			convCtx.Convey("Then err should be nil.affectedRow should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceSpecialEdit(t *testing.T) {
	convey.Convey("SpecialEdit", t, func(convCtx convey.C) {
		var (
			c       = &bm.Context{}
			special = &model.UpSpecial{}
			id      = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			affectedRow, err := s.SpecialEdit(c, special, id)
			convCtx.Convey("Then err should be nil.affectedRow should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceSpecialGet(t *testing.T) {
	convey.Convey("SpecialGet", t, func(convCtx convey.C) {
		var (
			c   = &bm.Context{}
			arg = &model.GetSpecialArg{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, total, err := s.SpecialGet(c, arg)
			convCtx.Convey("Then err should be nil.res,total should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(total, convey.ShouldNotBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceListUpBase(t *testing.T) {
	convey.Convey("ListUpBase", t, func(convCtx convey.C) {
		var (
			c        = &bm.Context{}
			size     = int(0)
			lastID   = int64(0)
			activity = []int64{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			mids, newLastID, err := s.ListUpBase(c, size, lastID, activity)
			convCtx.Convey("Then err should be nil.mids,newLastID should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(newLastID, convey.ShouldNotBeNil)
				convCtx.So(mids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpInfoActivitys(t *testing.T) {
	convey.Convey("UpInfoActivitys", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			req = &upgrpc.UpListByLastIDReq{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.UpInfoActivitys(c, req)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceSpecialGetByMid(t *testing.T) {
	convey.Convey("SpecialGetByMid", t, func(convCtx convey.C) {
		var (
			c   = &bm.Context{}
			arg = &model.GetSpecialByMidArg{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.SpecialGetByMid(c, arg)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpSpecial(t *testing.T) {
	convey.Convey("UpSpecial", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			req = &upgrpc.UpSpecialReq{Mid: 27515314}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.UpSpecial(c, req)
			convCtx.Convey("No return values", func(convCtx convey.C) {
				convCtx.So(res, convey.ShouldNotBeNil)
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceUpsSpecial(t *testing.T) {
	convey.Convey("UpsSpecial", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			req = &upgrpc.UpsSpecialReq{Mids: []int64{27515314}}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.UpsSpecial(c, req)
			convCtx.Convey("No return values", func(convCtx convey.C) {
				convCtx.So(res, convey.ShouldNotBeNil)
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
