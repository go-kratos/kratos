package v1

import (
	"context"
	"flag"
	"fmt"
	v1pb "go-common/app/service/live/gift/api/grpc/v1"
	"go-common/app/service/live/gift/conf"
	"go-common/app/service/live/gift/dao"
	"go-common/app/service/live/resource/sdk"

	"go-common/app/service/live/gift/model"
	"os"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

var (
	s *GiftService
	//ctx = context.TODO()
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "")
		flag.Set("conf_token", "")
		flag.Set("tree_id", "")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = NewGiftService(conf.Conf)
	dao.InitApi()
	s.TickerReloadGift()
	titansSdk.Init(conf.Conf.Titan)
	os.Exit(m.Run())
}

func TestV1RoomGiftList(t *testing.T) {
	convey.Convey("RoomGiftList", t, func(c convey.C) {
		var (
			ctx = context.Background()
			req = &v1pb.RoomGiftListReq{RoomId: 1}
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			resp, err := s.RoomGiftList(ctx, req)
			c.Convey("Then err should be nil.resp should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(resp, convey.ShouldNotBeNil)
				c.So(resp.List, convey.ShouldNotBeNil)
				c.So(resp.ShowCountMap, convey.ShouldEqual, 0)
				c.So(resp.OldList, convey.ShouldHaveLength, 1)
				c.So(resp.OldList[0].Id, convey.ShouldEqual, 1)
			})
		})
	})
}

func TestV1GiftConfig(t *testing.T) {
	convey.Convey("GiftConfig", t, func(c convey.C) {
		var (
			ctx = context.Background()
			req = &v1pb.GiftConfigReq{
				Platform: "pc",
				Build:    0,
			}
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			resp, err := s.GiftConfig(ctx, req)
			c.Convey("Then err should be nil.resp should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(resp, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestV1TickerReloadGift(t *testing.T) {
	convey.Convey("TickerReloadGift", t, func(c convey.C) {
		c.Convey("When everything gose positive", func(c convey.C) {
			s.TickerReloadGift()
			c.Convey("No return values", func(c convey.C) {
			})
		})
	})
}

func TestV1DiscountGiftList(t *testing.T) {
	convey.Convey("DiscountGiftList", t, func(c convey.C) {
		var (
			ctx = context.Background()
		)
		OK, ok := discountAtomic.Load().(*DiscountCache)
		fmt.Printf("$$$%v%v", OK, ok)
		dCache.discountPlan = map[int64]*model.DiscountInfo{
			1: {
				Id: 1, SceneKey: 99, SceneValue: []int64{0}, Platform: 0,
				List: map[int64]map[string]interface{}{
					20003: {
						"corner_mark_0":     "-20%",
						"corner_position_0": 2,
						"discount_id":       1,
						"discount_price_0":  360000,
					},
				},
			},
			63: {
				Id: 63, SceneKey: 0, SceneValue: []int64{0}, Platform: 4,
				List: map[int64]map[string]interface{}{
					20012: {
						"corner_mark_1":     "-20%",
						"corner_position_1": 2,
						"discount_id":       1,
						"discount_price_1":  360000,
					},
				},
			},
			65: {
				Id: 65, SceneKey: 0, SceneValue: []int64{0}, Platform: 4,
				List: map[int64]map[string]interface{}{
					20013: {
						"corner_mark_1":     "-20%",
						"corner_position_1": 2,
						"discount_id":       1,
						"discount_price_1":  360000,
					},
				},
			},
		}
		dCache.discountMap = map[int64]map[int64]map[int64]int64{
			0: {2: {89: 63}, 99: {0: 1}},
			1: {3: {459879: 63, 460460: 63}, 99: {0: 1}},
			2: {3: {459879: 63, 460460: 63}, 99: {0: 1}},
			4: {3: {123456: 65, 460460: 63}, 99: {0: 1}},
		}
		dCache.cacheTime = time.Now().Unix()
		//discountAtomic.Store(dCache)
		resp, err := s.DiscountGiftList(ctx, &v1pb.DiscountGiftListReq{
			Uid:      1,
			Roomid:   1,
			AreaV2Id: 89,
		})
		c.So(err, convey.ShouldBeNil)
		c.So(resp.DiscountList[0].GiftId, convey.ShouldEqual, 20012)
		c.So(resp.DiscountList[0].DiscountPrice, convey.ShouldEqual, 360000)
		c.So(resp.DiscountList[0].CornerMark, convey.ShouldEqual, "-20%")
		c.So(resp.DiscountList[0].CornerColor, convey.ShouldEqual, "#699553")

		//resp, err = s.DiscountGiftList(ctx, &v1pb.DiscountGiftListReq{
		//	Uid:      1,
		//	Roomid:   123456,
		//	AreaV2Id: 89,
		//	Platform: "android",
		//})
		//c.So(err, convey.ShouldBeNil)
		//c.So(resp.DiscountList[0].GiftId, convey.ShouldEqual, 20013)
		//c.So(resp.DiscountList[0].DiscountPrice, convey.ShouldEqual, 360000)
		//c.So(resp.DiscountList[0].CornerMark, convey.ShouldEqual, "-20%")
		//c.So(resp.DiscountList[0].CornerColor, convey.ShouldEqual, "#699553")

	})
}

func TestV1DailyBag(t *testing.T) {
	convey.Convey("DailyBag", t, func(c convey.C) {
		var (
			ctx = context.Background()
			req = &v1pb.DailyBagReq{Uid: 88895029}
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			resp, err := s.DailyBag(ctx, req)
			c.Convey("Then err should be nil.resp should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(resp, convey.ShouldNotBeNil)
				c.So(resp.BagStatus, convey.ShouldNotBeNil)
				c.So(resp.BagExpireStatus, convey.ShouldNotBeNil)
				c.So(resp.BagToast, convey.ShouldNotBeNil)
				c.So(resp.BagList, convey.ShouldNotBeNil)
			})
		})
	})
}
