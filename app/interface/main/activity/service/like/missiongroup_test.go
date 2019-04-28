package like

import (
	"context"
	"fmt"
	"testing"

	l "go-common/app/interface/main/activity/model/like"

	"encoding/json"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_MissionLike(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		lid, err := svf.MissionLike(context.Background(), 10293, 15555180)
		So(err, ShouldBeNil)
		fmt.Printf("%d", lid)
	}))
}

func TestService_MissionInfo(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		res, err := svf.MissionInfo(context.Background(), 10292, 2, 442549)
		So(err, ShouldBeNil)
		fmt.Printf("%+v", res)
	}))
}

func TestService_MissionLikeAct(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		res, err := svf.MissionLikeAct(context.Background(), &l.ParamMissionLikeAct{Sid: 10292, Lid: 17}, 216761)
		So(err, ShouldBeNil)
		bs, _ := json.Marshal(res)
		Printf("%+v", string(bs))
	}))
}

func TestService__MissionTops(t *testing.T) {
	Convey("test set like act", t, WithService(func(svf *Service) {
		acts, err := svf.MissionTops(context.Background(), 10292, 50)
		So(err, ShouldBeNil)
		for _, v := range acts {
			fmt.Printf("%+v", v)
		}
	}))
}

func TestService_CalculateAchievement(t *testing.T) {
	Convey("test set like act", t, WithService(func(svf *Service) {
		err := svf.CalculateAchievement(context.Background(), 10292, 15555180, 10)
		So(err, ShouldBeNil)
	}))
}

func TestService_MissionRank(t *testing.T) {
	Convey("test set like act", t, WithService(func(svf *Service) {
		res, err := svf.MissionRank(context.Background(), 10292, 14137123)
		So(err, ShouldBeNil)
		fmt.Printf("%+v", res)
	}))
}

func TestService_MissionFriendsList(t *testing.T) {
	Convey("test set like act", t, WithService(func(svf *Service) {
		res, err := svf.MissionFriendsList(context.Background(), &l.ParamMissionFriends{Sid: 10292, Lid: 11, Size: 5}, 216761)
		So(err, ShouldBeNil)
		for _, v := range res {
			fmt.Printf("%+v", v)
		}
	}))
}

func TestService_MissionAward(t *testing.T) {
	Convey("test set like act", t, WithService(func(svf *Service) {
		res, err := svf.MissionAward(context.Background(), 10292, 2089809)
		So(err, ShouldBeNil)
		for _, v := range res {
			fmt.Printf("%+v", v)
		}
	}))
}

func TestService_MissionAchieve(t *testing.T) {
	Convey("test set like act", t, WithService(func(svf *Service) {
		res, err := svf.MissionAchieve(context.Background(), 10292, 4, 2089809)
		So(err, ShouldBeNil)
		Printf("%d", res)
	}))
}

func TestService_MissionUser(t *testing.T) {
	Convey("test set like act", t, WithService(func(svf *Service) {
		res, err := svf.MissionUser(context.Background(), 10292, 12)
		So(err, ShouldBeNil)
		Printf("%+v", res)
	}))
}
