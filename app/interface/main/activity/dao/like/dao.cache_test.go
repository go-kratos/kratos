package like

import (
	"context"
	"testing"

	"fmt"

	"go-common/app/interface/main/activity/model/like"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeLike(t *testing.T) {
	convey.Convey("Like", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Like(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeActSubject(t *testing.T) {
	convey.Convey("ActSubject", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ActSubject(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeLikeMissionBuff(t *testing.T) {
	convey.Convey("LikeMissionBuff", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(10256)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LikeMissionBuff(c, id, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeMissionGroupItems(t *testing.T) {
	convey.Convey("MissionGroupItems", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			keys = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MissionGroupItems(c, keys)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeActMission(t *testing.T) {
	convey.Convey("ActMission", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(10256)
			lid = int64(7)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ActMission(c, id, lid, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeActLikeAchieves(t *testing.T) {
	convey.Convey("ActLikeAchieves", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ActLikeAchieves(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeActMissionFriends(t *testing.T) {
	convey.Convey("ActMissionFriends", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(10256)
			lid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ActMissionFriends(c, id, lid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeActUserAchieve(t *testing.T) {
	convey.Convey("ActUserAchieve", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ActUserAchieve(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeMatchSubjects(t *testing.T) {
	convey.Convey("MatchSubjects", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			keys = []int64{10256}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MatchSubjects(c, keys)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeLikeContent(t *testing.T) {
	convey.Convey("LikeContent", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			keys = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LikeContent(c, keys)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestActSubjectProtocol(t *testing.T) {
	convey.Convey("LikeContent", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10298)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ActSubjectProtocol(c, sid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestCacheActSubjectProtocol(t *testing.T) {
	convey.Convey("LikeContent", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10298)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheActSubjectProtocol(c, sid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestAddCacheActSubjectProtocol(t *testing.T) {
	convey.Convey("LikeContent", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			sid      = int64(10256)
			protocol = &like.ActSubjectProtocol{ID: 1, Sid: 10256}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheActSubjectProtocol(c, sid, protocol)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
