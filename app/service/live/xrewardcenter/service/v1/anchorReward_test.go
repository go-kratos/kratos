package v1

import (
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"go-common/app/service/live/xrewardcenter/api/grpc/v1"
	"go-common/app/service/live/xrewardcenter/conf"
	"testing"
)

var (
	AnchorReward *AnchorRewardService
)

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	AnchorReward = NewAnchorTaskService(conf.Conf)
}

// go test  -test.v -test.run TestAnchorRewardService_MyReward
func TestAnchorRewardService_MyReward(t *testing.T) {
	Convey("TestMyReward", t, func() {

		ctx := metadata.NewContext(bm.Context{}, metadata.MD{
			"mid": 10000,
		})

		res, err := AnchorReward.MyReward(ctx, &v1.AnchorTaskMyRewardReq{
			Page: 1,
			Uid:  10000,
		})
		t.Logf("%+v", res)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestIsViewed
func TestAnchorRewardService_IsViewed(t *testing.T) {
	Convey("TestAnchorRewardService_IsViewed", t, func() {

		ctx := metadata.NewContext(bm.Context{}, metadata.MD{
			"mid": 10000,
		})

		res, err := AnchorReward.IsViewed(ctx, &v1.AnchorTaskIsViewedReq{
			Uid: 10000,
		})
		t.Logf("%+v", res)
		So(err, ShouldBeNil)
	})
}

func TestAnchorRewardService_UseRecord(t *testing.T) {
	Convey("TestAnchorRewardService_UseRecord", t, func() {
		ctx := metadata.NewContext(bm.Context{}, metadata.MD{
			"mid": 10000,
		})

		res, err := AnchorReward.UseRecord(ctx, &v1.AnchorTaskUseRecordReq{
			Page: 0,
			Uid:  10000,
		})
		t.Logf("%+v", res)

		So(err, ShouldBeNil)
	})
}

func TestAnchorRewardService_UseReward(t *testing.T) {
	Convey("TestAnchorRewardService_UseReward", t, func() {
		ctx := metadata.NewContext(bm.Context{}, metadata.MD{
			"mid": 10000,
		})

		res, err := AnchorReward.UseReward(ctx, &v1.AnchorTaskUseRewardReq{
			Id:      1,
			Uid:     10000,
			UsePlat: "ios",
		})
		t.Logf("%+v", res)

		So(err, ShouldNotBeNil)
	})
}

func TestAnchorRewardService_AddReward(t *testing.T) {
	Convey("TestAnchorRewardService_AddReward", t, func() {
		ctx := metadata.NewContext(bm.Context{}, metadata.MD{
			"mid": 10000,
		})

		res, err := AnchorReward.AddReward(ctx, &v1.AnchorTaskAddRewardReq{
			RewardId: 1,
			Roomid:   1,
			Source:   1,
			Uid:      10000,
			OrderId:  "test",
		})
		t.Logf("%+v", res)

		So(err, ShouldBeNil)
	})
}
