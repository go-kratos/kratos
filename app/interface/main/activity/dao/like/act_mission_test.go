package like

import (
	"context"
	"fmt"
	"testing"

	l "go-common/app/interface/main/activity/model/like"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikemissionLikeLimitKey(t *testing.T) {
	convey.Convey("missionLikeLimitKey", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := missionLikeLimitKey(sid, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikelikeActMissionScoreKey(t *testing.T) {
	convey.Convey("likeActMissionScoreKey", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := likeActMissionScoreKey(sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikelikeMissionScoreStrKey(t *testing.T) {
	convey.Convey("likeMissionScoreStrKey", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
			lid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := likeMissionScoreStrKey(sid, lid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikescoreMaxTime(t *testing.T) {
	convey.Convey("scoreMaxTime", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := scoreMaxTime()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikescoreMaxNum(t *testing.T) {
	convey.Convey("scoreMaxNum", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := scoreMaxNum()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikebuildRankScore(t *testing.T) {
	convey.Convey("buildRankScore", t, func(ctx convey.C) {
		var (
			num   = int64(3)
			ctime = int64(15487588)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := buildRankScore(num, ctime)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeRawActMission(t *testing.T) {
	convey.Convey("RawActMission", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
			lid = int64(123)
			mid = int64(123)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawActMission(c, sid, lid, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeMissionLikeLimit(t *testing.T) {
	convey.Convey("MissionLikeLimit", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MissionLikeLimit(c, sid, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeInrcMissionLikeLimit(t *testing.T) {
	convey.Convey("InrcMissionLikeLimit", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
			mid = int64(77)
			val = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.InrcMissionLikeLimit(c, sid, mid, val)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeSetMissionTop(t *testing.T) {
	convey.Convey("SetMissionTop", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10256)
			lid   = int64(123)
			score = int64(1)
			ctime = int64(158748596)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.SetMissionTop(c, sid, lid, score, ctime)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", count)
			})
		})
	})
}

func TestLikeMissionLidScore(t *testing.T) {
	convey.Convey("MissionLidScore", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
			lid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			score, err := d.MissionLidScore(c, sid, lid)
			ctx.Convey("Then err should be nil.score should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(score, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeMissionLidRank(t *testing.T) {
	convey.Convey("MissionLidRank", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
			lid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rank, err := d.MissionLidRank(c, sid, lid)
			ctx.Convey("Then err should be nil.rank should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rank, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeMissionScoreList(t *testing.T) {
	convey.Convey("MissionScoreList", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10256)
			start = int(1)
			end   = int(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			data, err := d.MissionScoreList(c, sid, start, end)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(data, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeAddActMission(t *testing.T) {
	convey.Convey("AddActMission", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			act = &l.ActMissionGroup{Sid: 10256, Lid: 1235, Mid: 77, IPv6: make([]byte, 0)}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			actID, err := d.AddActMission(c, act)
			ctx.Convey("Then err should be nil.actID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(actID, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeRawActMissionFriends(t *testing.T) {
	convey.Convey("RawActMissionFriends", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
			lid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawActMissionFriends(c, sid, lid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
