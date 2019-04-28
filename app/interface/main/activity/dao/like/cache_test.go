package like

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikelikeKey(t *testing.T) {
	convey.Convey("likeKey", t, func(ctx convey.C) {
		var (
			id = int64(7)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := likeKey(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeactSubjectKey(t *testing.T) {
	convey.Convey("actSubjectKey", t, func(ctx convey.C) {
		var (
			id = int64(7)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := actSubjectKey(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeactSubjectMaxIDKey(t *testing.T) {
	convey.Convey("actSubjectMaxIDKey", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := actSubjectMaxIDKey()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikelikeMaxIDKey(t *testing.T) {
	convey.Convey("likeMaxIDKey", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := likeMaxIDKey()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikelikeMissionBuffKey(t *testing.T) {
	convey.Convey("likeMissionBuffKey", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := likeMissionBuffKey(sid, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikelikeMissionGroupIDkey(t *testing.T) {
	convey.Convey("likeMissionGroupIDkey", t, func(ctx convey.C) {
		var (
			lid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := likeMissionGroupIDkey(lid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikelikeActMissionKey(t *testing.T) {
	convey.Convey("likeActMissionKey", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
			lid = int64(1)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := likeActMissionKey(sid, lid, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeactAchieveKey(t *testing.T) {
	convey.Convey("actAchieveKey", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := actAchieveKey(sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeactMissionFriendsKey(t *testing.T) {
	convey.Convey("actMissionFriendsKey", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
			lid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := actMissionFriendsKey(sid, lid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeactUserAchieveKey(t *testing.T) {
	convey.Convey("actUserAchieveKey", t, func(ctx convey.C) {
		var (
			id = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := actUserAchieveKey(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeactUserAchieveAwardKey(t *testing.T) {
	convey.Convey("actUserAchieveAwardKey", t, func(ctx convey.C) {
		var (
			id = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := actUserAchieveAwardKey(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikesubjectStatKey(t *testing.T) {
	convey.Convey("subjectStatKey", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := subjectStatKey(sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeviewRankKey(t *testing.T) {
	convey.Convey("viewRankKey", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := viewRankKey(sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikelikeContentKey(t *testing.T) {
	convey.Convey("likeContentKey", t, func(ctx convey.C) {
		var (
			lid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := likeContentKey(lid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestSubjectProtocolKey(t *testing.T) {
	convey.Convey("subjectProtocolKey", t, func(ctx convey.C) {
		var (
			sid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := subjectProtocolKey(sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
