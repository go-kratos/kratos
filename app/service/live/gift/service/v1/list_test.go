package v1

import (
	"context"
	v1pb "go-common/app/service/live/gift/api/grpc/v1"
	"go-common/app/service/live/gift/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestV1getOnlinePlanGiftList(t *testing.T) {
	convey.Convey("getOnlinePlanGiftList", t, func(c convey.C) {
		var (
			ctx          = context.Background()
			roomID       = int64(0)
			areaParentID = int64(0)
			areaID       = int64(0)
			platform     = ""
			build        = int64(0)
			mobiApp      = "android"
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			resp, err := s.GetOnlinePlanGiftList(ctx, roomID, areaParentID, areaID, platform, build, mobiApp)
			c.Convey("Then err should be nil.resp should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(resp, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1LoadGiftCache(t *testing.T) {
	convey.Convey("LoadGiftCache", t, func(c convey.C) {
		var (
			ctx        = context.Background()
			needReload bool
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := s.LoadGiftCache(ctx, needReload)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestV1SyncLocalCache(t *testing.T) {
	convey.Convey("SyncLocalCache", t, func(c convey.C) {
		var (
			ctx = context.Background()
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			err := s.SyncLocalCache(ctx)
			c.Convey("Then err should be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestV1isPrivilegeGift(t *testing.T) {
	convey.Convey("isPrivilegeGift", t, func(c convey.C) {
		var (
			giftID = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := s.isPrivilegeGift(giftID)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1parsePlatform(t *testing.T) {
	convey.Convey("parsePlatform", t, func(c convey.C) {
		pf := s.parsePlatform(0)
		c.So(pf, convey.ShouldNotBeNil)
		c.So(pf, convey.ShouldHaveLength, 1)

		c.So(pf, convey.ShouldContain, "default")

		pf = s.parsePlatform(6)
		c.So(pf, convey.ShouldNotBeNil)
		c.So(pf, convey.ShouldHaveLength, 2)
		c.So(pf, convey.ShouldContain, "ios")
		c.So(pf, convey.ShouldContain, "android")
	})
}

func TestV1addDefaultPlan(t *testing.T) {
	convey.Convey("addDefaultPlan", t, func(c convey.C) {
		var (
			plan = make(map[string]map[int64]map[int64]*model.GiftPlan)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			s.addDefaultPlan(plan)
			c.Convey("No return values", func(c convey.C) {
			})
		})
	})
}

func TestV1ConvertGiftList(t *testing.T) {
	convey.Convey("ConvertGiftList", t, func(c convey.C) {
		var (
			list     = "1,25"
			planID   = int64(0)
			platform = "pc"
			build    = int64(0)
			mobiApp  = "android"
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			data := s.ConvertGiftList(list, planID, platform, build, mobiApp)
			c.Convey("Then data should not be nil.", func(c convey.C) {
				c.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1CheckGiftVersion(t *testing.T) {
	convey.Convey("CheckGiftVersion", t, func(c convey.C) {
		var (
			giftID   = int64(30047)
			platform = ""
			build    = int64(0)
			mobiApp  = "android"
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := s.CheckGiftVersion(giftID, platform, build, mobiApp)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1AddPrivilegeGift(t *testing.T) {
	convey.Convey("AddPrivilegeGift", t, func(c convey.C) {
		var (
			list     = []*v1pb.RoomGiftListResp_List{}
			platform = ""
			build    = int64(0)
			mobiApp  = "android"
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := s.AddPrivilegeGift(list, platform, build, mobiApp)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1NewGiftList(t *testing.T) {
	convey.Convey("NewGiftList", t, func(c convey.C) {
		var (
			id       = int64(0)
			position = int64(0)
			planID   = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			list := s.NewGiftList(id, position, planID)
			c.Convey("Then list should not be nil.", func(c convey.C) {
				c.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1getOldList(t *testing.T) {
	convey.Convey("getOldList", t, func(c convey.C) {
		c.Convey("When everything gose positive", func(c convey.C) {
			old := s.getOldList()
			c.Convey("Then old should not be nil.", func(c convey.C) {
				c.So(old, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetAllConfig(t *testing.T) {
	convey.Convey("GetAllConfig", t, func(c convey.C) {
		var (
			ctx = context.Background()
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			all, err := s.GetAllConfig(ctx)
			c.Convey("Then err should be nil.all should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(all, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1ConvertDB2Config(t *testing.T) {
	convey.Convey("ConvertDB2Config", t, func(c convey.C) {
		var (
			ol = &model.GiftOnline{}
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res := s.ConvertDB2Config(ol)
			c.Convey("Then res should not be nil.", func(c convey.C) {
				c.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetStayTime(t *testing.T) {
	convey.Convey("GetStayTime", t, func(c convey.C) {
		var (
			ol = &model.GiftOnline{}
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			time := s.GetStayTime(ol)
			c.Convey("Then time should not be nil.", func(c convey.C) {
				c.So(time, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1IsValidGift(t *testing.T) {
	convey.Convey("IsValidGift", t, func(c convey.C) {
		var (
			giftID = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := s.IsValidGift(giftID)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1IsSpecialGift(t *testing.T) {
	convey.Convey("IsSpecialGift", t, func(c convey.C) {
		var (
			giftID = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			p1 := s.IsSpecialGift(giftID)
			c.Convey("Then p1 should not be nil.", func(c convey.C) {
				c.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1GetGiftInfoByID(t *testing.T) {
	convey.Convey("GetGiftInfoByID", t, func(c convey.C) {
		var (
			ctx    = context.Background()
			giftID = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			gift := s.GetGiftInfoByID(ctx, giftID)
			c.Convey("Then gift should not be nil.", func(c convey.C) {
				c.So(gift, convey.ShouldNotBeNil)
			})
		})
	})
}
