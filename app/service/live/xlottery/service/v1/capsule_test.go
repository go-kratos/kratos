package v1

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	pb "go-common/app/service/live/xlottery/api/grpc/v1"
	"go-common/app/service/live/xlottery/conf"
)

var (
	s *CapsuleService
)

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = NewCapsuleService(conf.Conf)
}

func TestCapsuleService_UpdateCoin(t *testing.T) {
	areaIdsTest := make([]*pb.UpdateCoinConfigReq_AreaIds, 0)
	areaIdTest := &pb.UpdateCoinConfigReq_AreaIds{}
	areaIdTest.ParentId = 2
	areaIdTest.IsAll = 1
	areaIdsTest = append(areaIdsTest, areaIdTest)
	convey.Convey("TestAdd", t, func() {
		res, err := s.UpdateCoinConfig(context.TODO(), &pb.UpdateCoinConfigReq{
			Id:        1,
			Title:     "普通扭蛋",
			ChangeNum: 2,
			StartTime: 1541751430,
			EndTime:   1541891430,
			Status:    0,
			GiftType:  2,
			GiftIds:   []int64{1, 2, 3},
			AreaIds:   areaIdsTest,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCapsuleService_UpdatePool(t *testing.T) {
	convey.Convey("TestUpdatePool", t, func() {
		res, err := s.UpdatePool(context.TODO(), &pb.UpdatePoolReq{
			Id:        13,
			CoinId:    1,
			Title:     "普通扭蛋币aa",
			StartTime: 1,
			EndTime:   2114352000,
			Rule:      "使用价值累计达到10000瓜子的礼物（包含直接使用瓜子购买、道具包裹，但不包括产生梦幻扭蛋币的活动道具），可以获得1枚扭蛋币。使用扭蛋币可以参与抽奖。",
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCapsuleService_UpdatePoolDetail(t *testing.T) {
	convey.Convey("TestUpdatePoolDetail", t, func() {
		res, err := s.UpdatePoolPrize(context.TODO(), &pb.UpdatePoolPrizeReq{
			PoolId:      1,
			Id:          14,
			Type:        6,
			Num:         1,
			ObjectId:    25,
			Expire:      0,
			WebUrl:      "https://i0.hdslb.com/bfs/live/8e7a4dc8de374faee22fca7f9a3f801a1712a36b.png",
			MobileUrl:   "https://i0.hdslb.com/bfs/live/8e7a4dc8de374faee22fca7f9a3f801a1712a36b.png",
			Description: "小电视是一种直播虚拟礼物，可以在直播间送给自己喜爱的主播哦~",
			JumpUrl:     "https://i0.hdslb.com/bfs/live/8e7a4dc8de374faee22fca7f9a3f801a1712a36b.png",
			ProType:     1,
			Chance:      6000,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCapsuleService_UpdateCoinStatus(t *testing.T) {
	convey.Convey("TestUpdateCoinStatus", t, func() {
		res, err := s.UpdateCoinStatus(context.TODO(), &pb.UpdateCoinStatusReq{
			Id:     1,
			Status: 1,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCapsuleService_DeleteCoin(t *testing.T) {
	convey.Convey("TestDeleteCoin", t, func() {
		res, err := s.DeleteCoin(context.TODO(), &pb.DeleteCoinReq{
			Id: 1,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCapsuleService_DeletePool(t *testing.T) {
	convey.Convey("TestDeletePool", t, func() {
		res, err := s.DeletePool(context.TODO(), &pb.DeletePoolReq{
			Id: 7,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCapsuleService_UpdatePoolStatus(t *testing.T) {
	convey.Convey("TestUpdatePoolStatus", t, func() {
		res, err := s.UpdatePoolStatus(context.TODO(), &pb.UpdatePoolStatusReq{
			Id:     1,
			Status: 1,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCapsuleService_GetPoolDetail(t *testing.T) {
	convey.Convey("TestGetPoolDetail", t, func() {
		res, err := s.GetPoolPrize(context.TODO(), &pb.GetPoolPrizeReq{
			PoolId: 1,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCapsuleService_DeletePoolDetail(t *testing.T) {
	convey.Convey("TestDeletePool", t, func() {
		res, err := s.DeletePoolPrize(context.TODO(), &pb.DeletePoolPrizeReq{
			Id: 2,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCapsuleService_GetPoolConfig(t *testing.T) {
	convey.Convey("TestGetPool", t, func() {
		res, err := s.GetPoolList(context.TODO(), &pb.GetPoolListReq{
			Page:     1,
			PageSize: 16,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestCapsuleService_GetCoinList(t *testing.T) {
	convey.Convey("TestGetCoinList", t, func() {
		res, err := s.GetCoinList(context.TODO(), &pb.GetCoinListReq{
			Page:     1,
			PageSize: 16,
		})
		fmt.Println(res, err)
		t.Logf("%v,%s", res, err)
		convey.So(err, convey.ShouldBeNil)
	})
}
